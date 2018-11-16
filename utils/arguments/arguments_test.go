package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseArguments(t *testing.T) {
	var tests = []struct {
		input    string
		expected map[string]string
	}{
		{
			"-v=\"value\"",
			map[string]string{"v": "value"},
		},
		{
			"-a=\"another value\"",
			map[string]string{"a": "another value"},
		},
		{
			"-v=\"value\" -a=\"another value\"",
			map[string]string{"v": "value", "a": "another value"},
		},
		{
			"-v=value",
			map[string]string{},
		},
		{
			"-v='value'",
			map[string]string{},
		},
		{
			"",
			map[string]string{},
		},
	}

	for _, test := range tests {
		assert.Equal(t, ParseArguments(test.input), test.expected)
	}
}
