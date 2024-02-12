package discord_bot

import (
	"cmp"
	"discordgo-basic/discord_bot/handlers"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"regexp"
)

var commandHandlers = map[Command]func(b *BotImpl, s *discordgo.Session, i *discordgo.InteractionCreate){
	helloCommand: func(b *BotImpl, bot *discordgo.Session, i *discordgo.InteractionCreate) {
		handlers.Responses[handlers.HelloResponse].(handlers.NewResponseType)(bot, i)
	},
}

var autocompleteHandlers = map[Command]func(b *BotImpl, s *discordgo.Session, i *discordgo.InteractionCreate){}

var modalHandlers = map[Command]func(b *BotImpl, s *discordgo.Session, i *discordgo.InteractionCreate){}

func getOpts(data discordgo.ApplicationCommandInteractionData) map[CommandOption]*discordgo.ApplicationCommandInteractionDataOption {
	options := data.Options
	optionMap := make(map[CommandOption]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[CommandOption(opt.Name)] = opt
	}
	return optionMap
}

// If FieldType and ValueType are the same, then we attempt to assert FieldType to value.Value
// Otherwise, we return the interface conversion to the caller to do manual type conversion
//
// Example:
//
//	if int64Val, ok := interfaceConvertAuto[int, int64](&queue.Steps, stepOption, optionMap, parameters); ok {
//		queue.Steps = int(*int64Val)
//	}
//
// (*discordgo.ApplicationCommandInteractionDataOption).IntValue() actually uses float64 for the interface conversion, so use float64 for integers, numbers, etc.
// and then convert to the desired type.
// Only string and float64 are supported for V as that's what the discordgo API returns.
// If the field is nil, then we don't assign the value to the field.
// Instead, we return *V and bool to indicate whether the conversion was successful.
// This is useful for when we want to convert to a type that is not the same as the field type.
func interfaceConvertAuto[F any, V string | float64](field *F, option CommandOption, optionMap map[CommandOption]*discordgo.ApplicationCommandInteractionDataOption, parameters map[CommandOption]string) (*V, bool) {
	if value, ok := optionMap[option]; ok {
		vToField, ok := value.Value.(F)
		if ok && field != nil {
			*field = vToField
		}
		valueType, ok := value.Value.(V)
		return &valueType, ok
	}
	if value, ok := parameters[option]; ok {
		if field != nil {
			_, err := fmt.Sscanf(value, "%v", field)
			if err != nil {
				return nil, false
			}
		}
		var out V
		_, err := fmt.Sscanf(value, "%v", &out)
		if err != nil {
			return nil, false
		}
		return &out, true
	}
	return nil, false
}

func between[T cmp.Ordered](value, minimum, maximum T) T {
	return min(max(minimum, value), maximum)
}

func sanitizeTooltip(input string) string {
	tooltipRegex := regexp.MustCompile(`[✨❌](.+) 🪄:([\d.]+)$|[✨❌](.+)`)
	sanitizedTooltip := tooltipRegex.FindStringSubmatch(input)

	if sanitizedTooltip != nil {
		log.Printf("Removing tooltip: %#v", sanitizedTooltip)

		switch {
		case sanitizedTooltip[1] != "":
			input = sanitizedTooltip[1] + ":" + sanitizedTooltip[2]
		case sanitizedTooltip[3] != "":
			input = sanitizedTooltip[3]
		}
		log.Printf("Sanitized input: %v", input)
	}
	return input
}

//	type TextInput struct {
//	   CustomID    string         `json:"custom_id"`
//	   Label       string         `json:"label"`
//	   Style       TextInputStyle `json:"style"`
//	   Placeholder string         `json:"placeholder,omitempty"`
//	   Value       string         `json:"value,omitempty"`
//	   Required    bool           `json:"required"`
//	   MinLength   int            `json:"min_length,omitempty"`
//	   MaxLength   int            `json:"max_length,omitempty"`
//	}
func getModalData(data discordgo.ModalSubmitInteractionData) map[handlers.Component]*discordgo.TextInput {
	var options = make(map[handlers.Component]*discordgo.TextInput)
	for _, actionRow := range data.Components {
		for _, c := range actionRow.(*discordgo.ActionsRow).Components {
			switch c := c.(type) {
			case *discordgo.TextInput:
				options[handlers.Component(c.CustomID)] = c
			default:
				log.Fatalf("Wrong component type: %T, skipping...", c)
			}
		}
	}
	return options
}
