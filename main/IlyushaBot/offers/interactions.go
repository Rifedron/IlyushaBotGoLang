package offers

import (
	"awesomeProject/main/IlyushaBot"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

var OfferCommands = []*discordgo.ApplicationCommand{
	{
		Name: "Управлять предложением",
		Type: discordgo.MessageApplicationCommand,
	},
}

var OfferInteractions = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	//Context commands
	"Управлять предложением": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		messageID := i.ApplicationCommandData().TargetID
		o, valid := isOfferManageValid(s, i, messageID)
		if valid {
			_ = s.InteractionRespond(i.Interaction, offerManageSelectMenuMessage(s, messageID, o, i))
		}
	},
	//Select menu
	"statusSelectMenu": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		value := strings.Split(i.MessageComponentData().Values[0], "|")
		option := value[0]
		offerID := value[1]
		offer, valid := isOfferManageValid(s, i, offerID)
		if valid {
			message, err := s.ChannelMessage(i.ChannelID, offerID)
			if err != nil {
				s.InteractionRespond(i.Interaction, IlyushaBot.EphemeralTextResponse("Сообщение было удалено"))
				return
			}
			embed := *message.Embeds[0]
			switch option {
			case "feedback":
				_ = s.InteractionRespond(i.Interaction, feedbackModal(offerID))
				break
			case "delete":
				_ = s.InteractionRespond(i.Interaction, deletingModal(offerID))
				break
			case "deny":
				s.InteractionRespond(i.Interaction, denyModal(offerID))
			default:
				newStatus := getStatusByID(option)
				mergeEmbedByStatus(&embed, newStatus)
				embed.Footer = embedFooter(s, i)
				go s.ChannelMessageEditEmbed(i.ChannelID, offerID, &embed)
				offer.Status = newStatus.StatusCode
				offer.Embed = &embed

				go updateOfferFile(offer)
				_ = s.InteractionRespond(i.Interaction, feedbackModal(offerID))
			}
		}
	},
	"selfOfferManage": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		value := strings.Split(i.MessageComponentData().Values[0], "|")
		option := value[0]
		offerID := value[1]
		offer, valid := isOfferManageValid(s, i, offerID)
		if valid {
			switch option {
			case "edit":
				if offer.Status != IGNORED.StatusCode {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseUpdateMessage,
						Data: &discordgo.InteractionResponseData{
							Content: "Нельзя изменять отвеченные предложения",
							Flags:   discordgo.MessageFlagsEphemeral,
						},
					})
					break
				}
				_ = s.InteractionRespond(i.Interaction, editModal(offerID))
				break
			case "deleteMy":
				go removeOfferFile(offerID)
				go s.ChannelMessageDelete(i.ChannelID, offerID)
				go s.ChannelMessageSendComplex(IlyushaBot.Cfg.OfferLogsChannelID, &discordgo.MessageSend{
					Content: fmt.Sprintf("Удалено создателем (<@%s>)", offer.AuthorID),
					Embed:   offer.Embed,
				})
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseUpdateMessage,
					Data: &discordgo.InteractionResponseData{
						Content: "Вы удалили своё предложение",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				break
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
	"deny": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		offerID := strings.Split(i.ModalSubmitData().CustomID, "|")[1]
		reason := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

		o, err := getOffer(offerID)
		if err != nil {
			return
		}

		message, err := s.ChannelMessage(i.ChannelID, offerID)
		if err != nil {
			return
		}
		embed := message.Embeds[0]
		embed.Footer = embedFooter(s, i)
		mergeEmbedByStatus(embed, DENIED)
		embed.Fields = []*discordgo.MessageEmbedField{
			{
				Name:  "Причина",
				Value: reason,
			},
		}

		o.Status = DENIED.StatusCode
		o.Embed = embed
		go s.ChannelMessageEditEmbed(i.ChannelID, offerID, embed)

		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Вы отклонили предложение",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})

		updateOfferFile(o)
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
	"edit": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		offerID := strings.Split(i.ModalSubmitData().CustomID, "|")[1]
		message, err := s.ChannelMessage(i.ChannelID, offerID)
		if err != nil {
			return
		}
		embed := *message.Embeds[0]
		embed.Description = i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
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
}

func vote(s *discordgo.Session, i *discordgo.InteractionCreate, messageID string, voteType VoteType) {
	offer, err := getOffer(messageID)
	if err != nil {
		return
	}
	if offer.AuthorID == i.Member.User.ID {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Вы не можете голосовать за своё предложение",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	mutateVotedOffer(offer, voteType, i.Member.User.ID)

	go updateOfferFile(offer)
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: updatedMessage(i.Message, offer),
	})
}

func mergeEmbedByStatus(embed *discordgo.MessageEmbed, status *Status) {
	embed.Title = "Предложение | " + status.DisplayName
	embed.Color = status.Color
}

func embedFooter(s *discordgo.Session, i *discordgo.InteractionCreate) *discordgo.MessageEmbedFooter {
	member, _ := s.GuildMember(i.GuildID, i.Member.User.ID)
	return &discordgo.MessageEmbedFooter{
		Text:    "Ответил " + member.DisplayName(),
		IconURL: member.AvatarURL(""),
	}
}

func isOfferManageValid(s *discordgo.Session, i *discordgo.InteractionCreate, messageID string) (*offer, bool) {
	offer, b := offerExists(s, i, messageID)
	return offer, b
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
