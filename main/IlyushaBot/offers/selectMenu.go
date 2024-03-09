package offers

import (
	"awesomeProject/main/IlyushaBot"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func offerManageSelectMenuMessage(s *discordgo.Session, messageID string, o *offer, i *discordgo.InteractionCreate) *discordgo.InteractionResponse {
	if o.AuthorID == i.Member.User.ID {
		return selfOfferManageMenu(messageID)
	}
	gRoles, err := s.GuildRoles(i.GuildID)
	if err != nil {
		fmt.Println(err.Error())
	}
	if IlyushaBot.MemberHasSpecifiedRoleOrUpper(i.Member, IlyushaBot.Cfg.ModeratorRoleID, gRoles) {
		return offerReplyMenu(messageID)
	}
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			TTS:     false,
			Content: "Вы не можете управлять этим предложением",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}
}

func selfOfferManageMenu(messageID string) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							MenuType: discordgo.StringSelectMenu,
							CustomID: "selfOfferManage",
							Options: []discordgo.SelectMenuOption{
								{
									Label: "Изменить текст предложекния",
									Value: "edit|" + messageID,
									Emoji: &discordgo.ComponentEmoji{Name: "📝"},
								},
								{
									Label: "Удалить предложение",
									Value: "deleteMy|" + messageID,
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

func offerReplyMenu(messageID string) *discordgo.InteractionResponse {
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
