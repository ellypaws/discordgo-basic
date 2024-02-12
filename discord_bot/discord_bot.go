package discord_bot

import (
	"discordgo-basic/discord_bot/handlers"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"strings"

	"github.com/charmbracelet/log"
)

type BotImpl struct {
	botSession         *discordgo.Session
	guildID            string
	registeredCommands map[Command]*discordgo.ApplicationCommand
	imagineCommand     *Command
	config             *Config
}

type Config struct {
	BotToken       string
	GuildID        string
	RemoveCommands bool
}

func New(cfg *Config) (*BotImpl, error) {
	if cfg.BotToken == "" {
		return nil, errors.New("missing bot token")
	}

	handlers.Token = &cfg.BotToken

	if cfg.GuildID == "" {
		//return nil, errors.New("missing guild ID")
		log.Printf("Guild ID not provided, commands will be registered globally")
	}

	botSession, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		return nil, err
	}

	botSession.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err = botSession.Open()
	if err != nil {
		return nil, err
	}

	bot := &BotImpl{
		botSession:         botSession,
		registeredCommands: make(map[Command]*discordgo.ApplicationCommand),
		config:             cfg,
	}

	err = bot.registerCommands()
	if err != nil {
		return nil, err
	}

	bot.registerHandlers(botSession)

	return bot, nil
}

func (b *BotImpl) registerHandlers(session *discordgo.Session) {
	session.AddHandler(func(session *discordgo.Session, i *discordgo.InteractionCreate) {
		var h func(b *BotImpl, s *discordgo.Session, i *discordgo.InteractionCreate)
		var ok bool
		switch i.Type {
		// commands
		case discordgo.InteractionApplicationCommand:
			h, ok = commandHandlers[Command(i.ApplicationCommandData().Name)]
		// buttons
		case discordgo.InteractionMessageComponent:
			log.Printf("Component with customID `%v` was pressed, attempting to respond\n", i.MessageComponentData().CustomID)
			h, ok = componentHandlers[handlers.Component(i.MessageComponentData().CustomID)]
		// autocomplete
		case discordgo.InteractionApplicationCommandAutocomplete:
			h, ok = autocompleteHandlers[Command(i.ApplicationCommandData().Name)]
		// modals
		case discordgo.InteractionModalSubmit:
			h, ok = modalHandlers[Command(i.ModalSubmitData().CustomID)]
		default:
			log.Printf("Unknown interaction type '%v'", i.Type)
		}

		if !ok || h == nil {
			var interactionType string = "unknown"
			var interactionName string = "unknown"
			switch i.Type {
			case discordgo.InteractionApplicationCommand:
				interactionType = "command"
				interactionName = i.ApplicationCommandData().Name
			case discordgo.InteractionMessageComponent:
				interactionType = "component"
				interactionName = i.MessageComponentData().CustomID
			case discordgo.InteractionApplicationCommandAutocomplete:
				interactionType = "autocomplete"

				data := i.ApplicationCommandData()
				for _, opt := range data.Options {
					if !opt.Focused {
						continue
					}
					interactionName = fmt.Sprintf("command: /%v option: %v (%v)", data.Name, opt.Name)
					break
				}
			case discordgo.InteractionModalSubmit:
				interactionType = "modal"
				interactionName = i.ModalSubmitData().CustomID
			}
			log.Printf("WARNING: Cannot find handler for interaction [%v] '%v'", interactionType, interactionName)
			return
		}

		h(b, session, i)
	})

	log.Debugf("Registered handlers %v", commandHandlers)
}

func (b *BotImpl) registerCommands() error {
	b.registeredCommands = make(map[Command]*discordgo.ApplicationCommand, len(commands))
	for key, command := range commands {
		if command.Name == "" {
			// clean the key because it might be a description of some sort
			// only get the first word, and clean to only alphanumeric characters or -
			sanitized := strings.ReplaceAll(string(key), " ", "-")
			sanitized = strings.ToLower(sanitized)

			// remove all non-valid characters
			for _, c := range sanitized {
				if (c < 'a' || c > 'z') && (c < '0' || c > '9') && c != '-' {
					sanitized = strings.ReplaceAll(sanitized, string(c), "")
				}
			}
			command.Name = sanitized
		}

		cmd, err := b.botSession.ApplicationCommandCreate(b.botSession.State.User.ID, b.guildID, command)
		if err != nil {
			return errors.New(fmt.Sprintf("Cannot create '%v' command: %v", command.Name, err))
		}
		b.registeredCommands[key] = cmd

		log.Debugf("Registered %v command as: /%v", key, cmd.Name)
	}

	return nil
}

func (b *BotImpl) rebuildMap(
	f func(*BotImpl) Command,
	key *Command,
	m map[Command]*discordgo.ApplicationCommand,
	h map[Command]func(b *BotImpl, s *discordgo.Session, i *discordgo.InteractionCreate,
	)) {
	oldKey := *key

	*key = f(b)
	if *key == oldKey {
		return
	}
	log.Printf("Rebuilding map for '%v' to '%v'", oldKey, *key)

	m[*key] = m[oldKey]
	m[*key].Name = string(*key)
	h[*key] = h[oldKey]
	delete(m, oldKey)
	delete(h, oldKey)
}

func (b *BotImpl) Start() {
	StartPolling()

	err := b.teardown()
	if err != nil {
		log.Printf("Error tearing down bot: %v", err)
	}
}

func StartPolling() {
	log.Print("Press Ctrl+C to exit")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

Polling:
	for {
		select {
		case <-stop:
			break Polling
		}
	}
	log.Printf("Polling stopped...\n")
}

func (b *BotImpl) teardown() error {
	// Delete all commands added by the bot
	if b.config.RemoveCommands {
		log.Printf("Removing all commands added by bot...")

		for key, v := range b.registeredCommands {
			log.Printf("Removing command [key:%v], '%v'...", key, v.Name)

			err := b.botSession.ApplicationCommandDelete(b.botSession.State.User.ID, b.guildID, v.ID)
			if err != nil {
				log.Fatalf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	return b.botSession.Close()
}

func shortenString(s string) string {
	if len(s) > 90 {
		log.Debugf("Shortened string to: %v", s[:90])
		return s[:90]
	}
	return s
}
