package ffmpeg

import (
	"fmt"
)

type CreatePosterOptions struct {
	Input string
	Time  float64
}

func CreatePoster(opts CreatePosterOptions) ([]byte, error) {
	args := fmt.Sprintf("-loglevel error -hide_banner -ss %.4f -i %s "+
		"-vframes 1 pipe:.jpg", opts.Time, opts.Input)

	return pipe(args)
}
