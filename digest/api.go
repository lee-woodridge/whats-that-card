package digest

import (
	. "github.com/lee-woodridge/whats-that-card/types"

	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	hearthstoneAPI = "https://omgvamp-hearthstone-v1.p.mashape.com/cards"
)

func GetCardsFromAPI() (CardSets, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", hearthstoneAPI, nil)
	if err != nil {
		return nil, err
	}
	// Get mashape api key from file (as not to check in to source control).
	file, err := ioutil.ReadFile("./mashape.key")
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Mashape-Key", string(file))
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var cards CardSets
	if err := json.Unmarshal(bytes, &cards); err != nil {
		return nil, err
	}
	return cards, nil
}
