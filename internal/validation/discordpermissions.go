package validation

import "github.com/bwmarrin/discordgo"

// Given a permission, checks if the user has that permission in the guild
func HasPermission(s *discordgo.Session, m *discordgo.MessageCreate, permission_requested int64) bool {

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
