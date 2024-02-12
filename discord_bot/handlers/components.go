package handlers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Component string

const (
	DeleteButton Component = "delete_error_message"

	dismissButton Component = "dismiss_error_message"
	urlButton     Component = "url_button"
	urlDelete     Component = "url_delete"

	readmoreDismiss Component = "readmore_dismiss"

	paginationButtons Component = "pagination_button"
	okCancelButtons   Component = "ok_cancel_buttons"

	Cancel    Component = "cancel"
	Interrupt Component = "interrupt"

	CancelDisabled    Component = "cancel_disabled"
	InterruptDisabled Component = "interrupt_disabled"

	JSONInput Component = "raw"

	roleSelect = "role_select"
)

var minValues = 1

var Components = map[Component]discordgo.MessageComponent{
	DeleteButton: discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Delete this message",
				Style:    discordgo.DangerButton,
				CustomID: string(DeleteButton),
				Emoji: &discordgo.ComponentEmoji{
					Name: "🗑️",
				},
			},
		},
	},
	urlButton: discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label: "Read more",
				Style: discordgo.LinkButton,
			},
		},
	},
	urlDelete: discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label: "Read more",
				Style: discordgo.LinkButton,
				Emoji: &discordgo.ComponentEmoji{
					Name: "📜",
				},
			},
			discordgo.Button{
				Label:    "Delete",
				Style:    discordgo.DangerButton,
				CustomID: string(DeleteButton),
				Emoji: &discordgo.ComponentEmoji{
					Name: "🗑️",
				},
			},
		},
	},
	dismissButton: discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Dismiss",
				Style:    discordgo.SecondaryButton,
				CustomID: string(DeleteButton),
			},
		},
	},
	readmoreDismiss: discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Read more",
				Style:    discordgo.LinkButton,
				CustomID: string(urlButton),
			},
			discordgo.Button{
				Label:    "Dismiss",
				Style:    discordgo.SecondaryButton,
				CustomID: string(DeleteButton),
			},
		},
	},

	paginationButtons: discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Previous",
				Style:    discordgo.SecondaryButton,
				CustomID: string(paginationButtons + "_previous"),
			},
			discordgo.Button{
				Label:    "Next",
				Style:    discordgo.SecondaryButton,
				CustomID: string(paginationButtons + "_next"),
			},
		},
	},
	okCancelButtons: discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "OK",
				Style:    discordgo.SuccessButton,
				CustomID: string(okCancelButtons + "_ok"),
			},
			discordgo.Button{
				Label:    "Cancel",
				Style:    discordgo.DangerButton,
				CustomID: string(okCancelButtons + "_cancel"),
			},
		},
	},
	roleSelect: discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.SelectMenu{
				MenuType:    discordgo.RoleSelectMenu,
				CustomID:    roleSelect,
				Placeholder: "Pick a role",
			},
		},
	},
	Cancel: discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Cancel",
				Style:    discordgo.DangerButton,
				CustomID: string(Cancel),
			},
		},
	},
	CancelDisabled: discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Cancel",
				Style:    discordgo.DangerButton,
				CustomID: string(Cancel),
				Disabled: true,
			},
		},
	},
	Interrupt: discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Interrupt",
				Style:    discordgo.DangerButton,
				CustomID: string(Interrupt),
				Emoji: &discordgo.ComponentEmoji{
					Name: "⚠️",
				},
				Disabled: false,
			},
		},
	},
	InterruptDisabled: discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Interrupt",
				Style:    discordgo.DangerButton,
				CustomID: string(Interrupt),
				Emoji: &discordgo.ComponentEmoji{
					Name: "⚠️",
				},
				Disabled: true,
			},
		},
	},

	JSONInput: discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID:    string(JSONInput),
				Label:       "JSON blob",
				Style:       discordgo.TextInputParagraph,
				Placeholder: "{\"height\":768,\"width\":512,\"prompt\":\"masterpiece\"}",
				Value:       "",
				Required:    true,
				MinLength:   1,
				MaxLength:   4000,
			},
		},
	},
}

func ModelSelectMenu(ID Component) discordgo.ActionsRow {
	display := strings.TrimPrefix(string(ID), "imagine_")
	display = strings.TrimSuffix(string(ID), "_model_name_menu")
	return discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.SelectMenu{
				CustomID:    string(ID),
				Placeholder: fmt.Sprintf("Change %s Model", display),
				MinValues:   &minValues,
				MaxValues:   1,
				Options: []discordgo.SelectMenuOption{
					{
						Label:       display,
						Value:       "Placeholder",
						Description: "Placeholder",
						Default:     false,
					},
				},
			},
		},
	}
}
