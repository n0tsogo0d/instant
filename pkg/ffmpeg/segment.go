package ffmpeg

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// CreateVideoSegmentOptions for ffmpeg args on video segment transcoding
type CreateVideoSegmentOptions struct {
	Input    string
	Start    float64
	Duration float64
	Bitrate  int
	Height   int
	Width    int
}

// CreateVideoSegmentOptions for ffmpeg args on audio segment transcoding
type CreateAudioSegmentOptions struct {
	Input    string
	Start    float64
	Duration float64
	Bitrate  int
}

// func CreateAudioSegment(opts CreateAudioSegmentOptions) ([]byte, error) {
// 	args := fmt.Sprintf("-loglevel error -hide_banner"+
// 		" -ss %.4f -t %.4f -i %s "+
// 		"-vn -c:a aac -b:a %d -ac 2 pipe:.aac", opts.Start,
// 		opts.Duration, opts.Input, opts.Bitrate)

// 	cmd := exec.Command("ffmpeg", strings.Fields(args)...)
// 	var stdout, stderr bytes.Buffer
// 	cmd.Stderr = &stderr
// 	cmd.Stdout = &stdout
// 	err := cmd.Run()
// 	if err != nil {
// 		return nil, errors.New(stderr.String())
// 	}

// 	return stdout.Bytes(), nil
// }

func CreateAudioSegment(opts CreateAudioSegmentOptions) ([]byte, error) {
	// args := fmt.Sprintf("-loglevel error -hide_banner "+
	// 	"-ss %.4f -to %.4f -i %s -sn -vn -muxdelay 0 -muxpreload 0 "+
	// 	"-c:a aac -b:a %d -ac 2 "+
	// 	"-copyts -f hls pipe:1", opts.Start,
	// 	opts.Start+opts.Duration, opts.Input, opts.Bitrate)

	args := fmt.Sprintf("-loglevel error -hide_banner "+
		" -ss %.4f -t %.4f -i %s "+
		"-vn -muxdelay 0 -muxpreload 0 -c:a aac -b:a %d -ac 2 pipe:.aac", opts.Start,
		opts.Duration, opts.Input, opts.Bitrate)

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

func CreateVideoSegment(opts CreateVideoSegmentOptions) ([]byte, error) {
	// args := fmt.Sprintf("-loglevel error -hide_banner "+
	// 	"-ss %.4f -to %.4f -i %s -sn -an -muxdelay 0 -muxpreload 0 "+
	// 	"-vf scale=w=%d:h=%d:force_original_aspect_ratio=decrease "+
	// 	"-b:v %d -minrate %d -maxrate %d -bufsize %d -profile:v main "+
	// 	"-copyts -f hls -preset ultrafast pipe:1", opts.Start,
	// 	opts.Start+opts.Duration, opts.Input, opts.Width, opts.Height,
	// 	opts.Bitrate, opts.Bitrate, opts.Bitrate, opts.Bitrate/2)

	args := fmt.Sprintf("-loglevel error -hide_banner "+
		"-ss %.4f -to %.4f -i %s -sn -an -muxdelay 0 -muxpreload 0 "+
		"-vf scale=w=%d:h=%d:force_original_aspect_ratio=decrease "+
		"-b:v %d -minrate %d -maxrate %d -bufsize %d -profile:v main "+
		"-copyts -f hls -preset ultrafast pipe:1", opts.Start,
		opts.Start+opts.Duration, opts.Input, opts.Width, opts.Height,
		opts.Bitrate, opts.Bitrate, opts.Bitrate, opts.Bitrate/2)

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

type CreateSubtitleSegmentOptions struct {
	Input    string
	Start    float64
	Duration float64
	Index    int
}

// CreateSubtitleSegment creates WEBVTT subtitles for a given subtitle stream
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
