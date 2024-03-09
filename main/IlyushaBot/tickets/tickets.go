package tickets

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
)

var tickets = parseTickets()

func parseTickets() []*activeTicket {
	file, err := os.Open("tickets.json")
	defer file.Close()
	var tickets []*activeTicket
	if err != nil {
		_, _ = os.Create("tickets.json")
		return tickets
	}

	decoder := json.NewDecoder(file)
	_ = decoder.Decode(&tickets)
	return tickets
}

func getTicket(channelID string) *activeTicket {
	index := slices.IndexFunc(tickets, func(ticket *activeTicket) bool {
		return ticket.ChannelID == channelID
	})
	if index == -1 {
		return nil
	}
	return tickets[index]
}

func updateTickets() {
	bytes, _ := json.Marshal(tickets)
	err := os.WriteFile("tickets.json", bytes, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
}
