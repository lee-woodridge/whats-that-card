package digest

import (
	. "github.com/lee-woodridge/whats-that-card/types"

	"encoding/json"
	"io/ioutil"
)

func GetCardsFromFile(filename string) (CardSets, error) {
	file, e := ioutil.ReadFile(filename)
	if e != nil {
		return nil, e
	}
	var cards map[string][]Card
	if err := json.Unmarshal(file, &cards); err != nil {
		return nil, err
	}
	return cards, nil
}
