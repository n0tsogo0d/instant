package ffmpeg

import (
	"fmt"
)

type CreateVideoSegmentOptions struct {
	Input        string
	Start        float64
	Duration     float64
	VideoBitrate int
	AudioBitrate int
	Height       int
	Width        int
}
type CreateSubtitleSegmentOptions struct {
	Input    string
	Start    float64
	Duration float64
	Index    int
}

func CreateVideoSegment(opts CreateVideoSegmentOptions) ([]byte, error) {
	args := fmt.Sprintf("-loglevel error -hide_banner "+
		// input
		"-ss %.4f -to %.4f -i %s "+
		// subtitles
		" -sn "+
		// video
		"-c:v libx264 "+
		"-vf scale=w=%d:h=%d:force_original_aspect_ratio=decrease "+
		"-b:v %d -minrate %d -maxrate %d -bufsize %d -profile:v main "+
		// audio
		"-c:a aac -b:a %d -ac 2 "+
		// output
		"-crf 26 -copyts -preset ultrafast -muxpreload 0 -muxdelay 0 pipe:.ts",
		opts.Start, opts.Start+opts.Duration, opts.Input, opts.Width,
		opts.Height, opts.VideoBitrate, opts.VideoBitrate, opts.VideoBitrate,
		opts.VideoBitrate/2, opts.AudioBitrate)

	return pipe(args)
}

func CreateSubtitleSegment(opts CreateSubtitleSegmentOptions) ([]byte, error) {
	args := fmt.Sprintf("-loglevel error -hide_banner "+
		"-ss %.4f -to %.4f -i %s -map 0:%d "+
		"-f webvtt pipe:1",
		opts.Start, opts.Duration+opts.Start, opts.Input, opts.Index)

	return pipe(args)
}
