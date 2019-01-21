package utils

import "regexp"

// GetGroups gets regexp groups and insert them to a map
func GetGroups(regexpStr, str string) map[string]string {
	compiledRegexp := regexp.MustCompile(regexpStr)
	match := compiledRegexp.FindStringSubmatch(str)

	groups := make(map[string]string)
	for i, name := range compiledRegexp.SubexpNames() {
		if i > 0 && i <= len(match) {
			groups[name] = match[i]
		}
	}
	return groups
}
