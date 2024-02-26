package offers

import (
	"awesomeProject/main/IlyushaBot"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

var OfferEvents = []interface{}{

	func(s *discordgo.Session, e *discordgo.MessageCreate) {
		if !e.Author.Bot {
			if e.ChannelID == IlyushaBot.Cfg.OffersChannelId {
				err := s.ChannelMessageDelete(e.ChannelID, e.Message.ID)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				offerMessage, err := s.ChannelMessageSendComplex(e.ChannelID, offerMessage(s, e.Message))
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				newOffer := &offer{
					AuthorID:  e.Author.ID,
					MessageID: offerMessage.ID,
					Embed:     offerMessage.Embeds[0],
					Status:    0,
					Voters:    []*voter{},
				}
				createOfferFile(newOffer)
			}
		}
	},

	func(s *discordgo.Session, e *discordgo.MessageDelete) {
		_ = removeOfferFile(e.Message.ID)
	},
}

func offerMessage(s *discordgo.Session, msg *discordgo.Message) *discordgo.MessageSend {
	authorMember, _ := s.GuildMember(msg.GuildID, msg.Author.ID)
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    authorMember.DisplayName(),
			IconURL: authorMember.AvatarURL(""),
		},
		Title:       "Предложение",
		Color:       0xFFFFFF,
		Description: msg.Content,
	}
	return &discordgo.MessageSend{
		Components: votingButtons,
		Embed:      embed,
	}
}

var votingButtons = []discordgo.MessageComponent{
	discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "0",
				CustomID: "halal",
				Style:    discordgo.SuccessButton,
				Emoji:    &discordgo.ComponentEmoji{Name: "👍"},
			},
			discordgo.Button{
				Label:    "0",
				CustomID: "haram",
				Style:    discordgo.DangerButton,
				Emoji:    &discordgo.ComponentEmoji{Name: "👎"},
			},
		},
	},
}
