package offers

import (
	"awesomeProject/main/IlyushaBot"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

var OfferCommands = []*discordgo.ApplicationCommand{
	{
		Name: "Ответить на предложение",
		Type: discordgo.MessageApplicationCommand,
	},
}

var OfferInteractions = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	//Context commands
	"Ответить на предложение": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		messageID := i.ApplicationCommandData().TargetID
		_, valid := isOfferReplyValid(s, i, messageID)
		if valid {
			_ = s.InteractionRespond(i.Interaction, offerReplySelectMenuMessage(messageID))
		}
	},
	//Select menu
	"statusSelectMenu": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		value := strings.Split(i.MessageComponentData().Values[0], "|")
		option := value[0]
		offerID := value[1]
		offer, valid := isOfferReplyValid(s, i, offerID)
		if valid {
			message, _ := s.ChannelMessage(i.ChannelID, offerID)
			embed := *message.Embeds[0]
			switch option {
			case "feedback":
				_ = s.InteractionRespond(i.Interaction, feedbackModal(offerID))
				break
			case "delete":
				_ = s.InteractionRespond(i.Interaction, deletingModal(offerID))
				break
			default:
				newStatus := getStatusByID(option)
				embed.Color = newStatus.Color
				embed.Title = "Предложение | " + newStatus.DisplayName
				embed.Footer = embedFooter(s, i)
				go s.ChannelMessageEditEmbed(i.ChannelID, offerID, &embed)
				offer.Status = newStatus.StatusCode
				offer.Embed = &embed

				go updateOfferFile(offer)
				_ = s.InteractionRespond(i.Interaction, feedbackModal(offerID))
			}
		}
	},
	//Buttons
	"halal": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		vote(s, i, i.Message.ID, HALAL)
	},
	"haram": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		vote(s, i, i.Message.ID, HARAM)
	},
	//Modals
	"feedback": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		offerID := strings.Split(i.ModalSubmitData().CustomID, "|")[1]
		message, err := s.ChannelMessage(i.ChannelID, offerID)
		if err != nil {
			return
		}
		embed := *message.Embeds[0]
		embed.Footer = embedFooter(s, i)
		embed.Fields = []*discordgo.MessageEmbedField{
			{
				Name:  "Комментарий",
				Value: i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
			},
		}
		go s.ChannelMessageEditEmbed(i.ChannelID, offerID, &embed)
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Вы оставили комментарий",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		offer, _ := getOffer(offerID)
		offer.Embed = &embed
		updateOfferFile(offer)
	},
	"delete": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		offerID := strings.Split(i.ModalSubmitData().CustomID, "|")[1]
		reason := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		go s.ChannelMessageDelete(i.ChannelID, offerID)
		deletedOffer, _ := getOffer(offerID)
		go s.ChannelMessageSendComplex(IlyushaBot.Cfg.OfferLogsChannelID, &discordgo.MessageSend{
			Content: fmt.Sprintf("Удалено модератором %s\n**Причина:** `%s`", i.Member.Mention(), reason),
			Embed:   deletedOffer.Embed,
		})
		go removeOfferFile(offerID)

		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Вы удалили предложение",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	},
}

func vote(s *discordgo.Session, i *discordgo.InteractionCreate, messageID string, voteType VoteType) {
	offer, err := getOffer(messageID)
	if err != nil {
		return
	}
	mutateVotedOffer(offer, voteType, i.Member.User.ID)

	go updateOfferFile(offer)
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: updatedMessage(i.Message, offer),
	})
}

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

func embedFooter(s *discordgo.Session, i *discordgo.InteractionCreate) *discordgo.MessageEmbedFooter {
	member, _ := s.GuildMember(i.GuildID, i.Member.User.ID)
	return &discordgo.MessageEmbedFooter{
		Text:    "Ответил " + member.DisplayName(),
		IconURL: member.AvatarURL(""),
	}
}

func isOfferReplyValid(s *discordgo.Session, i *discordgo.InteractionCreate, messageID string) (*offer, bool) {
	offer, b := offerExists(s, i, messageID)
	return offer, b && memberHasReplierRole(s, i, i.Member)
}

func offerExists(s *discordgo.Session, i *discordgo.InteractionCreate, messageID string) (*offer, bool) {
	offer, err := getOffer(messageID)
	if err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Предложение не существует",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return nil, false
	}
	return offer, true
}

func memberHasReplierRole(s *discordgo.Session, i *discordgo.InteractionCreate, member *discordgo.Member) bool {

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "У вас нет прав на рассмотрение предложений",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	return false
}

func roleById(roles discordgo.Roles, id string) *discordgo.Role {
	for _, role := range roles {
		if role.ID == id {
			return role
		}
	}
	return nil
}

func offerReplySelectMenuMessage(messageID string) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							MenuType: discordgo.StringSelectMenu,
							CustomID: "statusSelectMenu",
							Options: []discordgo.SelectMenuOption{
								menuOptionFromStatus(IMPLEMENTED, messageID),
								menuOptionFromStatus(ACCEPTED, messageID),
								menuOptionFromStatus(DENIED, messageID),
								{
									Label: "Изменить комментарий",
									Value: "feedback|" + messageID,
									Emoji: &discordgo.ComponentEmoji{Name: "📝"},
								},
								{
									Label: "Удалить предложение",
									Value: "delete|" + messageID,
									Emoji: &discordgo.ComponentEmoji{Name: "🗑"},
								},
							},
						},
					},
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	}
}

func menuOptionFromStatus(status *Status, messageID string) discordgo.SelectMenuOption {
	return discordgo.SelectMenuOption{
		Label: status.DisplayName,
		Value: fmt.Sprintf("%s|%s", status.ID, messageID),
		Emoji: &status.Emoji,
	}
}

func mutateVotedOffer(offer *offer, voteType VoteType, voterID string) {
	for i, v := range offer.Voters {
		if v.VoterID == voterID {
			if v.VoteType == voteType {
				offer.Voters = append(offer.Voters[:i], offer.Voters[i+1:]...)
				return
			}
			v.VoteType = voteType
			return
		}
	}
	offer.Voters = append(offer.Voters, &voter{
		VoterID:  voterID,
		VoteType: voteType,
	})
}

func updatedMessage(message *discordgo.Message, offer *offer) *discordgo.InteractionResponseData {
	embeds := message.Embeds
	actionsRow := message.Components[0].(*discordgo.ActionsRow)
	halalButton := actionsRow.Components[0].(*discordgo.Button)
	haramButton := actionsRow.Components[1].(*discordgo.Button)
	halalButton.Label = fmt.Sprintf("%d", votersCount(offer, HALAL))
	haramButton.Label = fmt.Sprintf("%d", votersCount(offer, HARAM))
	return &discordgo.InteractionResponseData{
		Embeds:     embeds,
		Components: []discordgo.MessageComponent{actionsRow},
	}
}

func votersCount(offer *offer, voteType VoteType) int {
	count := 0
	for _, v := range offer.Voters {
		if v.VoteType == voteType {
			count++
		}
	}
	return count
}
