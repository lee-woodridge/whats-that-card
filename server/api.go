package server

import (
	. "github.com/lee-woodridge/whats-that-card/card"
	// "github.com/lee-woodridge/whats-that-card/prep"

	"fmt"
	"net/http"
	"os"
	"strings"
)

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
