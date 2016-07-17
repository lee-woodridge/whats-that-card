package card

import (
	"reflect"
	"sort"
)

type CardSets map[string][]Card
type Cards []Card

type Mechanic struct {
	Name string
}

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
	InPlayText   string
	Elite        bool
	HowToGetGold string
	CardSet      string
	Name         string
	Artist       string
	Rarity       string
	Race         string
	CardId       string
}

// Ugly code to loop over the map in alphabetical order.
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

func (c Card) GetAllStrings() []string {
	v := reflect.ValueOf(c)
	allStrings := []string{}
	for i := 0; i < v.NumField(); i++ {
		switch s := v.Field(i).Interface().(type) {
		case string:
			allStrings = append(allStrings, s)
		}
	}
	return allStrings
}
