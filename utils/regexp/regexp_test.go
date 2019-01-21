package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGroups(t *testing.T) {
	var tests = []struct {
		input    []string
		expected map[string]string
	}{
		{
			[]string{`(?P<testGroup>\d+)`, "bba110"},
			map[string]string{"testGroup": "110"},
		},
		{
			[]string{`(?P<testGroup1>\d+)\s(?P<testGroup2>[A-z]*)`, "bba110 foo"},
			map[string]string{"testGroup1": "110", "testGroup2": "foo"},
		},
	}

	for _, test := range tests {
		groups := GetGroups(test.input[0], test.input[1])
		assert.Equal(t, groups, test.expected)
	}
}
