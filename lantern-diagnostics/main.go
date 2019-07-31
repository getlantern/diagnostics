// Command lantern-diagnostics implements a diagnostics tool. This tool is intended for use by the
// flashlight client. The main use case is running diagnostics and attaching a report when a user
// submits an issue through the application.
//
// Only the report is written to stdout. All logs and error messages are written to stderr.
//
// Exits with status 1 if diagnostics could not be run. Exits 2 if there was a partial failure in
// running diagnostics. In the latter case, the report will contain partial results as well as
// details on failures.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/getlantern/diagnostics"
)

var (
	indent        = flag.Bool("indent", false, "print the report with newlines and indentations")
	pingAddresses = flag.String("ping-addresses", "", "comma-separated list of addresses for ping tests")
	pingCount     = flag.Int("ping-count", 1, "number of ping packets to send to each address")
)

func main() {
	flag.Parse()

	cfg := diagnostics.Config{}
	if *pingAddresses != "" {
		cfg.PingConfig = &diagnostics.PingConfig{
			Addresses: strings.Split(*pingAddresses, ","),
			Count:     *pingCount,
		}
	}
	report := diagnostics.Run(cfg)

	var marshal func(v interface{}) ([]byte, error)
	if *indent {
		marshal = func(v interface{}) ([]byte, error) { return json.MarshalIndent(v, "", "  ") }
	} else {
		marshal = json.Marshal
	}

	reportJSON, err := marshal(report)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to marshal report:", err)
		os.Exit(1)
	}
	fmt.Println(string(reportJSON))

	if report.HasErrors() {
		// Exit 3 when the report has errors (exit code 2 is used for flag parsing errors).
		os.Exit(3)
	}
}
