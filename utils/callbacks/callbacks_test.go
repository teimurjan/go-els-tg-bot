package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDeleteTrackingCallback(t *testing.T) {
	type returnValue struct {
		id int64
		e  error
	}
	var tests = []struct {
		input    string
		expected *returnValue
	}{
		{
			"/delete_tracking 1",
			&returnValue{1, nil},
		},
		{
			"/delete_tracking 757",
			&returnValue{757, nil},
		},
		{
			"/delete_tracking",
			&returnValue{
				id: 0,
				e: fmt.Errorf(
					"Invalid callback data: %s. Required data's format: /delete_tracking ID",
					"/delete_tracking",
				),
			},
		},
		{
			"/delete_tracking",
			&returnValue{
				id: 0,
				e: fmt.Errorf(
					"Invalid callback data: %s. Required data's format: /delete_tracking ID",
					"/delete_tracking",
				),
			},
		},
	}

	for _, test := range tests {
		id, err := ParseDeleteTrackingCallback(test.input)
		assert.Equal(t, err, test.expected.e)
		assert.Equal(t, id, test.expected.id)
	}
}
