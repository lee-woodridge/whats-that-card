package search

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/lee-woodridge/whats-that-card/card"
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

func collectWords(cards []card.Card) (map[string]map[string]struct{}, error) {
	wordToCardsMap := make(map[string]map[string]struct{})
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
				r := regexp.MustCompile("\\w+")
				for _, word := range r.FindAllString(s, -1) {
					if m, found := wordToCardsMap[word]; found {
						m[thisCard.CardID] = struct{}{}
					} else {
						wordToCardsMap[word] = make(map[string]struct{})
						wordToCardsMap[word][thisCard.CardID] = struct{}{}
					}
				}
			}
		}
	}
	return wordToCardsMap, nil
}

func Prep(cards []card.Card) error {
	val, err := collectWords(cards)
	if err != nil {
		return err
	}
	for k, v := range val {
		fmt.Printf("%s -> %s\n", k, v)
	}
	return nil
}
