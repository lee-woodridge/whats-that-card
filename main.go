package main

import (
	"github.com/lee-woodridge/whats-that-card/digest"
	"github.com/lee-woodridge/whats-that-card/prep"
	"github.com/lee-woodridge/whats-that-card/server"
)

func main() {
	cards, err := digest.GetCardsFromFile("./cards.json")
	if err != nil {
		panic(err.Error())
	}

	// Do pre-processing of the cards.
	searchInfo, err := prep.CardInfoPrep(cards)
	if err != nil {
		panic(err)
	}

	server.StartServer(searchInfo)
}
