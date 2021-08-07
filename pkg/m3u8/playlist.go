// m3u8 generates m3u playlist files from user input
//
// master playlist format:
// #EXTM3U
// #EXT-X-VERSION:3
// #EXT-X-MEDIA:TYPE=SUBTITLES,GROUP-ID="subtitles",NAME="English",DEFAULT=NO,AUTOSELECT=NO,FORCED=NO,LANGUAGE="en",URI="subtitle_playlist.m3u8"
//
// #EXT-X-STREAM-INF:BANDWIDTH=10000000,AVERAGE-BANDWIDTH=10000000,VIDEO-RANGE=SDR,CODECS="avc1.640029,mp4a.40.2",RESOLUTION=1920x1080,FRAME-RATE=23.975,SUBTITLES="subtitles"
// playlist.m3u8
//
//
// video/audio/subtitle playlist format:
// #EXTM3U
// #EXT-X-PLAYLIST-TYPE:VOD
// #EXT-X-VERSION:3
// #EXT-X-TARGETDURATION:10
// #EXT-X-MEDIA-SEQUENCE:0
// #EXTINF:10.0000,
// segment.ts
// #EXTINF:10.0000,
// #EXT-X-ENDLIST
package m3u8

import "fmt"

const (
	TypeSubtitles = "SUBTITLES"
	TypeAudio     = "AUDIO"

	CodecVideo    = "video"
	CodecAudio    = "audio"
	CodecSubtitle = "subtitle"
)

type (
	Playlist struct {
		Type           string
		Version        int
		TargetDuration float64
		MediaSequence  int
		Segments       []SegmentItem
		Playlists      []PlaylistItem
		Media          []MediaItem
	}
	MediaItem struct {
		Type       string
		GroupID    string
		Name       string
		Default    bool
		Autoselect bool
		Forced     bool
		Language   string
		URI        string
	}
	PlaylistItem struct {
		Bandwidth        int
		AverageBandwidth int
		VideoRange       string
		Codecs           string
		Framerate        float64
		Subtitles        string
		Audio            string
		URI              string
		Width            int
		Height           int
	}
	SegmentItem struct {
		Duration float64
		URI      string
	}
)

func (p Playlist) String() string {
	str := "#EXTM3U\n"
	str += fmt.Sprintf("#EXT-X-VERSION:%d\n", p.Version)

	if p.Segments == nil {
		str += "#EXT-X-INDEPENDENT-SEGMENTS\n"
		str += "\n"
		for _, m := range p.Media {
			str += fmt.Sprintf("#EXT-X-MEDIA:TYPE=%s,GROUP-ID=\"%s\","+
				"LANGUAGE=\"%s\",NAME=\"%s\",AUTOSELECT=%s,DEFAULT=%s,"+
				"URI=\"%s\"\n", m.Type, m.GroupID, m.Language, m.Name,
				boolToyesNo(m.Autoselect), boolToyesNo(m.Default), m.URI)
		}
		str += "\n"

		// todo work this out a bit better...
		for _, p := range p.Playlists {
			str += fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%d,"+
				"AVERAGE-BANDWIDTH=%d,VIDEO-RANGE=%s,CODECS=\"%s\","+
				"RESOLUTION=%dx%d,FRAME-RATE=%.3f",
				p.Bandwidth, p.AverageBandwidth, p.VideoRange, p.Codecs,
				p.Width, p.Height, p.Framerate)

			// add audio if available
			if len(p.Audio) > 0 {
				str += fmt.Sprintf(",AUDIO=\"%s\"", p.Audio)
			}

			// add subtitle group if available
			if len(p.Subtitles) > 0 {
				str += fmt.Sprintf(",SUBTITLES=\"%s\"", p.Subtitles)
			}

			str += "\n"
			str += p.URI
			str += "\n"
		}

		return str
	}

	str += fmt.Sprintf("#EXT-X-TARGETDURATION:%.4f\n", p.TargetDuration)
	for _, s := range p.Segments {
		str += fmt.Sprintf("#EXTINF:%.4f,\n", s.Duration)
		str += fmt.Sprintf("%s\n", s.URI)
	}
	str += "#EXT-X-ENDLIST\n"

	return str
}

func (p Playlist) Bytes() []byte {
	return []byte(p.String())
}

func boolToyesNo(b bool) string {
	if b {
		return "YES"
	}

	return "NO"
}
