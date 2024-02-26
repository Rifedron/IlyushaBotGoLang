package privateVCs

type privateVC struct {
	ChannelID      string `json:"channelID"`
	CreatorID      string `json:"creatorID"`
	CurrentOwnerID string `json:"currentOwnerID"`
}
