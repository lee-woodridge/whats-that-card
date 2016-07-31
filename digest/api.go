package digest

import (
	. "github.com/lee-woodridge/whats-that-card/card"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	hearthstoneAPI = "https://omgvamp-hearthstone-v1.p.mashape.com/cards"
)

// GetCardsFromFile gets the card information from the mashape hearthstone API
// and marshals it into our card struct.
func GetCardsFromAPI() (Cards, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", hearthstoneAPI, nil)
	if err != nil {
		return nil, err
	}
	// Get mashape api key from env variable (as not to check in to source control).
	// Can push env variable to heroku with:
	//		heroku config:add MASHAPE_KEY="$MASHAPE_KEY"
	mashapeKey := os.Getenv("MASHAPE_KEY")
	req.Header.Set("X-Mashape-Key", mashapeKey)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var cs CardSets
	if err := json.Unmarshal(bytes, &cs); err != nil {
		return nil, err
	}
	cards := cs.AllCards()
	return cards, nil
}
