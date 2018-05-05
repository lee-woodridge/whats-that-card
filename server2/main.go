package server2

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/lee-woodridge/whats-that-card/card"
	"github.com/lee-woodridge/whats-that-card/digest"
	"github.com/lee-woodridge/whats-that-card/search"
)

const (
	cardsFilename = "./cards.json"
)

type CardScore struct {
	Card  card.Card `json:"card"`
	Score float64   `json:"score"`
}

type CardScores []CardScore

func (cs CardScores) Len() int {
	return len(cs)
}

func (cs CardScores) Swap(i, j int) {
	cs[i], cs[j] = cs[j], cs[i]
}

func (cs CardScores) Less(i, j int) bool {
	return cs[i].Score < cs[j].Score
}

// Cosine Similarity (query, doc) = Dot product(query, doc) / ||doc|| * ||query||
// ||doc|| = sqrt(TF-IDF(w1, doc)^2 + TF-IDF(w2, doc)^2)
// ||query|| = sqrt(TF-IDF(w1, query)^2 + TF-IDF(w2, query)^2)
// Dot product(query, doc) = (TF-IDF(w1, doc)*TF-IDF(w1, query)) + (TF-IDF(w2, doc)*TF-IDF(w2, query))
func cosineSimilarity(dotProduct, cardMagnitude, queryMagnitude float64) float64 {
	return dotProduct / (cardMagnitude * queryMagnitude)
}

func rankCards(queryTdIdfs map[string]float64, searchInfo *search.Info) CardScores {
	// Precompute query vector magnitude (so we don't have to every loop)
	var queryMagnitude float64
	for _, tdidf := range queryTdIdfs {
		queryMagnitude += (tdidf * tdidf)
	}
	queryMagnitude = math.Sqrt(queryMagnitude)

	cardToMagnitude := make(map[string]float64)
	cardToDotProduct := make(map[string]float64)
	for word, queryTdIdf := range queryTdIdfs {
		for card, cardTdIdf := range searchInfo.Tfidf[word] {
			cardToMagnitude[card] += (cardTdIdf * cardTdIdf)
			cardToDotProduct[card] += (queryTdIdf * cardTdIdf)
		}
	}
	var cardScores CardScores = []CardScore{}
	for card, dotProduct := range cardToDotProduct {
		cardMagnitude := math.Sqrt(cardToMagnitude[card])
		cardScores = append(cardScores, CardScore{
			Card:  searchInfo.CardMap[card],
			Score: cosineSimilarity(dotProduct, cardMagnitude, queryMagnitude)})
	}
	sort.Sort(cardScores)
	return cardScores
}

func tdidfForQuery(query string, searchInfo *search.Info) map[string]float64 {
	query = strings.ToLower(query)
	// Split string into words
	r := regexp.MustCompile("\\w+")
	queryWords := r.FindAllString(query, -1)
	queryTfIdfs := make(map[string]float64)
	for _, word := range queryWords {
		tf := (float64(1) / float64(len(queryWords)))
		idf := searchInfo.Idf[word]
		queryTfIdfs[word] = tf * idf
	}
	return queryTfIdfs
}

func query(searchInfo *search.Info) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := strings.SplitN(r.URL.Path, "/", 3)[2]
		fmt.Printf("query: %s\n", query)
		queryTdIdf := tdidfForQuery(query, searchInfo)
		fmt.Printf("query tdidf: %#v\n", queryTdIdf)
		cardScores := rankCards(queryTdIdf, searchInfo)
		fmt.Printf("card scores: %#v\n", cardScores)
		json.NewEncoder(w).Encode(cardScores)
		return
	}
}

func StartServer() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4201"
	}

	// Create search info.
	cards, err := digest.GetCardsFromFile(cardsFilename)
	if err != nil {
		return err
	}
	searchInfo, err := search.Prep(cards)
	if err != nil {
		return err
	}

	http.HandleFunc("/query/", query(searchInfo))
	fmt.Printf("Now serving on port: %s\n", port)
	http.ListenAndServe(":"+port, nil)
	return nil
}
