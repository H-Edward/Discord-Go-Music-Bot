package state

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type CommandSourceType int

const (
	SourceTypeUnknown     CommandSourceType = iota
	SourceTypeInteraction                   // Slash commands
	SourceTypeMessage                       // Text commands
)

type Context struct {
	SourceType  CommandSourceType // Where the command came from (i.e., interaction or message)
	Session     *discordgo.Session
	Interaction *discordgo.InteractionCreate // Will be nil if not an interaction
	Message     *discordgo.MessageCreate     // Will be nil if not a message
	User        *discordgo.User              // Caller of the command
	GuildID     string                       // Guild ID where the command was called
	ChannelID   string                       // Channel ID where the command was called
	Arguments   map[string]interface{}
	CommandName string // Name of the command being executed, used for determining argument keys
}

// Wrappers

func (c *Context) Reply(message string) {
	if c.SourceType == SourceTypeInteraction && c.Interaction != nil {
		c.Session.InteractionRespond(c.Interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message,
			},
		})
	}
	if c.SourceType == SourceTypeMessage && c.Message != nil {
		c.Session.ChannelMessageSend(c.Message.ChannelID, message)
	}

}

// Getters

// Naive getters
func (c *Context) GetSession() *discordgo.Session {
	return c.Session
}

func (c *Context) GetInteraction() *discordgo.InteractionCreate {
	return c.Interaction
}

func (c *Context) GetMessage() *discordgo.MessageCreate {
	return c.Message
}

// Lighter getters

func (c *Context) GetUser() *discordgo.User {
	return c.User
}

func (c *Context) GetGuildID() string {
	return c.GuildID
}

func (c *Context) GetChannelID() string {
	return c.ChannelID
}

func (c *Context) getArgument(key string) (interface{}, bool) {
	val, exists := c.Arguments[key]
	return val, exists
}

// Setters

func NewInteractionContext(s *discordgo.Session, i *discordgo.InteractionCreate) *Context {
	ctx := &Context{
		SourceType:  SourceTypeInteraction,
		Session:     s,
		Interaction: i,
		User:        i.Member.User,
		GuildID:     i.GuildID,
		ChannelID:   i.Message.ChannelID,
		Arguments:   make(map[string]interface{}),
		CommandName: i.ApplicationCommandData().Name,
	}

	if data := i.ApplicationCommandData(); data.Name != "" {
		for _, option := range data.Options {
			ctx.Arguments[option.Name] = option.Value
		}
	}
	return ctx
}

func NewMessageContext(s *discordgo.Session, m *discordgo.MessageCreate, command string) *Context {
	ctx := &Context{
		SourceType:  SourceTypeMessage,
		Session:     s,
		Message:     m,
		User:        m.Author,
		ChannelID:   m.ChannelID,
		GuildID:     m.GuildID,
		Arguments:   make(map[string]interface{}),
		CommandName: command,
	}
	determineCommandNameFromMessage(ctx)
	// Determine arguments from message content if needed
	determineArgumentsFromMessage(ctx)
	return ctx

}

func determineCommandNameFromMessage(ctx *Context) string {
	command := strings.Fields(ctx.GetMessage().Content)[0]
	if len(command) > 0 && command[0] == '!' {
		return command[1:] // Remove the '!' prefix
	}
	return ""
}

func determineArgumentsFromMessage(ctx *Context) {
	// These aren't to be considered sanitised or validated, just parsed out of the message content
	switch ctx.CommandName {
	case "play":
		// everything after !play is the url
		if len(ctx.Message.Content) > 6 {
			ctx.Arguments["url"] = ctx.Message.Content[6:]
		} else {
			ctx.Arguments["url"] = ""
		}

	case "search":
		if len(ctx.Message.Content) > 8 {
			ctx.Arguments["query"] = ctx.Message.Content[8:]
		} else {
			ctx.Arguments["query"] = ""
		}
	case "volume":
		if len(ctx.Message.Content) > 8 {
			ctx.Arguments["level"] = ctx.Message.Content[8:]
		} else {
			ctx.Arguments["level"] = ""
		}
	case "nuke":
		if len(ctx.Message.Content) > 6 {
			ctx.Arguments["count"] = ctx.Message.Content[6:]
		} else {
			ctx.Arguments["count"] = ""
		}
	default:
		// No arguments to parse for other commands
	}
}
