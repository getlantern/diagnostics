package diagchild

import (
	"fmt"
	"testing"

	"github.com/getlantern/diagnostics"
	"github.com/stretchr/testify/require"
)

// Debugging
// TODO: remove me
func TestRunAsChildProcess(t *testing.T) {
	b, err := JSONAsChildProcess(diagnostics.Config{
		PingConfig: &diagnostics.PingConfig{
			Addresses: []string{"8.8.8.8"},
		},
	}, "give me your password!")
	require.NoError(t, err)

	fmt.Println(string(b))
}
