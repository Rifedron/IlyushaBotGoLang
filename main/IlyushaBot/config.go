package IlyushaBot

import (
	"encoding/json"
	"os"
)

type config struct {
	Token              string `json:"token"`
	OffersChannelId    string `json:"offersChannelId"`
	OfferReplierRoleID string `json:"offerReplierRoleID"`
	OfferLogsChannelID string `json:"offerLogsChannelID"`
	PrivatesFabricID   string `json:"privatesFabricID"`
	PrivatesCategoryID string `json:"privatesCategoryID"`
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
