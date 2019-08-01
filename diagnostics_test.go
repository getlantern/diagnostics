package diagnostics

import (
	"os"
	"testing"

	"github.com/go-yaml/yaml"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	results := Run(2, &Ping{Address: "8.8.8.8", Count: 1, Force: true}, &Ping{Address: "999.999.999.999", Count: 1, Force: true})
	err := yaml.NewEncoder(os.Stdout).Encode(results)
	assert.NoError(t, err)
}
