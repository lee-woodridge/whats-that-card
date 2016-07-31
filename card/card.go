package card

import (
	"sort"
)

// CardInfo used for holding relevant search information.
type CardInfo struct {
	RawCard    *Card
	Score      float32
	WordsFound []string
	Seen       int
}

// CardSets is the format of the cards we get from the API.
type CardSets map[string][]Card

// Cards is just an array of Card structs.
type Cards []Card

type Mechanic struct {
	Name string
}

// Card is the core struct of the card service. It holds all the information
// for a specific card in the game.
type Card struct {
	Mechanics    []Mechanic // needs it's own type [{"name": "Charge"},{"name": "Divine Shield"}]
	Durability   int
	Locale       string
	Text         string
	HowToGet     string
	ImgGold      string
	Cost         int
	Flavor       string
	PlayerClass  string
	Img          string
	Attack       int
	Health       int
	Type         string
	Collectible  bool
	Faction      string
	Elite        bool
	HowToGetGold string
	CardSet      string
	Name         string
	Artist       string
	Rarity       string
	Race         string
	CardId       string
}

// AllCards takes the map format we get from the API and returns a simple array
// of all the cards. Sorts the keys to ensure we get deterministic loop ordering.
func (cs CardSets) AllCards() []Card {
	keys := []string{}
	for k, _ := range cs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	all := []Card{}
	for _, k := range keys {
		for _, card := range cs[k] {
			all = append(all, card)
		}
	}
	return all
}
