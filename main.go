package main

import (
	"github.com/lee-woodridge/whats-that-card/digest"
	"github.com/lee-woodridge/whats-that-card/server"

	"fmt"
)

func main() {
	cards, err := digest.GetCardsFromFile("./cards.json")
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("cards[Basic][0]: %#v\n", cards[0])

	server.StartServer(cards)

	// if err := upload.InsertCardsToElastic(cards); err != nil {
	// 	fmt.Errorf("%s", err.Error())
	// }
}
