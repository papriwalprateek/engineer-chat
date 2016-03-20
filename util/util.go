package util

import "regexp"

// ParseMsg parses message contents (/{action} {message})
func ParseMsg(message string) (string, string) {
	actionRegex, _ := regexp.Compile(`^\/([^\s]*)\s*(.*)$`)
	res := actionRegex.FindAllStringSubmatch(message, -1)
	if len(res) == 1 {
		return res[0][1], res[0][2]
	}
	return "", ""
}
