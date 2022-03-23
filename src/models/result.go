package models

import (
	"net/url"
	"regexp"
	"strings"
)

var regexpNonAuthorizedChars = regexp.MustCompile("[^a-zA-Z0-9-_]")
var regexpMultipleDashes = regexp.MustCompile("-+")

type Result struct {
	Url url.URL
}

func (result *Result) ReplaceVars(cmd string) string {
	cmd = strings.ReplaceAll(cmd, "{result}", result.Url.String())
	cmd = strings.ReplaceAll(cmd, "{scheme}", result.Url.Scheme)
	cmd = strings.ReplaceAll(cmd, "{host}", result.Url.Host)
	cmd = strings.ReplaceAll(cmd, "{hostname}", result.Url.Hostname())
	cmd = strings.ReplaceAll(cmd, "{port}", result.Url.Port())
	cmd = strings.ReplaceAll(cmd, "{slug}", result.Slug())
	return cmd
}

func (result *Result) Slug() string {
	slug := regexpNonAuthorizedChars.ReplaceAllString(result.Url.String(), "-")
	return regexpMultipleDashes.ReplaceAllString(slug, "-")
}
