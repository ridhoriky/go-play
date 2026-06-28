package utils

import (
	"regexp"
	"strings"
)

func GenerateSlug(name string) string {
	slug := strings.ToLower(name)

	reg := regexp.MustCompile("[^a-z0-9]+")
	slug = reg.ReplaceAllString(slug, "-")

	slug = strings.Trim(slug, "-")

	return slug
}
