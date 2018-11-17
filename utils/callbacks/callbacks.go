package utils

import (
	"fmt"
	"regexp"
	"strconv"
)

const deleteTrackingCallbackRegex = `^\/delete_tracking (\d+)$`

// ParseDeleteTrackingCallback parses delete tracking callback data
func ParseDeleteTrackingCallback(callbackData string) (int64, error) {
	compiled := regexp.MustCompile(deleteTrackingCallbackRegex)
	match := compiled.FindAllStringSubmatch(callbackData, -1)

	if len(match) != 1 {
		return 0, fmt.Errorf(
			"Invalid callback data: %s. Required data's format: /delete_tracking ID",
			callbackData,
		)
	}

	// regex parses only numbers
	id, _ := strconv.ParseInt(match[0][1], 10, 64)

	return id, nil
}
