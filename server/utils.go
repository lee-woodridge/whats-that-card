package server

import (
	. "github.com/lee-woodridge/whats-that-card/card"

	"sort"
)

type CardInfoArray []CardInfo

func (c CardInfoArray) Len() int      { return len(c) }
func (c CardInfoArray) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

// Return cards of the same score in alphabetical order (on Name).
func (c CardInfoArray) Less(i, j int) bool {
	if c[i].Score < c[j].Score {
		return false
	} else if c[i].Score == c[j].Score {
		return c[i].RawCard.Name < c[j].RawCard.Name
	} else {
		return true
	}
}

// CombineResults takes an array of result maps which are returned from
// our Trie functions, and returns a list of CardInfo ordered by their importance.
func CombineResults(in []map[string][]interface{}) []CardInfo {
	combine := make(map[*Card]CardInfo)
	for _, m := range in {
		for _, inters := range m {
			for _, inter := range inters {
				card, _ := inter.(CardInfo)
				if info, ok := combine[card.RawCard]; ok {
					// Combine scores.
					info.Score += card.Score
					info.WordsFound = append(info.WordsFound, card.WordsFound...)
				} else {
					combine[card.RawCard] = card
				}
			}
		}
	}
	res := make(CardInfoArray, len(combine))
	i := 0
	for _, info := range combine {
		res[i] = info
		i++
	}
	sort.Sort(res)
	return res
}
