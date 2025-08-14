package validation

import "regexp"

// Checks if a search query is safe and valid
func IsValidSearchQuery(query string) bool {
	var safeSearch = regexp.MustCompile(`^[a-zA-Z0-9\s]+$`)

	if !safeSearch.MatchString(query) {
		return false
	}

	if query == "" || len(query) > 200 {
		return false
	}
	return true
}
