package offers

import (
	"encoding/json"
	"fmt"
	"os"
)

func getOffer(messageID string) (*offer, error) {
	offerFile, err := os.Open("offers/" + messageID + ".json")
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(offerFile)

	newOffer := &offer{}
	err = decoder.Decode(newOffer)
	_ = offerFile.Close()

	if err != nil {
		fmt.Println("Возникла ошибка при считывании файла " + offerFile.Name())
		fmt.Println(err.Error())
	}

	return newOffer, err
}

func updateOfferFile(offer *offer) {
	bytes, _ := json.Marshal(offer)
	_ = os.WriteFile("offers/"+offer.MessageID+".json", bytes, os.ModePerm)
}

func createOfferFile(offer *offer) {
	newOfferFile, err := os.Create(`offers/` + offer.MessageID + ".json")
	if err != nil {
		fmt.Println(err.Error())
	}
	encoder := json.NewEncoder(newOfferFile)
	err = encoder.Encode(offer)
	if err != nil {
		fmt.Println(err.Error())
	}
	_ = newOfferFile.Close()
}

func removeOfferFile(messageID string) error {
	return os.Remove("offers/" + messageID + ".json")
}
