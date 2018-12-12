package utils

import (
	"regexp"
)

const argumentsRegex = "-(?P<flag>[a-z])=[\"“”](?P<value>.{1,}?)[\"“”]"

func ParseArguments(arguments string) map[string]string {
	compiled := regexp.MustCompile(argumentsRegex)
	match := compiled.FindAllStringSubmatch(arguments, -1)

	parsedArguments := make(map[string]string)
	if len(match) == 0 {
		return parsedArguments
	}

	for _, matchItem := range match {
		parsedArguments[matchItem[1]] = matchItem[2]
	}

	return parsedArguments
}
