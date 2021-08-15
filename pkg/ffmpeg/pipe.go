package ffmpeg

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

// pipe executes an ffmpeg command and returns stdout and stderr
// It is expected that the ffmpeg output gets written to stdout
func pipe(args string) ([]byte, error) {
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
