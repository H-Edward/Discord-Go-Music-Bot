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

// Sanitises a search query by removing unsafe characters and trimming length
// A false return value indicates the query is empty or invalid after sanitisation and shouldn't be used
func SanitiseSearchQuery(query string) (string, bool) {
	// Remove any non-alphanumeric characters except spaces
	var safeSearch = regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
	sanitised := safeSearch.ReplaceAllString(query, "")

	sanitised = regexp.MustCompile(`^\s+|\s+$`).ReplaceAllString(sanitised, "")

	sanitised = regexp.MustCompile(`\s+`).ReplaceAllString(sanitised, " ")

	good := IsValidSearchQuery(sanitised)
	if !good {
		return "", false
	}
	

	return sanitised, true
}
