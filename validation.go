package main

import (
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Checks if a search query is safe and valid
func isValidSearchQuery(query string) bool {
	var safeSearch = regexp.MustCompile(`^[a-zA-Z0-9\s]+$`)

	if !safeSearch.MatchString(query) {
		return false
	}

	if query == "" || len(query) > 200 {
		return false
	}
	return true
}

// Makes sure a URL entered is a valid YouTube URL
func isValidURL(input string) bool {
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

// Given a permission, checks if the user has that permission in the guild
func hasPermission(s *discordgo.Session, m *discordgo.MessageCreate, permission_requested int64) bool {

	member, err := s.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		return false
	}
	for _, role := range member.Roles {
		roleData, err := s.State.Role(m.GuildID, role)
		if err != nil {
			continue
		}
		if roleData.Permissions&permission_requested == permission_requested {
			return true
		}
		if roleData.Permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
			// If the user has the Administrator permission, they have all permissions
			return true
		}
	}

	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		return false
	}

	if guild.OwnerID == m.Author.ID {
		return true
	}

	// If no roles matched, check the member's permissions directly
	if member.Permissions&permission_requested == permission_requested {
		return true
	}
	return false

}
