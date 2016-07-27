package server

import (
	// . "github.com/lee-woodridge/whats-that-card/card"
	"github.com/lee-woodridge/whats-that-card/prep"

	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// StartServer is the top level function for creating our card service.
//
// It pre-processes the card information such as setting up Tries for querying,
// then starts handling requests on the PORT environment variable.
func StartServer(cards prep.SearchInfo) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/card/", getCard(cards))
	http.HandleFunc("/search", search(cards))
	fmt.Printf("Now serving on port: %s\n", port)
	http.ListenAndServe(":"+port, nil)
}

// getCard is the function which handles the /card/ endpoint.
// It returns the card information for the card ID given.
func getCard(cards prep.SearchInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cardId := strings.SplitN(r.URL.Path, "/", 3)[2]
		for _, card := range cards.Cards {
			if card.CardId == cardId {
				json.NewEncoder(w).Encode(card)
				// w.Write([]byte(fmt.Sprintf("%#v", card)))
				return
			}
		}
		http.Error(w, "Card Id doesn't exist", http.StatusBadRequest)
	}
}

// search is the function which handles the /search endpoint.
// Queries are json encoded, and assumed to have the form:
// {
// 		"search": "query string"
// }
// It returns a list of cards which match the search term, json encoded.
func search(cards prep.SearchInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var query struct {
			Search string
		}
		if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		search := strings.ToLower(query.Search)
		words := strings.Split(search, " ")
		fullLength := len(words) - 1
		// Create storage for results.
		results := make([]map[string][]interface{}, fullLength+1)
		// Split search terms into full words and incomplete final word.
		full, prefix := words[:fullLength], words[fullLength]
		// Do fuzzy search for each full word.
		for i, word := range full {
			results[i] = cards.Trie.FuzzySearch(word, 0)
		}
		// Do prefix search for incomplete word.
		results[fullLength] = cards.Trie.PrefixSearch(prefix)

		// Combine results.
		combined := CombineResults(results)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(combined)
	}
}
