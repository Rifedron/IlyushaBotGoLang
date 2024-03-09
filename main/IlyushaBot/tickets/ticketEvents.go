package tickets

import (
	"awesomeProject/main/IlyushaBot"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

var TicketEvents = []interface{}{
	func(s *discordgo.Session, e *discordgo.MessageCreate) {
		if e.Content == "%ticketsPanel" {
			gRoles, _ := s.GuildRoles(e.GuildID)
			if IlyushaBot.MemberHasSpecifiedRoleOrUpper(e.Member, IlyushaBot.Cfg.HighStaffRoleID, gRoles) {
				_ = s.ChannelMessageDelete(e.ChannelID, e.Message.ID)
				_, err := s.ChannelMessageSendComplex(e.ChannelID, ticketPanelMessage)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	},
}

var ticketPanelMessage = &discordgo.MessageSend{
	Embeds: []*discordgo.MessageEmbed{
		{
			Title:       "📝Приветствуем в канале поддержки!",
			Description: "Здесь вы можете оставлять жалобы и обращения к администрации",
			Color:       0xd8958f,
		},
	},
	Components: []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Style: discordgo.SecondaryButton,
					Emoji: &discordgo.ComponentEmoji{
						Name: "✏",
					},
					CustomID: "ticketCreate",
				},
			},
		},
	},
}
