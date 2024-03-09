package IlyushaBot

import (
	"encoding/json"
	"os"
)

type config struct {
	Token              string `json:"token"`
	OffersChannelId    string `json:"offersChannelId"`
	OfferLogsChannelID string `json:"offerLogsChannelID"`

	HighStaffRoleID string `json:"highStaffRoleID"`
	ElderModRoleID  string `json:"elderModRoleID"`
	ModeratorRoleID string `json:"moderatorRoleID"`

	PrivatesFabricID   string `json:"privatesFabricID"`
	PrivatesCategoryID string `json:"privatesCategoryID"`

	TicketsActiveCategoryID        string `json:"ticketsActiveCategoryID"`
	TicketsConsiderationCategoryID string `json:"ticketsConsiderationCategoryID"`
	TicketsClosedCategoryID        string `json:"ticketsClosedCategoryID"`
}

func parseConfig() *config {
	cfg := &config{}

	file, err := os.Open("config.json")
	if err != nil {
		_, _ = os.Create("config.json")
		panic("Баля где файл config.json паяснете")
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(cfg)
	if err != nil {
		panic("В конфиге насрано. Исправляй")
	}

	return cfg
}

var Cfg = parseConfig()
