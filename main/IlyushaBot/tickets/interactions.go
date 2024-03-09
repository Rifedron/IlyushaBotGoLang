package tickets

import (
	"awesomeProject/main/IlyushaBot"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"slices"
)

var TicketInteractions = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	//Buttons
	"ticketCreate": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if ticketCreationValid(s, i) {
			ticketChannel, err := s.GuildChannelCreateComplex(i.GuildID, ticketChannelCreateData(i.Member, IlyushaBot.GuildDefaultRole(s, i.GuildID).ID))
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			go s.InteractionRespond(i.Interaction, IlyushaBot.EphemeralTextResponse("**Ваш тикет успешно создан**"+ticketChannel.Mention()))
			go func(s *discordgo.Session, mention string, ch *discordgo.Channel) {
				msg, _ := s.ChannelMessageSend(ch.ID, mention)
				_ = s.ChannelMessageDelete(ch.ID, msg.ID)
			}(s, i.Member.Mention(), ticketChannel)
			go s.ChannelMessageSendComplex(ticketChannel.ID, newTicketsChannelEmbed(i.Member))
			tickets = append(tickets, &activeTicket{
				InitiatorID: i.Member.User.ID,
				ChannelID:   ticketChannel.ID,
				Taken:       false,
			})
			updateTickets()
		}
	},
	"takeTicket": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		ticket := getTicket(i.ChannelID)
		if ticketTakingValid(s, i, ticket) {
			go s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{Content: i.Member.Mention() + "** закрепил тикет**"},
			})
			ticket.TakerID = i.Member.User.ID
			ticket.Taken = true
			go updateTickets()
			_, err := s.ChannelEdit(i.ChannelID, takenTicketChannelEdit(ticket, IlyushaBot.GuildDefaultRole(s, i.GuildID).ID))
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	},
	"closeTicket": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		ticket := getTicket(i.ChannelID)
		if ticketClosingValid(s, i, ticket) {
			go s.ChannelEdit(i.ChannelID, closedTicketChannelEdit(ticket, IlyushaBot.GuildDefaultRole(s, i.GuildID).ID))
			go s.InteractionRespond(i.Interaction, IlyushaBot.TextResponse("**Тикет закрыт**"))
			rmIndex := slices.Index(tickets, ticket)
			tickets = append(tickets[:rmIndex], tickets[rmIndex+1:]...)
			updateTickets()
		}
	},
}

func ticketCreationValid(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	index := slices.IndexFunc(tickets, func(t *activeTicket) bool { return t.InitiatorID == i.Member.User.ID })
	if index != -1 {
		ticket := tickets[index]
		_, err := s.Channel(ticket.ChannelID)
		if err != nil {
			tickets = append(tickets[:index], tickets[index+1:]...)
			return true
		}
		_ = s.InteractionRespond(i.Interaction, IlyushaBot.EphemeralTextResponse(
			"У вас уже есть открытый тикет <#"+tickets[index].ChannelID+">",
		))
		return false
	}
	return true
}

func ticketClosingValid(s *discordgo.Session, i *discordgo.InteractionCreate, t *activeTicket) bool {
	if !ticketInteractionValid(s, i, t) {
		return false
	}
	gRoles, err := s.GuildRoles(i.GuildID)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if !IlyushaBot.MemberHasSpecifiedRoleOrUpper(i.Member, IlyushaBot.Cfg.ElderModRoleID, gRoles) {
		_ = s.InteractionRespond(i.Interaction, IlyushaBot.EphemeralTextResponse("У вас нет прав на закрытие тикета"))
		return false
	}
	if t.TakerID != i.Member.User.ID {
		_ = s.InteractionRespond(i.Interaction, IlyushaBot.EphemeralTextResponse("Вы не можете закрыть тикет, который не брали"))
		return false
	}
	return true
}

func ticketTakingValid(s *discordgo.Session, i *discordgo.InteractionCreate, t *activeTicket) bool {
	if !ticketInteractionValid(s, i, t) {
		return false
	}
	roles, err := s.GuildRoles(i.GuildID)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if t.Taken {
		_ = s.InteractionRespond(i.Interaction, IlyushaBot.EphemeralTextResponse("Тикет уже взят <@"+t.TakerID+">"))
		return false
	}
	if t.InitiatorID == i.Member.User.ID {
		_ = s.InteractionRespond(i.Interaction, IlyushaBot.EphemeralTextResponse("Ты еблан?"))
		return false
	}
	if !IlyushaBot.MemberHasSpecifiedRoleOrUpper(i.Member, IlyushaBot.Cfg.ElderModRoleID, roles) {
		_ = s.InteractionRespond(i.Interaction, IlyushaBot.EphemeralTextResponse("У вас нет прав на взятие тикета"))
		return false
	}
	return true
}

func ticketInteractionValid(s *discordgo.Session, i *discordgo.InteractionCreate, t *activeTicket) bool {
	if t == nil {
		_ = s.InteractionRespond(i.Interaction, IlyushaBot.EphemeralTextResponse("Тикет уже закрыт либо не существует"))
		return false
	}
	return true
}

func ticketChannelCreateData(m *discordgo.Member, defRoleID string) discordgo.GuildChannelCreateData {
	return discordgo.GuildChannelCreateData{
		Name: "тикет-" + m.User.Username,
		Type: discordgo.ChannelTypeGuildText,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:    m.User.ID,
				Type:  discordgo.PermissionOverwriteTypeMember,
				Allow: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages,
			},
			{
				ID:    IlyushaBot.Cfg.ElderModRoleID,
				Type:  discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages,
			},
			{
				ID:   defRoleID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages,
			},
		},
		ParentID: IlyushaBot.Cfg.TicketsActiveCategoryID,
	}
}

func takenTicketChannelEdit(ticket *activeTicket, defRoleID string) *discordgo.ChannelEdit {
	return &discordgo.ChannelEdit{
		ParentID: IlyushaBot.Cfg.TicketsConsiderationCategoryID,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:    ticket.InitiatorID,
				Type:  discordgo.PermissionOverwriteTypeMember,
				Allow: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages,
			},
			{
				ID:    ticket.TakerID,
				Type:  discordgo.PermissionOverwriteTypeMember,
				Allow: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages,
			},
			{
				ID:   defRoleID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages,
			},
		},
	}
}

func closedTicketChannelEdit(ticket *activeTicket, defRoleID string) *discordgo.ChannelEdit {
	return &discordgo.ChannelEdit{
		ParentID: IlyushaBot.Cfg.TicketsClosedCategoryID,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:    ticket.InitiatorID,
				Type:  discordgo.PermissionOverwriteTypeMember,
				Allow: discordgo.PermissionViewChannel,
				Deny:  discordgo.PermissionSendMessages,
			},

			{
				ID:    IlyushaBot.Cfg.ElderModRoleID,
				Type:  discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionViewChannel,
				Deny:  discordgo.PermissionSendMessages,
			},
			{
				ID:   defRoleID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages,
			},
		},
	}
}

func newTicketsChannelEmbed(m *discordgo.Member) *discordgo.MessageSend {
	return &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Тикет " + m.DisplayName(),
				Description: "Опишите причину своего обращения здесь\n" +
					"В ближайшее время один из модераторов возьмёт ваш тикет",
				Color: 0xd8958f,
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Style: discordgo.PrimaryButton,
						Emoji: &discordgo.ComponentEmoji{
							Name: "📌",
						},
						CustomID: "takeTicket",
					},
					discordgo.Button{
						Style: discordgo.DangerButton,
						Emoji: &discordgo.ComponentEmoji{
							Name: "🔒",
						},
						CustomID: "closeTicket",
					},
				},
			},
		},
	}
}
