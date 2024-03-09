package privateVCs

import (
	"awesomeProject/main/IlyushaBot"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

var PrivateVcEvents = []interface{}{
	func(s *discordgo.Session, e *discordgo.VoiceStateUpdate) {
		if e.BeforeUpdate != nil {
			go voiceQuit(s, e.BeforeUpdate)
		}
		if e.VoiceState.ChannelID != "" {
			voiceJoin(s, e.VoiceState)
		}
	},
	func(s *discordgo.Session, e *discordgo.Ready) {
		i := 0
		for _, private := range privates {
			channel, err := s.Channel(private.ChannelID)
			if err != nil {
				privates = append(privates[:i], privates[i+1:]...)
				continue
			}
			if voiceEmpty(s, channel.GuildID, channel.ID) {
				go s.ChannelDelete(channel.ID)
				privates = append(privates[:i], privates[i+1:]...)
				continue
			}
			if !isCurrentOwnerInVoice(s, channel.GuildID, private) {
				go s.ChannelPermissionDelete(channel.ID, private.CurrentOwnerID)
			}
			i++
		}
		updatePrivates()
	},
}

func voiceJoin(s *discordgo.Session, st *discordgo.VoiceState) {
	if st.ChannelID == IlyushaBot.Cfg.PrivatesFabricID {
		newPrivateVc(s, st)
		return
	}
	private := privateVoice(st.ChannelID)
	if private != nil {
		if private.CurrentOwnerID == st.UserID {
			err := s.ChannelPermissionSet(private.ChannelID, private.CurrentOwnerID, discordgo.PermissionOverwriteTypeMember,
				discordgo.PermissionManageChannels|discordgo.PermissionVoiceMoveMembers, 0)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}
func voiceQuit(s *discordgo.Session, st *discordgo.VoiceState) {
	private := privateVoice(st.ChannelID)
	if private != nil {
		if voiceEmpty(s, st.GuildID, private.ChannelID) {
			deletePrivate(s, private)
			return
		}
		if private.CurrentOwnerID == st.UserID {
			ch, err := s.Channel(private.ChannelID)
			if err != nil {
				return
			}
			for _, overwrite := range ch.PermissionOverwrites {
				if overwrite.ID == private.CurrentOwnerID {
					err := s.ChannelPermissionDelete(private.ChannelID, private.CurrentOwnerID)
					if err != nil {
						fmt.Println(err.Error())
					}
					break
				}
			}
		}
	}
}

func newPrivateVc(s *discordgo.Session, st *discordgo.VoiceState) {
	channel, err := s.GuildChannelCreateComplex(st.GuildID, discordgo.GuildChannelCreateData{
		Name:                 st.Member.DisplayName(),
		Type:                 discordgo.ChannelTypeGuildVoice,
		PermissionOverwrites: privateVoicePermOverwrites(st.UserID),
		ParentID:             IlyushaBot.Cfg.PrivatesCategoryID,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	go s.ChannelMessageSendEmbed(channel.ID, newPrivateEmbed(s, st))
	err = s.GuildMemberMove(st.GuildID, st.UserID, &channel.ID)
	if err != nil {
		_, _ = s.ChannelDelete(channel.ID)
		return
	}
	newPrivate := &privateVC{
		ChannelID:      channel.ID,
		CreatorID:      st.Member.User.ID,
		CurrentOwnerID: st.Member.User.ID,
	}
	privates = append(privates, newPrivate)
	updatePrivates()

	time.Sleep(500 * time.Millisecond)
	if voiceEmpty(s, st.GuildID, newPrivate.ChannelID) {
		deletePrivate(s, newPrivate)
	}
}

func newPrivateEmbed(s *discordgo.Session, st *discordgo.VoiceState) *discordgo.MessageEmbed {
	m, _ := s.GuildMember(st.GuildID, st.UserID)
	return &discordgo.MessageEmbed{
		Title: "Добро пожаловать в ваш приватный канал!",
		Description: "**/set owner** - назначить владельца\n\n" +
			"**/claim** - получить права на канал (если нет активного владельца)\n\n",
		Color: 0xd8958f,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Приятного общения!",
			IconURL: "https://cdn.discordapp.com/attachments/1190817984747405375/1193607671589372034/icons.png?ex=65ad54c5&is=659adfc5&hm=35488a51e3cb2a05ca4e9b0409ebde83c201e4fa8cfc7fce2f88915b76a25bd3&",
		},
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Создатель канала: " + m.DisplayName(),
			IconURL: m.AvatarURL(""),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Примечание:",
				Value: "первоначальный создатель комнаты может отобрать права даже у активного владельца",
			},
		},
	}
}

func deletePrivate(s *discordgo.Session, private *privateVC) {
	go s.ChannelDelete(private.ChannelID)
	for i, vc := range privates {
		if vc == private {
			privates = append(privates[:i], privates[i+1:]...)
			go updatePrivates()
			break
		}
	}
}

func voiceEmpty(s *discordgo.Session, guildID string, privateID string) bool {
	g, _ := s.State.Guild(guildID)
	for _, state := range g.VoiceStates {
		if state.ChannelID == privateID {
			return false
		}
	}
	return true
}

func privateVoicePermOverwrites(userID string) []*discordgo.PermissionOverwrite {
	return []*discordgo.PermissionOverwrite{
		{
			ID:    userID,
			Type:  discordgo.PermissionOverwriteTypeMember,
			Allow: discordgo.PermissionManageChannels | discordgo.PermissionVoiceMoveMembers,
		},
	}
}
