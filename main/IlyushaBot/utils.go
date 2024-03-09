package IlyushaBot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"slices"
)

func MemberHasSpecifiedRoleOrUpper(m *discordgo.Member, specifiedRoleID string, gRoles []*discordgo.Role) bool {
	specifiedRoleIndex := slices.IndexFunc(gRoles, func(r *discordgo.Role) bool {
		return r.ID == specifiedRoleID
	})
	if specifiedRoleIndex == -1 {
		return false
	}
	specifiedRole := gRoles[specifiedRoleIndex]

	return slices.ContainsFunc(m.Roles, func(rID string) bool {
		rIndex := slices.IndexFunc(gRoles, func(r *discordgo.Role) bool { return r.ID == rID })
		role := gRoles[rIndex]
		return role.Position >= specifiedRole.Position
	})
}

func GuildDefaultRole(s *discordgo.Session, gID string) *discordgo.Role {
	roles, err := s.GuildRoles(gID)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	i := slices.IndexFunc(roles, func(role *discordgo.Role) bool { return role.Position == 0 })
	return roles[i]
}

func TextResponse(txt string) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: txt,
		},
	}
}

func EphemeralTextResponse(txt string) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: txt,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}
}
