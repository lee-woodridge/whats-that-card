package server2

import (
	"github.com/lee-woodridge/whats-that-card/digest"
	"github.com/lee-woodridge/whats-that-card/search"
)

const (
	cardsFilename = "./cards.json"
)

func Start() error {
	cards, err := digest.GetCardsFromFile(cardsFilename)
	if err != nil {
		return err
	}
	err = search.Prep(cards)
	if err != nil {
		return err
	}
	return nil
}
