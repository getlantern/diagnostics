// Package diagchild provides support for running the diagnostics as a child process.
package diagchild

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
	
	"github.com/getlantern/diagnostics"
	"github.com/getlantern/diagnostics/internal/diagbin"
	"github.com/getlantern/errors"
)

var executablePath string

func init() {
	createTmpExecutable := func() (string, error) {
		tf, err := ioutil.TempFile("", "lantern-diagnostics")
		if err != nil {
			return "", err
		}
		return tf.Name(), tf.Chmod(0744)
	}

	path, err := createTmpExecutable()
	if err != nil {
		// This will cause the executable to be created in byteexec's default location.
		executablePath = "lantern-diagnostics"
	}
	executablePath = path
}

// Express the configuration as flags, in the form expected by the lantern-diagnostics command.
func configToFlags(cfg diagnostics.Config) []string {
	args := []string{}
	if cfg.PingConfig != nil {
		args = append(args, pingConfigToFlags(*cfg.PingConfig)...)
	}
	return args
}

func pingConfigToFlags(cfg diagnostics.PingConfig) []string {
	args := []string{}
	if len(cfg.Addresses) > 0 {
		args = append(args, "-ping-addresses", strings.Join(cfg.Addresses, ","))
	}
	if cfg.Count != 0 {
		args = append(args, "-ping-count", strconv.Itoa(cfg.Count))
	}
	return args
}

// JSONAsChildProcess is like RunAsChildProcess, but returns a JSON-encoded Report. This is a
// convenient shortcut provided as the child process already outputs the report in JSON.
func JSONAsChildProcess(cfg diagnostics.Config, prompt string) ([]byte, error) {
	return diagbin.RunFromBinary(executablePath, prompt, configToFlags(cfg)...)
}

// RunAsChildProcess is like Run, but runs the diagnostics in a separate process. This is useful for
// running the diagnostics with a different set of permissions than those assigned to the main
// process. The diagnostics process will prompt for elevated permissions if necessary.
func RunAsChildProcess(cfg diagnostics.Config, prompt string) (*diagnostics.Report, error) {
	b, err := JSONAsChildProcess(cfg, prompt)
	if err != nil {
		return nil, err
	}
	r := new(diagnostics.Report)
	if err := json.Unmarshal(b, r); err != nil {
		return nil, errors.New("failed to unmarshal report: %v", err)
	}
	return r, nil
}
