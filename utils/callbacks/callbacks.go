package utils

import (
	"fmt"
	"regexp"
	"strconv"
)

const deleteTrackingCallbackRegex = `^\/delete_tracking (\d+)$`

func ParseDeleteTrackingCallback(callbackData string) (int64, error) {
	compiled := regexp.MustCompile(deleteTrackingCallbackRegex)
	match := compiled.FindAllStringSubmatch(callbackData, -1)

	if len(match) != 1 {
		return 0, fmt.Errorf(
			"Invalid callback data: %s. Required data's format: /delete_tracking ID.",
			callbackData,
		)
	}

	id, err := strconv.ParseInt(match[0][1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf(
			"Given ID %s is not a number.",
			match[0][1],
		)
	}

	return id, nil
}
