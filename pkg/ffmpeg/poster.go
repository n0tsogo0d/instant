package ffmpeg

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type CreatePosterOptions struct {
	Input string
	Time  float64
}

func CreatePoster(opts CreatePosterOptions) ([]byte, error) {
	args := fmt.Sprintf("-loglevel error -hide_banner -ss %.4f -i %s "+
		"-vframes 1 pipe:.jpg", opts.Time, opts.Input)

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
