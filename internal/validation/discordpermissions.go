package validation

import (
	"discord-go-music-bot/internal/state"

	"github.com/bwmarrin/discordgo"
)

// Given a permission, checks if the user has that permission in the guild
func HasPermission(ctx state.Context, permission_requested int64) bool {

	member, err := ctx.GetSession().GuildMember(ctx.GetGuildID(), ctx.GetUser().ID)
	if err != nil {
		return false
	}
	for _, role := range member.Roles {
		roleData, err := ctx.GetSession().State.Role(ctx.GetGuildID(), role)
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

	guild, err := ctx.GetSession().State.Guild(ctx.GetGuildID())
	if err != nil {
		return false
	}

	if guild.OwnerID == ctx.GetUser().ID {
		return true
	}

	// If no roles matched, check the member's permissions directly
	if member.Permissions&permission_requested == permission_requested {
		return true
	}
	return false

}
