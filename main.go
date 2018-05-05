package main

import (
	"fmt"
	"os"

	"github.com/lee-woodridge/whats-that-card/fetch"
	"github.com/lee-woodridge/whats-that-card/images"
	"github.com/lee-woodridge/whats-that-card/server2"
)

func main() {
	if len(os.Args) < 2 {
		os.Stderr.WriteString("argument required")
		return
	}
	switch os.Args[1] {
	case "fetch":
		fetch.Fetch()
		return
	case "images":
		images.GetImages()
		return
	case "server":
		err := server2.StartServer()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%#v", err)
		}
		return
	default:
		fmt.Fprintf(os.Stderr, "%s is not a valid argument", os.Args[1])
		return
	}

	// cards, err := digest.GetCardsFromAPI()
	// if err != nil {
	// 	panic(err.Error())
	// }

	// // Do pre-processing of the cards.
	// searchInfo, err := prep.CardInfoPrep(cards)
	// if err != nil {
	// 	panic(err)
	// }

	// server.StartServer(searchInfo)
}
