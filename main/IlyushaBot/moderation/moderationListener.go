package moderation

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

var ModerationEvents = []interface{}{
	func(s *discordgo.Session, e *discordgo.MessageCreate) {
		_, found := strings.CutPrefix(e.Content, "/")
		if found {

			if i != -1 {

			} else {
				s.ChannelMessageSendEmbedReply(e.ChannelID, noPermsCmdReply, e.MessageReference)
			}
		}
	},
}

var ModertationMessageCommands = map[string]func(s *discordgo.Session, e *discordgo.MessageCreate, args string){
	"mute": func(s *discordgo.Session, e *discordgo.MessageCreate, args string) {

	},
}

var noPermsCmdReply = &discordgo.MessageEmbed{
	Description: "❌У вас нет прав на использование команд",
	Color:       0xFF0000,
}
