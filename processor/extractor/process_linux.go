package extractor

import (
	"io"
	"os/exec"
)

// Build a Process to execute & attach pipe to it here
func buildProcess(args string, command string) (*exec.Cmd, io.ReadCloser) {
	args = command + " " + args
	cmd := exec.Command("bash", "-c", args)

	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	handleErrors(err, "Metadata - StdoutPipe")

	return cmd, stdout
}
