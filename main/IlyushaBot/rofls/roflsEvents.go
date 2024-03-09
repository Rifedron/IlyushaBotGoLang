package rofls

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

var RoflsEvents = []interface{}{
	func(s *discordgo.Session, e *discordgo.GuildMemberAdd) {
		if e.User.ID == "1157367120921907220" {
			channel, err := s.UserChannelCreate("1157367120921907220")
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			for i := 0; i < 50; i++ {
				go s.ChannelMessageSend(channel.ID, "Подошёл близко начел лоскат пизду)))\nНачел лоскат пизду ещё бестрее")
			}
		}
	},
	func(s *discordgo.Session, e *discordgo.MessageCreate) {
		if e.Content == "%boba" {
			roles, err := s.GuildRoles(e.GuildID)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			var msgContent string
			for _, role := range roles {
				msgContent = fmt.Sprintf("%s %s %d\n", msgContent, role.Name, role.Position)
			}
			_, err = s.ChannelMessageSend(e.ChannelID, msgContent)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	},
}
