package digest

import (
	. "github.com/lee-woodridge/whats-that-card/card"

	"encoding/json"
	"io/ioutil"
)

// GetCardsFromFile gets the card information from a local file and
// marshals it into our card struct.
func GetCardsFromFile(filename string) (Cards, error) {
	file, e := ioutil.ReadFile(filename)
	if e != nil {
		return nil, e
	}
	var cs CardSets
	if err := json.Unmarshal(file, &cs); err != nil {
		return nil, err
	}
	cards := cs.AllCards()
	return cards, nil
}
