package search

import (
	"fmt"

	"github.com/lee-woodridge/whats-that-card/card"
)

type CardSet struct {
	m map[*card.Card]struct{}
}

var (
	exists = struct{}{}
)

func NewCardSet(c *card.Card) *CardSet {
	var cs CardSet
	cs.m = make(map[*card.Card]struct{})
	cs.m[c] = exists
	return &cs
}

func (cs *CardSet) Add(c *card.Card) {
	cs.m[c] = exists
}

func (cs *CardSet) Remove(c *card.Card) {
	delete(cs.m, c)
}

func (cs *CardSet) Contains(c *card.Card) bool {
	_, found := cs.m[c]
	return found
}

// func (cs *CardSet) String() string {
// 	var buffer bytes.Buffer
// 	for k := range cs.m {
// 		buffer.WriteString(k.Name)
// 	}
// 	return buffer.String()
// }

func (cs *CardSet) print() {
	for k := range cs.m {
		fmt.Printf("%s ", k.Name)
	}
	fmt.Printf("\n")
}
