package main

import (
	"discordgo-basic/discord_bot"
	"flag"
	"github.com/joho/godotenv"
	"log"
	"os"
)

// Bot parameters
var (
	guildID            = flag.String("guild", "", "Guild ID. If not passed - bot registers commands globally")
	botToken           = flag.String("token", "", "Bot access token")
	removeCommandsFlag = flag.Bool("remove", false, "Delete all commands when bot exits")
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
	log.Println(".env file loaded successfully")

	if botToken == nil || *botToken == "" {
		tokenEnv := os.Getenv("BOT_TOKEN")
		if tokenEnv == "YOUR_BOT_TOKEN_HERE" {
			log.Fatalf("Invalid bot token: %v\n"+
				"Did you edit the .env or run the program with -token ?", tokenEnv)
		}
		if tokenEnv != "" {
			botToken = &tokenEnv
		}
	}

	if guildID == nil || *guildID == "" {
		guildEnv := os.Getenv("GUILD_ID")
		if guildEnv != "" {
			guildID = &guildEnv
		}
	}

	if removeCommandsFlag == nil || !*removeCommandsFlag {
		removeCommandsEnv := os.Getenv("REMOVE_COMMANDS")
		if removeCommandsEnv != "" {
			removeCommandsFlag = new(bool)
			*removeCommandsFlag = removeCommandsEnv == "true"
		}
	}
}

func main() {
	flag.Parse()

	if botToken == nil || *botToken == "" {
		log.Fatalf("Bot token flag is required")
	}

	var removeCommands bool

	if removeCommandsFlag != nil && *removeCommandsFlag {
		removeCommands = *removeCommandsFlag
	}

	bot, err := discord_bot.New(&discord_bot.Config{
		BotToken:       *botToken,
		GuildID:        *guildID,
		RemoveCommands: removeCommands,
	})
	if err != nil {
		log.Fatalf("Error creating Discord bot: %v", err)
	}

	bot.Start()

	log.Println("Gracefully shutting down.")
}
