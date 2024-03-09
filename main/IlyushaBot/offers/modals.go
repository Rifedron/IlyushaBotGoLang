package offers

import "github.com/bwmarrin/discordgo"

func feedbackModal(offerID string) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			Title:    "Комментарий к предложению",
			CustomID: "feedback|" + offerID,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "text",
							Label:       "Текст комментария",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "Предложение имба",
							Required:    true,
							MaxLength:   500,
						},
					},
				},
			},
		},
	}
}

func deletingModal(offerID string) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			Title:    "Удалить предложение",
			CustomID: "delete|" + offerID,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "text",
							Label:       "Причина удаления",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "Предложение мегаступид",
							Required:    true,
							MaxLength:   500,
						},
					},
				},
			},
		},
	}
}

func denyModal(offerID string) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			Title:    "Отклонить предложение",
			CustomID: "deny|" + offerID,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "text",
							Label:       "Причина отказа",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "Предложение ступид",
							Required:    true,
							MaxLength:   500,
						},
					},
				},
			},
		},
	}
}

func editModal(offerID string) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			Title:    "Изменить текст",
			CustomID: "edit|" + offerID,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:  "text",
							Label:     "Новый текст",
							Style:     discordgo.TextInputParagraph,
							Required:  true,
							MaxLength: 500,
						},
					},
				},
			},
		},
	}
}
