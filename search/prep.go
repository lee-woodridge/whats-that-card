package search

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strings"

	"github.com/lee-woodridge/whats-that-card/card"

	"github.com/surgebase/porter2"
)

var rankings = map[string]int{
	"Mechanics":       5,
	"Durability":      0,
	"Locale":          0,
	"Text":            8,
	"HowToGet":        0,
	"ImgGold":         0,
	"Cost":            0,
	"Flavor":          1,
	"PlayerClass":     3,
	"Img":             0,
	"Attack":          0,
	"Health":          0,
	"Type":            3,
	"Collectible":     0,
	"Faction":         5,
	"Elite":           0,
	"HowToGetGold":    0,
	"CardSet":         5,
	"Name":            10,
	"Artist":          1,
	"Rarity":          5,
	"Race":            5,
	"CardId":          0,
	"DbfID":           0,
	"MultiClassGroup": 1,
	"Classes":         1,
	"Armor":           0,
}

var containsWords = map[string]bool{
	"Mechanics":       true,
	"Durability":      false,
	"Locale":          false,
	"Text":            true,
	"HowToGet":        false,
	"ImgGold":         false,
	"Cost":            false,
	"Flavor":          true,
	"PlayerClass":     true,
	"Img":             false,
	"Attack":          false,
	"Health":          false,
	"Type":            true,
	"Collectible":     false,
	"Faction":         true,
	"Elite":           false,
	"HowToGetGold":    false,
	"CardSet":         true,
	"Name":            true,
	"Artist":          true,
	"Rarity":          true,
	"Race":            true,
	"CardId":          false,
	"DbfID":           false,
	"MultiClassGroup": true,
	"Classes":         true,
	"Armor":           false,
}

// Map from word seen -> card id -> term frequency
func calculateTermFrequency(cards []card.Card) (map[string]map[string]float64, error) {
	wordToCardsMap := make(map[string]map[string]float64)
	for _, thisCard := range cards {
		// Get the card structs info.
		v := reflect.ValueOf(thisCard)
		// For each field, check the type is a string.
		for i := 0; i < v.NumField(); i++ {
			switch s := v.Field(i).Interface().(type) {
			case string:
				// Skip this field if its score is 0.
				if !containsWords[v.Type().Field(i).Name] {
					continue
				}
				s = strings.ToLower(s)
				// Split string into words
				r := regexp.MustCompile("\\w+")
				wordList := r.FindAllString(s, -1)
				// For each word, add it to the word map, and add the card to the set.
				for _, word := range wordList {
					// Stem the word.
					word = porter2.Stem(word)
					if m, found := wordToCardsMap[word]; found {
						m[thisCard.CardID] += (float64(1) / float64(len(wordList)))
					} else {
						wordToCardsMap[word] = make(map[string]float64)
						wordToCardsMap[word][thisCard.CardID] = (float64(1) / float64(len(wordList)))
					}
				}
			}
		}
	}
	return wordToCardsMap, nil
}

// word -> idf
func calculateInverseDocumentFrequency(tf map[string]map[string]float64, numCards int) map[string]float64 {
	numCardsFloat := float64(numCards)
	idf := make(map[string]float64)
	for word, cards := range tf {
		idf[word] = math.Log(1 + (numCardsFloat / float64(len(cards))))
	}
	return idf
}

// Modifies the input tf to tf-idf
// word -> card -> tf-idf
func calculateTfIdf(tf map[string]map[string]float64, idf map[string]float64) {
	for word, cards := range tf {
		for card := range cards {
			tf[word][card] *= idf[word]
		}
	}
}

func Prep(cards []card.Card) error {
	tf, err := calculateTermFrequency(cards)
	if err != nil {
		return err
	}
	for k, v := range tf {
		fmt.Printf("%s -> %s\n", k, v)
	}
	idf := calculateInverseDocumentFrequency(tf, len(cards))
	calculateTfIdf(tf, idf)
	return nil
}
