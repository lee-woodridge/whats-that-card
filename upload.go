package main

import (
	"github.com/lee-woodridge/whats-that-card/digest"
	"github.com/lee-woodridge/whats-that-card/upload"
)

func main() {
	cards, err := digest.GetCardsFromFile("./cards.json")
	if err != nil {
		panic(err.Error())
	}

	err = upload.InsertCardsToElastic(cards)
	if err != nil {
		panic(err.Error())
	}
}
