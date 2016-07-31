package card

import (
	"sort"
)

// CardInfo used for holding relevant search information.
type CardInfo struct {
	RawCard    *Card       `json:"rawCard"`
	Score      float32     `json:"score"`
	WordsFound []string    `json:"wordsFound"`
	Highlights []Highlight `json:"highlights"`
	Seen       int         `json:"seen"`
}

// Hightlights are for returning information to the UI about
// where we found particular search terms.
type Highlight struct {
	Field string `json:"field"` // name of the field we found the search term
	Text  string `json:"text"`  // snippet surrounding the search term to show context
}

// CardSets is the format of the cards we get from the API.
type CardSets map[string][]Card

// Cards is just an array of Card structs.
type Cards []Card

type Mechanic struct {
	Name string `json:"name"`
}

// Card is the core struct of the card service. It holds all the information
// for a specific card in the game.
type Card struct {
	// needs it's own type [{"name": "Charge"},{"name": "Divine Shield"}]
	Mechanics    []Mechanic `json:"mechanics"`
	Durability   int        `json:"durability"`
	Locale       string     `json:"locale"`
	Text         string     `json:"text"`
	HowToGet     string     `json:"howToGet"`
	ImgGold      string     `json:"imgGold"`
	Cost         int        `json:"cost"`
	Flavor       string     `json:"flavor"`
	PlayerClass  string     `json:"playerClass"`
	Img          string     `json:"img"`
	Attack       int        `json:"attack"`
	Health       int        `json:"health"`
	Type         string     `json:"type"`
	Collectible  bool       `json:"collectible"`
	Faction      string     `json:"faction"`
	Elite        bool       `json:"elite"`
	HowToGetGold string     `json:"howToGetGold"`
	CardSet      string     `json:"cardSet"`
	Name         string     `json:"name"`
	Artist       string     `json:"artist"`
	Rarity       string     `json:"rarity"`
	Race         string     `json:"race"`
	CardId       string     `json:"cardId"`
}

func NewCardInfo() CardInfo {
	return CardInfo{
		WordsFound: []string{},
		Highlights: []Highlight{},
		Score:      0,
		Seen:       0,
		RawCard:    nil,
	}
}

func CopyCardInfo(c CardInfo) CardInfo {
	n := NewCardInfo()
	n.WordsFound = make([]string, len(c.WordsFound))
	for i, word := range c.WordsFound {
		n.WordsFound[i] = word
	}
	n.Highlights = make([]Highlight, len(c.Highlights))
	for i, highlight := range c.Highlights {
		n.Highlights[i] = highlight
	}
	n.Score = c.Score
	n.Seen = c.Seen
	n.RawCard = c.RawCard
	return n
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
