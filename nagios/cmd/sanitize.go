package cmd

import (
	"strings"
)

var sanitizer = strings.NewReplacer(
	"\n", "\\n",
	";", ":",
)

// Sanitize should be used to ensure strings inserted into nagios commands
// are properly escaped.
func Sanitize(s string) string {
	return sanitizer.Replace(s)
}
