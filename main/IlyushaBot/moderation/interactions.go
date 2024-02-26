package moderation

import "discordgo"

var moderationCommands = []*discordgo.ApplicationCommand{
	{
		Type: discordgo.ChatApplicationCommand,
		Name: "mute",

		DefaultPermission:        nil,
		DefaultMemberPermissions: nil,
		DMPermission:             nil,
		NSFW:                     nil,
		Description:              "",
		DescriptionLocalizations: nil,
		Options:                  nil,
	},
}
