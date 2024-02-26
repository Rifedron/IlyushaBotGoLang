package privateVCs

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

var PrivatesCommands = []*discordgo.ApplicationCommand{
	{
		Type:        discordgo.ChatApplicationCommand,
		Name:        "claim",
		Description: "Выдаёт права на этот приватный канал если их ни у кого нет",
	},
	{
		Type:        discordgo.ChatApplicationCommand,
		Name:        "set",
		Description: "Выдаёт права на этот канал выбранному человеку",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "owner",
				Description: "Выдаёт права на этот канал выбранному человеку",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "user",
						Description: "Тот, кто получит права на этот канал",
						Required:    true,
					},
				},
			},
		},
	},
}

var PrivatesInteractions = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"claim": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		private, b := isClaimInteractionValid(s, i)
		if b {
			go s.ChannelMessageSend(private.ChannelID, fmt.Sprintf(
				"**%s теперь владелец этого канала**", i.Member.Mention()))
			go s.ChannelPermissionDelete(private.ChannelID, private.CurrentOwnerID)
			go s.ChannelPermissionSet(private.ChannelID, i.Member.User.ID, discordgo.PermissionOverwriteTypeMember,
				discordgo.PermissionManageChannels|discordgo.PermissionVoiceMoveMembers, 0)
			private.CurrentOwnerID = i.Member.User.ID
			go updatePrivates()
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Теперь вы владелец этой приватки",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
	},
	"set": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		UserID := i.ApplicationCommandData().Options[0].Options[0].UserValue(s).ID
		private, b := isSetOwnerInteractionValid(s, i)
		if b {
			newOwner, _ := s.GuildMember(i.GuildID, UserID)
			go s.ChannelPermissionDelete(private.ChannelID, private.CurrentOwnerID)
			go s.ChannelPermissionSet(private.ChannelID, UserID, discordgo.PermissionOverwriteTypeMember,
				discordgo.PermissionManageChannels|discordgo.PermissionVoiceMoveMembers, 0)
			go s.ChannelMessageSend(private.ChannelID, fmt.Sprintf("**%s назначил %s владельцем этого канала**",
				i.Member.Mention(), newOwner.Mention()))
			private.CurrentOwnerID = UserID
			go updatePrivates()
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Вы назначили %s владельцем этого канала", newOwner.DisplayName()),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
	},
}

func isPrivateVcCommandValid(s *discordgo.Session, i *discordgo.InteractionCreate) (*privateVC, bool) {
	voiceID := currentMemberVoiceID(s, i.GuildID, i.Member.User.ID)
	if voiceID == "" {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Вы не находитесь в голосовом канале",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return nil, false
	}
	private := privateVoice(voiceID)
	if private == nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Вы не находитесь в приватке",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return nil, false
	}
	return private, true
}

func isClaimInteractionValid(s *discordgo.Session, i *discordgo.InteractionCreate) (*privateVC, bool) {
	private, b := isPrivateVcCommandValid(s, i)
	if !b {
		return private, b
	}
	if private.CurrentOwnerID == i.Member.User.ID {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Вы уже владеец этого канала",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return nil, false
	}
	if private.CreatorID == i.Member.User.ID {
		return private, true
	}
	if isCurrentOwnerInVoice(s, i.GuildID, private) {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "У этого канала уже есть владелец",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return nil, false
	}
	return private, true
}

func isSetOwnerInteractionValid(s *discordgo.Session, i *discordgo.InteractionCreate) (*privateVC, bool) {
	private, b := isPrivateVcCommandValid(s, i)
	if !b {
		return nil, false
	}
	if i.Member.User.ID != private.CurrentOwnerID {
		go s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Вы не владелец этого канала",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return nil, false
	}
	UserID := i.ApplicationCommandData().Options[0].Options[0].UserValue(s).ID
	if currentMemberVoiceID(s, i.GuildID, UserID) != private.ChannelID {
		go s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Выбранный пользователь не находится в этом канале",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return nil, false
	}
	return private, true
}

func isCurrentOwnerInVoice(s *discordgo.Session, guildID string, private *privateVC) bool {
	g, _ := s.State.Guild(guildID)
	for _, st := range g.VoiceStates {
		if st.ChannelID == private.ChannelID {
			if st.UserID == private.CurrentOwnerID {
				return true
			}
		}
	}
	return false
}

func currentMemberVoiceID(s *discordgo.Session, gID string, memberID string) string {
	g, _ := s.State.Guild(gID)
	for _, state := range g.VoiceStates {
		if state.UserID == memberID {
			return state.ChannelID
		}
	}
	return ""
}
