package offers

import (
	"awesomeProject/main/IlyushaBot"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strconv"
)

var OfferEvents = []interface{}{

	func(s *discordgo.Session, e *discordgo.MessageCreate) {
		if !e.Author.Bot && e.Message.Content != "" {
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
		o, err := getOffer(e.Message.ID)
		if err != nil {
			return
		}
		msg, _ := s.ChannelMessageSendComplex(e.ChannelID, &discordgo.MessageSend{
			Embed:      o.Embed,
			Components: votingButtons(o),
		})
		o.MessageID = msg.ID
		updateOfferFile(o)
	},
}

func offerMessage(s *discordgo.Session, msg *discordgo.Message) *discordgo.MessageSend {
	authorMember, _ := s.GuildMember(msg.GuildID, msg.Author.ID)
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("%s\nID: %s", authorMember.DisplayName(), authorMember.User.ID),
			IconURL: authorMember.AvatarURL(""),
		},
		Title:       "Предложение",
		Color:       0xFFFFFF,
		Description: msg.Content,
	}
	return &discordgo.MessageSend{
		Components: votingButtons(&offer{}),
		Embed:      embed,
	}
}

func votingButtons(o *offer) []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    strconv.Itoa(votersCount(o, HALAL)),
					CustomID: "halal",
					Style:    discordgo.SuccessButton,
					Emoji:    &discordgo.ComponentEmoji{Name: "👍"},
				},
				discordgo.Button{
					Label:    strconv.Itoa(votersCount(o, HARAM)),
					CustomID: "haram",
					Style:    discordgo.DangerButton,
					Emoji:    &discordgo.ComponentEmoji{Name: "👎"},
				},
			},
		},
	}
}
