package utils

import (
	"regexp"
	"strings"
)

var regexpNonAuthorizedChars = regexp.MustCompile("[^a-zA-Z0-9-_]")
var regexpMultipleDashes = regexp.MustCompile("-+")

func Slugify(s string) string {
	slug := regexpNonAuthorizedChars.ReplaceAllString(s, "-")
	return strings.Trim(regexpMultipleDashes.ReplaceAllString(slug, "-"), "-")
}
