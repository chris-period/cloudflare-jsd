package cloudflare

import (
	"regexp"
	"strings"
)

var lzKeyRegex = regexp.MustCompile(`[^\s,]*\$[^\s,]*\+?[^\s,]*`)
var sKeyRegex = regexp.MustCompile(`\d+\.\d+:\d+:[^\s,]+`)

func parse(content string) (string, string) {
	foundLZKey := lzKeyRegex.FindAllString(content, -1)
	foundSKey := sKeyRegex.FindAllString(content, -1)

	if len(foundLZKey) == 2 && len(foundSKey) == 2 {
		lzKey := foundLZKey[1]
		if !strings.Contains(lzKey, "+") {
			lzKey = foundLZKey[0]
		}
		return lzKey, foundSKey[1]
	}

	return "", ""
}
