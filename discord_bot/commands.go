package discord_bot

import (
	"github.com/bwmarrin/discordgo"
)

type (
	Command       string
	CommandOption string
)

var (
	// Command names
	imagineCommand         Command = "imagine"
	imagineSettingsCommand Command = "imagine_settings"
)

const (
	helloCommand Command = "hello"
)

const (
	// Command options
	promptOption CommandOption = "prompt"
)

var commands = map[Command]*discordgo.ApplicationCommand{
	helloCommand: {
		Name: string(helloCommand),
		// All commands and options must have a description
		// Commands/options without description will fail the registration
		// of the command.
		Description: "Say hello to the bot",
		Type:        discordgo.ChatApplicationCommand,
	},
}

var commandOptions = map[CommandOption]*discordgo.ApplicationCommandOption{
	promptOption: {
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        string(promptOption),
		Description: "The text prompt to imagine",
		Required:    true,
	},
}

const (
	maskedUser    = "user"
	maskedChannel = "channel"
	maskedForum   = "threads"
	maskedRole    = "role"
)

var maskedOptions = map[string]*discordgo.ApplicationCommandOption{
	maskedUser: {
		Type:        discordgo.ApplicationCommandOptionUser,
		Name:        maskedUser,
		Description: "Choose a user",
		Required:    false,
	},
	maskedChannel: {
		Type:        discordgo.ApplicationCommandOptionChannel,
		Name:        maskedChannel,
		Description: "Choose a channel to close",
		// Channel type mask
		ChannelTypes: []discordgo.ChannelType{
			discordgo.ChannelTypeGuildText,
			discordgo.ChannelTypeGuildVoice,
		},
		Required: false,
	},
	maskedForum: {
		Type:        discordgo.ApplicationCommandOptionChannel,
		Name:        maskedForum,
		Description: "Choose a thread to mark as solved",
		ChannelTypes: []discordgo.ChannelType{
			discordgo.ChannelTypeGuildForum,
			discordgo.ChannelTypeGuildNewsThread,
			discordgo.ChannelTypeGuildPublicThread,
			discordgo.ChannelTypeGuildPrivateThread,
		},
	},
	maskedRole: {
		Type:        discordgo.ApplicationCommandOptionRole,
		Name:        maskedRole,
		Description: "Choose a role to add",
		Required:    false,
	},
}
