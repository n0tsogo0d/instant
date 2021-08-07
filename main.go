package main

import (
	"fmt"
	"log"
	"math"
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

		media = append(media, m3u8.MediaItem{
			Type:       m3u8.TypeAudio,
			GroupID:    "stereo",
			Name:       "English",
			Language:   "en",
			Default:    true,
			Autoselect: true,
			URI: fmt.Sprintf("playlist?type=audio&file=%s&bitrate=%d",
				probe.Format.Filename, 194000),
		})

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

		bitrate := 7_800_000
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
					AverageBandwidth: bitrate,
					Bandwidth:        bitrate,
					Codecs:           "avc1.640029,mp4a.40.2",
					Framerate:        framerate,
					VideoRange:       "SDR",
					Subtitles:        "subtitles",
					Audio:            "stereo",
					URI: fmt.Sprintf("playlist?type=video&file=%s&height=%d"+
						"&width=%d&bitrate=%d", probe.Format.Filename, height,
						width, bitrate),
				},
			},
		}
	case "video":
		length, _ := strconv.ParseFloat(probe.Format.Duration, 64)
		segments := make([]m3u8.SegmentItem, 0)
		duration := 6.00
		height := q.Get("height")
		width := q.Get("width")
		bitrate := q.Get("bitrate")

		for i := 0.00; i < length; i += duration {
			segments = append(segments, m3u8.SegmentItem{
				Duration: duration,
				URI: fmt.Sprintf("segment?type=video&file=%s&start=%.4f"+
					"&duration=%.4f&height=%s&width=%s&bitrate=%s",
					probe.Format.Filename, i, duration, height, width, bitrate),
			})
		}

		// this is important for the last segment
		// TODO: This has the wrong start, need to fix...
		remainder := math.Mod(length, duration)
		if remainder > 0 {
			segments = append(segments, m3u8.SegmentItem{
				Duration: duration,
				URI: fmt.Sprintf("segment?type=video&file=%s&start=%.4f"+
					"&duration=%.4f&height=%s&width=%s&bitrate=%s",
					probe.Format.Filename, duration-remainder, remainder,
					height, width, bitrate),
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
		duration, _ := strconv.ParseFloat(probe.Format.Duration, 64)

		playlist = m3u8.Playlist{
			Version:        3,
			Type:           "VOD",
			TargetDuration: duration,
			Segments: []m3u8.SegmentItem{{
				Duration: duration,
				URI: fmt.Sprintf("segment?type=subtitle&file=%s&index=%s"+
					"&start=0&duration=%.4f", probe.Format.Filename, index,
					duration),
			}},
		}
	case "audio":
		length, _ := strconv.ParseFloat(probe.Format.Duration, 64)
		duration := 60.0
		bitrate := q.Get("bitrate")
		segments := make([]m3u8.SegmentItem, 0)

		for i := 0.00; i < length; i += duration {
			segments = append(segments, m3u8.SegmentItem{
				Duration: duration,
				URI: fmt.Sprintf("segment?type=audio&file=%s&start=%.4f"+
					"&duration=%.4f&bitrate=%s", probe.Format.Filename, i,
					duration, bitrate),
			})
		}

		// this is important for the last segment
		// TODO: This has the wrong start, need to fix...
		remainder := math.Mod(length, duration)
		if remainder > 0 {
			segments = append(segments, m3u8.SegmentItem{
				Duration: duration,
				URI: fmt.Sprintf("segment?type=audio&file=%s&start=%.4f"+
					"&duration=%.4f&bitrate=%s", probe.Format.Filename,
					duration-remainder, remainder, bitrate),
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
		bitrate, _ := strconv.Atoi(q.Get("bitrate"))

		res, err = ffmpeg.CreateVideoSegment(ffmpeg.CreateVideoSegmentOptions{
			Input:    input,
			Start:    start,
			Duration: duration,
			Bitrate:  bitrate,
			Height:   height,
			Width:    width,
		})

	case "subtitle":
		index, _ := strconv.Atoi(q.Get("index"))

		res, err = ffmpeg.CreateSubtitleSegment(ffmpeg.CreateSubtitleSegmentOptions{
			Input:    input,
			Start:    start,
			Duration: duration,
			Index:    index,
		})
	case "audio":
		bitrate, _ := strconv.Atoi(q.Get("bitrate"))

		res, err = ffmpeg.CreateAudioSegment(ffmpeg.CreateAudioSegmentOptions{
			Input:    input,
			Start:    start,
			Duration: duration,
			Bitrate:  bitrate,
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
