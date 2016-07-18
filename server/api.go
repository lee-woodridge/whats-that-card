package server

import (
	. "github.com/lee-woodridge/whats-that-card/card"
	// "github.com/lee-woodridge/whats-that-card/prep"

	"fmt"
	"net/http"
	"os"
	"strings"
)

// StartServer is the top level function for creating our card service.
//
// It pre-processes the card information such as setting up Tries for querying,
// then starts handling requests on the PORT environment variable.
func StartServer(cards Cards) {
	// Do pre-processing of the cards.
	// searchInfo := prep.CardInfoPrep(cards)
	// fmt.Printf("sea", ...)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/card/", getCard(cards))
	http.ListenAndServe(":"+port, nil)
}

// getCard is the function which handles the /card/ endpoint.
// It returns the card information for the card ID given.
func getCard(cards Cards) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cardId := strings.SplitN(r.URL.Path, "/", 3)[2]
		for _, card := range cards {
			if card.CardId == cardId {
				w.Write([]byte(fmt.Sprintf("%#v", card)))
				return
			}
		}
		http.Error(w, "Card Id doesn't exist", http.StatusBadRequest)
	}
}
