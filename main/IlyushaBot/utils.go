package IlyushaBot

import (
	"github.com/bwmarrin/discordgo"
	"slices"
)

func moderatorValidCheck(s *discordgo.Session, guildsID string, m *discordgo.Member) bool {
	guildRoles, _ := s.GuildRoles(guildsID)
	modRole := guildRoles[slices.IndexFunc(guildRoles, func(role *discordgo.Role) bool {
		return role.ID == Cfg.OfferReplierRoleID
	})]
	for _, role := range m.Roles {
		if modRole.Position >= modRole.Position {
			return true
		}
	}
	return false
}
