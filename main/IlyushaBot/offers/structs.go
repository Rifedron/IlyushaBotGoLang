package offers

import "github.com/bwmarrin/discordgo"

type offer struct {
	AuthorID  string                  `json:"authorID"`
	MessageID string                  `json:"messageID"`
	Status    statusCode              `json:"status"`
	Embed     *discordgo.MessageEmbed `json:"embed"`
	Voters    []*voter                `json:"voters"`
}

type Status struct {
	StatusCode  statusCode               `json:"statusCode"`
	Color       int                      `json:"color"`
	ID          string                   `json:"ID"`
	DisplayName string                   `json:"displayName"`
	Emoji       discordgo.ComponentEmoji `json:"emoji"`
}

type statusCode uint8

var (
	IGNORED = &Status{
		StatusCode: 0,
		Color:      0xffffff,
	}
	DENIED = &Status{
		StatusCode:  1,
		Color:       0xff0000,
		DisplayName: "Отклонено",
		ID:          "deny",
		Emoji:       discordgo.ComponentEmoji{Name: "❌"},
	}
	ACCEPTED = &Status{
		StatusCode:  2,
		Color:       0x07f71f,
		DisplayName: "Одобрено",
		ID:          "accept",
		Emoji:       discordgo.ComponentEmoji{Name: "✅"},
	}
	IMPLEMENTED = &Status{
		StatusCode:  3,
		Color:       0x59ffac,
		DisplayName: "Введено",
		ID:          "impl",
		Emoji:       discordgo.ComponentEmoji{Name: "✨"},
	}
)

var validStatuses = []*Status{
	ACCEPTED,
	DENIED,
	IMPLEMENTED,
}

func getStatusByCode(code statusCode) *Status {
	for _, status := range validStatuses {
		if code == status.StatusCode {
			return status
		}
	}
	return IGNORED
}

func getStatusByID(id string) *Status {
	for _, status := range validStatuses {
		if status.ID == id {
			return status
		}
	}
	return IGNORED
}

type voter struct {
	VoterID  string   `json:"voterID"`
	VoteType VoteType `json:"voteType"`
}

type VoteType uint8

const (
	HALAL VoteType = 0
	HARAM VoteType = 1
)
