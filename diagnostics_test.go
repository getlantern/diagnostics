package diagnostics

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	S     *testStruct            `json:",omitempty"`
	M     map[string]*testStruct `json:",omitempty"`
	Error *string                `json:",omitempty"`
}

func (ts testStruct) String() string {
	b, err := json.MarshalIndent(ts, "", "  ")
	if err != nil {
		return fmt.Sprintf("%#v", ts)
	}
	return string(b)
}

func TestHasErrors(t *testing.T) {
	for _, testCase := range []struct {
		input    interface{}
		expected bool
	}{
		{
			input: testStruct{
				Error: sPtr(""),
			},
			expected: true,
		},
		{
			input: testStruct{
				S: &testStruct{
					Error: sPtr(""),
				},
			},
			expected: true,
		},
		{
			input: testStruct{
				M: map[string]*testStruct{
					"": &testStruct{
						Error: sPtr(""),
					},
				},
			},
			expected: true,
		},
		{
			input:    testStruct{},
			expected: false,
		},
		{
			input: testStruct{
				S: &testStruct{},
			},
			expected: false,
		},
		{
			input: testStruct{
				M: map[string]*testStruct{
					"": &testStruct{},
				},
			},
			expected: false,
		},
	} {
		assert.NotPanics(t, func() {
			assert.Equal(t, testCase.expected, hasErrors(testCase.input), "input:\n%v", testCase.input)
		}, "input:\n%v", testCase.input)
	}
}
