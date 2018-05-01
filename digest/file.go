package digest

import (
	"encoding/json"
	"os"

	"github.com/lee-woodridge/whats-that-card/card"
)

// GetCardsFromFile gets the card information from a local file and
// marshals it into our card struct.
func GetCardsFromFile(filename string) (card.Cards, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(f)
	// Error if any fields are not found in the Sets object.
	decoder.DisallowUnknownFields()

	var s card.Sets
	if err := decoder.Decode(&s); err != nil {
		return nil, err
	}
	cards := s.AllCards()
	return cards, nil
}
