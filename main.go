package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/n0tsogo0d/instant/pkg/ffmpeg"
	"github.com/n0tsogo0d/instant/pkg/m3u8"
)

func main() {
	r := http.NewServeMux()
	r.Handle("/", http.FileServer(http.Dir("web/")))
	r.HandleFunc("/api/poster", handlePoster)
	r.HandleFunc("/api/playlist", handlePlaylist)
	r.HandleFunc("/api/segment", handleSegment)

	srv := &http.Server{
		Handler: r,
		Addr:    ":8000",

		// long timeout since it could be that some segments
		// take longer, but it hasn't been tested yet
		WriteTimeout: 120 * time.Second,
		ReadTimeout:  120 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func handlePoster(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("file")

	probe, err := ffmpeg.Probe(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration, _ := strconv.ParseFloat(probe.Format.Duration, 64)
	res, err := ffmpeg.CreatePoster(ffmpeg.CreatePosterOptions{
		Input: input,
		Time:  duration / 2,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(res)
}

// handlePlaylist returns three types of playlists:
// - primary playlist: this contains information for all sub-playlists
//   like multiple video resolutions, audio streams, subtitle tracks, etc.
// - video playlist
// - audio playlist
// - subtitle playlist
// For more information simply stream something, open the developer
// console and see what requests are being made. It's actually pretty logical.
func handlePlaylist(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	file := q.Get("file")

	probe, err := ffmpeg.Probe(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var playlist m3u8.Playlist
	switch q.Get("type") {
	case "primary":
		// create master playlist
		var video ffmpeg.Stream
		media := make([]m3u8.MediaItem, 0)

		// exract all subtitles
		for _, s := range probe.Streams {
			if s.CodecType != m3u8.CodecSubtitle {
				// set video stream
				if s.CodecType == m3u8.CodecVideo {
					video = s
				}
				continue
			}

			media = append(media, m3u8.MediaItem{
				Type:       m3u8.TypeSubtitles,
				GroupID:    "subtitles",
				Name:       s.Tags.Title,
				Language:   s.Tags.Language,
				Default:    false,
				Autoselect: false,
				Forced:     false,
				URI: fmt.Sprintf("playlist?type=subtitle&file=%s&index=%d",
					probe.Format.Filename, s.Index),
			})
		}

		// check if there is a video stream
		if video.CodecType != m3u8.CodecVideo {
			http.Error(w, "no video stream", http.StatusInternalServerError)
			return
		}

		// framerate
		split := strings.Split(video.AvgFrameRate, "/")
		f1, _ := strconv.ParseFloat(split[0], 64)
		f2, _ := strconv.ParseFloat(split[1], 64)
		framerate := f1 / f2

		videoBitrate := 7_800_000
		audioBitrate := 192_000
		height := video.Height
		width := video.Width

		playlist = m3u8.Playlist{
			Version: 3,
			Type:    "VOD",
			Media:   media,
			Playlists: []m3u8.PlaylistItem{
				{
					Height:           height,
					Width:            width,
					AverageBandwidth: videoBitrate,
					Bandwidth:        videoBitrate,
					Codecs:           "avc1.640029,mp4a.40.2",
					Framerate:        framerate,
					VideoRange:       "SDR",
					Subtitles:        "subtitles",
					URI: fmt.Sprintf("playlist?type=video&file=%s&height=%d"+
						"&width=%d&video_bitrate=%d&audio_bitrate=%d",
						probe.Format.Filename, height, width, videoBitrate,
						audioBitrate),
				},
			},
		}
	case "video":
		length, _ := strconv.ParseFloat(probe.Format.Duration, 64)
		segments := make([]m3u8.SegmentItem, 0)
		duration := 6.00
		height := q.Get("height")
		width := q.Get("width")
		videoBitrate := q.Get("video_bitrate")
		audioBitrate := q.Get("audio_bitrate")

		for _, v := range m3u8.PlaylistSegments(length, duration) {
			segments = append(segments, m3u8.SegmentItem{
				Duration: v[1],
				URI: fmt.Sprintf("segment?type=video&file=%s&start=%.4f"+
					"&duration=%.4f&height=%s&width=%s&video_bitrate=%s"+
					"&audio_bitrate=%s", probe.Format.Filename, v[0], v[1],
					height, width, videoBitrate, audioBitrate),
			})
		}

		playlist = m3u8.Playlist{
			Version:        3,
			Type:           "VOD",
			TargetDuration: duration,
			Segments:       segments,
		}
	case "subtitle":
		index := q.Get("index")
		segments := make([]m3u8.SegmentItem, 0)
		length, _ := strconv.ParseFloat(probe.Format.Duration, 64)
		duration := 60.0

		for _, v := range m3u8.PlaylistSegments(length, duration) {
			segments = append(segments, m3u8.SegmentItem{
				Duration: v[1],
				URI: fmt.Sprintf("segment?type=subtitle&file=%s&index=%s"+
					"&start=%.4f&duration=%.4f", probe.Format.Filename, index,
					v[0], v[1]),
			})
		}

		playlist = m3u8.Playlist{
			Version:        3,
			Type:           "VOD",
			TargetDuration: duration,
			Segments:       segments,
		}
	default:
		http.Error(w, "invalid playlist type", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/x-mpegURL")
	w.Write(playlist.Bytes())
}

// handleSegment returns on-demand segments for audio/subtitle/video
func handleSegment(w http.ResponseWriter, r *http.Request) {
	var res []byte
	var err error

	q := r.URL.Query()
	input := q.Get("file")
	start, _ := strconv.ParseFloat(q.Get("start"), 64)
	duration, _ := strconv.ParseFloat(q.Get("duration"), 64)

	switch q.Get("type") {
	case "video":
		height, _ := strconv.Atoi(q.Get("height"))
		width, _ := strconv.Atoi(q.Get("width"))
		videoBitrate, _ := strconv.Atoi(q.Get("video_bitrate"))
		audioBitrate, _ := strconv.Atoi(q.Get("audio_bitrate"))

		res, err = ffmpeg.CreateVideoSegment(ffmpeg.CreateVideoSegmentOptions{
			Input:        input,
			Start:        start,
			Duration:     duration,
			VideoBitrate: videoBitrate,
			AudioBitrate: audioBitrate,
			Height:       height,
			Width:        width,
		})

	case "subtitle":
		index, _ := strconv.Atoi(q.Get("index"))

		res, err = ffmpeg.CreateSubtitleSegment(ffmpeg.CreateSubtitleSegmentOptions{
			Input:    input,
			Start:    start,
			Duration: duration,
			Index:    index,
		})
	default:
		http.Error(w, "invalid segment type", http.StatusInternalServerError)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/octet-stream")
	w.Write(res)
}
