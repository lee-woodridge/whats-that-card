package digest

import (
	. "github.com/lee-woodridge/whats-that-card/card"

	"encoding/json"
	"io/ioutil"
)

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
