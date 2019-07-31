// Package diagbin facilitates running the diagnostics as a pre-compiled binary.
package diagbin

import (
	"bytes"
	"os/exec"

	"github.com/getlantern/byteexec"
	"github.com/getlantern/elevate"
	"github.com/getlantern/errors"
)

const assetName = "diagbin"

// RunFromBinary runs the diagnostics command from a compiled binary.
func RunFromBinary(path, prompt string, arguments ...string) ([]byte, error) {
	asset, err := Asset(assetName)
	if err != nil {
		return nil, errors.New("could not find asset: %v", err)
	}
	executable, err := byteexec.New(asset, path)
	if err != nil {
		return nil, errors.New("failed to create executable from asset: %v", err)
	}

	cmd := elevate.WithPrompt(prompt).Command(executable.Filename, arguments...)
	cmd.Stdout, cmd.Stderr = new(bytes.Buffer), new(bytes.Buffer)
	if err := cmd.Run(); err != nil && exitCode(err) != 2 {
		if stderr := cmd.Stderr.(*bytes.Buffer).Bytes(); len(stderr) > 0 {
			return nil, errors.New("failed to run executable: %v: %s", err, string(stderr))
		}
		return nil, errors.New("failed to run executable: %v", err)
	}
	return cmd.Stdout.(*bytes.Buffer).Bytes(), nil
}

func exitCode(err error) int {
	// TODO: test on Windows and Linux
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode()
	}
	return -1
}
