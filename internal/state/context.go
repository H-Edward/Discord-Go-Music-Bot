package state

import (
	"strconv"
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
	SourceType   CommandSourceType // Where the command came from (i.e., interaction or message)
	Session      *discordgo.Session
	Interaction  *discordgo.InteractionCreate // Will be nil if not an interaction
	Message      *discordgo.MessageCreate     // Will be nil if not a message
	User         *discordgo.User              // Caller of the command
	GuildID      string                       // Guild ID where the command was called
	ChannelID    string                       // Channel ID where the command was called
	ArgumentsRaw map[string]interface{}       // Raw arguments from the command, type depends on source
	Arguments    map[string]string            // Standardised arguments, types are consistent
	CommandName  string                       // Name of the command being executed, used for determining argument keys
}

// Wrappers

func (ctx *Context) Reply(message string) {
	if ctx.SourceType == SourceTypeInteraction && ctx.Interaction != nil {
		ctx.Session.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message,
			},
		})
	}
	if ctx.SourceType == SourceTypeMessage && ctx.Message != nil {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, message)
	}

}

// Getters

// Naive getters
func (ctx *Context) GetSession() *discordgo.Session {
	return ctx.Session
}

func (ctx *Context) GetInteraction() *discordgo.InteractionCreate {
	return ctx.Interaction
}

func (ctx *Context) GetMessage() *discordgo.MessageCreate {
	return ctx.Message
}

// Lighter getters

func (ctx *Context) GetUser() *discordgo.User {
	return ctx.User
}

func (ctx *Context) GetGuildID() string {
	return ctx.GuildID
}

func (ctx *Context) GetChannelID() string {
	return ctx.ChannelID
}

func (ctx *Context) getArgument(key string) (interface{}, bool) {
	val, exists := ctx.Arguments[key]
	return val, exists
}
func (ctx *Context) getArgumentRaw(key string) (interface{}, bool) {
	val, exists := ctx.ArgumentsRaw[key]
	return val, exists
}

// Setters

func NewInteractionContext(s *discordgo.Session, i *discordgo.InteractionCreate) *Context {
	ctx := &Context{
		SourceType:   SourceTypeInteraction,
		Session:      s,
		Interaction:  i,
		User:         i.User,
		GuildID:      i.GuildID,
		ChannelID:    i.ChannelID,
		ArgumentsRaw: make(map[string]interface{}),
		Arguments:    make(map[string]string),
		CommandName:  i.ApplicationCommandData().Name,
	}

	if data := i.ApplicationCommandData(); data.Name != "" {
		for _, option := range data.Options {
			ctx.ArgumentsRaw[option.Name] = option.Value
		}
	}
	if ctx.User == nil && i.Member != nil {
		ctx.User = i.Member.User
	}

	ctx.standardiseArguments()
	return ctx
}

func NewMessageContext(s *discordgo.Session, m *discordgo.MessageCreate, command string) *Context {
	ctx := &Context{
		SourceType:   SourceTypeMessage,
		Session:      s,
		Message:      m,
		User:         m.Author,
		ChannelID:    m.ChannelID,
		GuildID:      m.GuildID,
		ArgumentsRaw: make(map[string]interface{}),
		Arguments:    make(map[string]string),
		CommandName:  command,
	}
	ctx.determineCommandNameFromMessage()

	ctx.determineArgumentsFromMessage()

	ctx.standardiseArguments()
	return ctx

}

func (ctx *Context) determineCommandNameFromMessage() {
	command := strings.Fields(ctx.GetMessage().Content)[0]
	if len(command) > 0 && command[0] == '!' {
		ctx.CommandName = command[1:]
		return
	}
	ctx.CommandName = ""
}

func (ctx *Context) determineArgumentsFromMessage() {
	// presume not sanitised
	switch ctx.CommandName {
	case "play":
		// everything after !play is the url
		if len(ctx.Message.Content) > 6 {
		} else {
			ctx.ArgumentsRaw["url"] = ""
		}
	case "search":
		if len(ctx.Message.Content) > 8 {
			ctx.ArgumentsRaw["query"] = ctx.Message.Content[8:]
		} else {
			ctx.ArgumentsRaw["query"] = ""
		}
	case "volume":
		if len(ctx.Message.Content) > 8 {
			ctx.ArgumentsRaw["level"] = ctx.Message.Content[8:]
		} else {
			ctx.ArgumentsRaw["level"] = ""
		}
	case "nuke":
		if len(ctx.Message.Content) > 6 {
			ctx.ArgumentsRaw["count"] = ctx.Message.Content[6:]
		} else {
			ctx.ArgumentsRaw["count"] = ""
		}
	default:
		// No arguments to parse for other commands
	}
}

// Convert raw arguments to standard types (not sanitised)
func (ctx *Context) standardiseArguments() {
	switch ctx.CommandName {
	case "play": // url string
		if val, exists := ctx.getArgumentRaw("url"); exists {
			if strVal, ok := val.(string); ok {
				ctx.Arguments["url"] = strVal
			} else {
				ctx.Arguments["url"] = ""
			}
		} else {
			ctx.Arguments["url"] = ""
		}
	case "search": // query string
		if val, exists := ctx.getArgumentRaw("query"); exists {
			if strVal, ok := val.(string); ok {
				ctx.Arguments["query"] = strVal
			} else {
				ctx.Arguments["query"] = ""
			}
		}

	case "volume": // level int (0-200)
		if val, exists := ctx.getArgumentRaw("level"); exists {
			switch v := val.(type) {
			case int:
				ctx.Arguments["level"] = strconv.Itoa(v)
			case float64:
				ctx.Arguments["level"] = strconv.Itoa(int(v))
			case string:
				ctx.Arguments["level"] = strings.TrimSpace(v)
			default:
				ctx.Arguments["level"] = ""
			}

		}
	case "nuke": // count int (1-100)
		if val, exists := ctx.getArgumentRaw("count"); exists {
			switch v := val.(type) {
			case int:
				ctx.Arguments["count"] = strconv.Itoa(v)
			case float64:
				ctx.Arguments["count"] = strconv.Itoa(int(v))
			case string:
				ctx.Arguments["count"] = strings.TrimSpace(v)
			default:
				ctx.Arguments["count"] = ""
			}

		}
	}

}
