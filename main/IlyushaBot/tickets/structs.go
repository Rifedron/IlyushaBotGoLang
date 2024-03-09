package tickets

type activeTicket struct {
	InitiatorID string `json:"initiatorID"`
	ChannelID   string `json:"channelID"`
	Taken       bool   `json:"taken"`
	TakerID     string `json:"takerID"`
}
