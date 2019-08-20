// Package diagnostics provides facilities for running tests and checks of a network or system.
package diagnostics

import (
	"errors"
	"runtime"
	"sync"
	"time"

	ping "github.com/sparrc/go-ping"
)

var (
	// ErrPingMissingAddress indicates that a ping test didn't have an address specified
	ErrPingMissingAddress = errors.New("ping missing address")

	// ErrPingUnsupportedPlatform indicates that a ping test was performed on a non-Windows platform and Force wasn't enabled
	ErrPingUnsupportedPlatform = errors.New("ping report is currently only supported on windows")
)

// Diagnostic is a diagnostic that can be run (such as a Ping)
type Diagnostic interface {
	// Type identifies the type of diagnostic
	Type() string

	// RunInSuite runs this diagnostic as part of as suite of diagnostics
	RunInSuite() (interface{}, error)
}

// DiagnosticResult is the result of running a diagnostic such as a Ping
type DiagnosticResult struct {
	Diagnostic string
	Result     interface{} `json:",omitempty"`
	Error      *string     `json:",omitempty"`
}

// Ping runs diagnostics using an ICMP ping utility.
type Ping struct {
	// Address specifies the address to ping
	Address string

	// Count is the number of packets sent per address. Defaults to 1.
	Count int `json:",omitempty"`

	// Force forces the ping report to run on non-Windows systems. Useful for testing, but requires
	// root permissions. See RunPingTest().
	Force bool `json:",omitempty"`
}

// Type implements the interface Diagnostic
func (p *Ping) Type() string {
	return "Ping"
}

// RunInSuite runs this Ping as part of as suite of diagnostics
func (p *Ping) RunInSuite() (interface{}, error) {
	return p.Run()
}

// PingResult is the result of running a Ping diagnostic
type PingResult struct {
	*Ping

	*ping.Statistics `json:",omitempty"`
}

// Run runs the Ping diagnostic
func (p *Ping) Run() (*PingResult, error) {
	if p.Address == "" {
		return nil, ErrPingMissingAddress
	}

	if runtime.GOOS != "windows" && !p.Force {
		// We need root permissions to ping on Linux and Mac OS:
		// https://github.com/sparrc/go-ping#note-on-windows-support
		//
		// We could just run the ping command on those systems and parse the output, but that
		// doesn't seem worth the effort at the moment.
		return nil, ErrPingUnsupportedPlatform
	}

	pinger, err := ping.NewPinger(p.Address)
	if err != nil {
		return nil, err
	}
	if p.Count <= 0 {
		pinger.Count = 1
	} else {
		pinger.Count = p.Count
	}

	// The default interval between packets is 1s. If we wait a bit longer than pinger.Count
	// seconds, we should receive all responses.
	pinger.Timeout = time.Second*time.Duration(pinger.Count) + time.Second

	pinger.SetPrivileged(true)
	pinger.Run()
	return &PingResult{Ping: p, Statistics: pinger.Statistics()}, nil
}

// Run runs multiple diagnostics, returning the corresponding results. It runs up to
// parallelism in parallel.
func Run(parallelism int, diagnostics ...Diagnostic) []*DiagnosticResult {
	var mx sync.Mutex
	requests := make(chan *diagnosticRequest, len(diagnostics))
	results := make([]*DiagnosticResult, len(diagnostics))
	for i, diagnostic := range diagnostics {
		requests <- &diagnosticRequest{diagnostic, i}
	}
	close(requests)
	if parallelism < 1 {
		parallelism = 1
	}
	var wg sync.WaitGroup
	wg.Add(parallelism)
	for i := 0; i < parallelism; i++ {
		go func() {
			defer wg.Done()
			for request := range requests {
				result, err := request.RunInSuite()
				diagnosticResult := &DiagnosticResult{
					Diagnostic: request.Type(),
					Result:     result,
					Error:      sPtr(err),
				}
				mx.Lock()
				results[request.order] = diagnosticResult
				mx.Unlock()
			}
		}()
	}

	wg.Wait()
	return results
}

type diagnosticRequest struct {
	Diagnostic

	order int
}

func sPtr(err error) *string {
	if err == nil {
		return nil
	}
	s := err.Error()
	return &s
}
