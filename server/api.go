package server

import (
	. "github.com/lee-woodridge/whats-that-card/card"
	"github.com/lee-woodridge/whats-that-card/prep"
	"github.com/lee-woodridge/whats-that-card/trie"

	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// SearchQuery is a struct reflecting the JSON request we expect.
type SearchQuery struct {
	Query    string
	Page     int
	PageSize int
}

// StartServer is the top level function for creating our card service.
//
// It pre-processes the card information such as setting up Tries for querying,
// then starts handling requests on the PORT environment variable.
func StartServer(cards prep.SearchInfo) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4201"
	}

	// Create cache for search.
	searchCache := NewSearchCache()

	http.HandleFunc("/card/", getCard(cards))
	http.HandleFunc("/search", search(cards, searchCache))
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

// sendResultJSON makes sure everything about sending a result is correctly
// handled, such as caching, setting response headers, paging etc.
func sendResultJSON(res []CardInfo, w http.ResponseWriter,
	searchCache *SearchCache, query *SearchQuery) {
	// Add result to cache.
	searchCache.AddResult(query.Query, res)
	totalResults := len(res)
	// Set headers (cors and json)
	w.Header().Set("Access-Control-Allow-Origin", "*") // TODO: do without "*"?
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if len(res) == 0 {
		json.NewEncoder(w).Encode([]interface{}{}) // Return {}
		return
	}
	// Slice result to get page we require.
	if query.Page*query.PageSize > (len(res) - 1) {
		http.Error(w, "Page size too high.", http.StatusBadRequest)
		return
	}
	lastIndex := (query.Page * query.PageSize) + query.PageSize
	if lastIndex > (len(res) - 1) {
		lastIndex = len(res) - 1
	}
	res = res[query.Page*query.PageSize : lastIndex]
	res = HighlightCards(res)
	// Add some metadata to the result for the UI.
	output := struct {
		Results []CardInfo `json:"results"`
		Count   int        `json:"count"`
		Total   int        `json:"total"`
	}{
		Results: res,
		Count:   len(res),
		Total:   totalResults,
	}
	json.NewEncoder(w).Encode(output)
}

// search is the function which handles the /search endpoint.
// Queries are json encoded, and assumed to have the form:
// {
// 		"search": "query string",
//		"page": 0, // indexed from 0 as first
//		"pageSize": 12 // any size
// }
// It returns a list of cards which match the search term, json encoded.
func search(cards prep.SearchInfo, searchCache *SearchCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := &SearchQuery{}
		if err := json.NewDecoder(r.Body).Decode(query); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		if res, ok := searchCache.GetResult(query.Query); ok {
			fmt.Printf("Result for %s found in cache.\n", query.Query)
			cards, ok := res.([]CardInfo)
			if !ok {
				http.Error(w, "Cache result of unexpected type.", http.StatusExpectationFailed)
			}
			sendResultJSON(cards, w, searchCache, query)
			return
		}

		search := strings.ToLower(query.Query)
		words := strings.Split(search, " ")
		fullLength := len(words) - 1
		// Create storage for results.
		var results []map[string][]interface{}
		// Split search terms into full words and incomplete final word.
		var full []string
		var prefix string
		names := []string{}
		if len([]rune(words[fullLength])) > 4 {
			results = make([]map[string][]interface{}, fullLength+2)
			full, prefix = words, words[fullLength]
		} else {
			results = make([]map[string][]interface{}, fullLength+1)
			full, prefix = words[:fullLength], words[fullLength]
		}

		// Do prefix search for incomplete word.
		i := 0
		if prefix != "" {
			results[0] = cards.Trie.PrefixSearch(prefix)
			names = append(names, prefix)
		}

		// Define our error function, which will run on each returned node.
		errorFunc := func(n *trie.Node, ld int) []interface{} {
			infos := n.Info()
			for i, info := range infos {
				card, _ := info.(CardInfo)
				// Reduce this card's score by 25% per Levenschtein Distance.
				card.Score *= pow(0.75, ld)
				infos[i] = card
			}
			return infos
		}
		// Do fuzzy search for each full word.
		for _, word := range full {
			i++
			errorMargin := 1
			if len([]rune(word)) < 6 {
				errorMargin = 1
			} else {
				errorMargin = 2
			}
			results[i] = cards.Trie.FuzzySearchFunc(word, errorMargin, errorFunc)
			names = append(names, word)
		}

		// Combine results.
		combined := CombineResults(results)

		fmt.Printf("Result for %s calculated with Trie.\n", query.Query)
		sendResultJSON(combined, w, searchCache, query)
	}
}
