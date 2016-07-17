package main

import (
	"github.com/lee-woodridge/whats-that-card/digest"

	"fmt"
)

func main() {
	cards, err := digest.GetCardsFromFile()
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("cards[Basic][0]: %#v\n", cards["Basic"][0])

	// if err := upload.InsertCardsToElastic(cards); err != nil {
	// 	fmt.Errorf("%s", err.Error())
	// }
}
