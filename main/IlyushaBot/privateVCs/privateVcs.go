package privateVCs

import (
	"encoding/json"
	"os"
)

var privates = getPrivates()

func privateVoice(channelID string) *privateVC {
	for _, private := range privates {
		if private.ChannelID == channelID {
			return private
		}
	}
	return nil
}

func getPrivates() []*privateVC {
	file, err := os.Open("privates.json")
	var pivatesCollection []*privateVC
	if err != nil {
		_, _ = os.Create("privates.json")
		return pivatesCollection
	}
	decoder := json.NewDecoder(file)
	_ = decoder.Decode(&pivatesCollection)
	return pivatesCollection
}

func updatePrivates() {
	bytes, _ := json.Marshal(privates)
	_ = os.WriteFile("privates.json", bytes, os.ModePerm)
}
