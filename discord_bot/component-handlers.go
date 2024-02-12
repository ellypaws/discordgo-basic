package discord_bot

import (
	"discordgo-basic/discord_bot/handlers"
	"github.com/bwmarrin/discordgo"
)

var componentHandlers = map[handlers.Component]func(bot *BotImpl, s *discordgo.Session, i *discordgo.InteractionCreate){
	handlers.DeleteButton: deleteMessage,
}

func deleteMessage(bot *BotImpl, s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.ChannelMessageDelete(i.ChannelID, i.Message.ID)
	if err != nil {
		handlers.ErrorEphemeralResponse(s, i.Interaction, err)
	}
}
