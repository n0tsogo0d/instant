package ffmpeg

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
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
		"-c:a aac -b:a %d -ac 6 "+
		// output
		"-crf 26 -copyts -preset ultrafast pipe:.ts", opts.Start,
		opts.Start+opts.Duration, opts.Input, opts.Width, opts.Height,
		opts.VideoBitrate, opts.VideoBitrate, opts.VideoBitrate,
		opts.VideoBitrate/2, opts.AudioBitrate)

	cmd := exec.Command("ffmpeg", strings.Fields(args)...)

	var stdout, stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return nil, errors.New(stderr.String())
	}

	return stdout.Bytes(), nil
}

func CreateSubtitleSegment(opts CreateSubtitleSegmentOptions) ([]byte, error) {
	args := fmt.Sprintf("-loglevel error -hide_banner -i %s -map 0:%d "+
		" -f webvtt pipe:1", opts.Input, opts.Index)

	cmd := exec.Command("ffmpeg", strings.Fields(args)...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, errors.New(stderr.String())
	}

	return stdout.Bytes(), nil
}
