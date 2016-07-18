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

	server.StartServer(cards)
}
