package validation

import (
	"regexp"
	"strings"
)

// Makes sure a URL entered is a valid YouTube URL
func IsValidURL(input string) bool {
	input = strings.TrimSpace(input)

	const prefix = `^(?:https?://)?(?:www\.)?`

	// watch URL: youtube.com/watch?v=VIDEO_ID
	watchRe := regexp.MustCompile(prefix + `(?:youtube\.com|m\.youtube\.com)/watch\?v=[\w-]{11}(?:&[\w=&-]*)?$`)

	// shorts URL: youtube.com/shorts/VIDEO_ID
	shortsRe := regexp.MustCompile(prefix + `(?:youtube\.com|m\.youtube\.com)/shorts/[\w-]{11}$`)

	// embed URL: youtube.com/embed/VIDEO_ID
	embedRe := regexp.MustCompile(prefix + `(?:youtube\.com|m\.youtube\.com)/embed/[\w-]{11}$`)

	// short URL: youtu.be/VIDEO_ID
	shortURLRe := regexp.MustCompile(prefix + `youtu\.be/[\w-]{11}$`)

	regexes := []*regexp.Regexp{watchRe, shortsRe, embedRe, shortURLRe}

	for _, re := range regexes {
		if re.MatchString(input) {
			return true
		}
	}
	return false
}
