package server

import (
	. "github.com/lee-woodridge/whats-that-card/card"

	"encoding/json"
	"fmt"
	"os"
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
	// For each word, we only want the highest score.
	// For example, we find the word "shield" 4 times on the blood knight card,
	// but we don't want that to weight the score so heavily.
	highScores := make([]map[string]CardInfo, len(in))
	// Turn each result set into a []CardInfo.
	for i, m := range in {
		scores := make(map[string]CardInfo)
		for _, inters := range m {
			for _, inter := range inters {
				card, _ := inter.(CardInfo)
				if ci, ok := scores[card.RawCard.CardId]; ok {
					if card.Score > ci.Score {
						scores[card.RawCard.CardId] = card
					}
				} else {
					scores[card.RawCard.CardId] = card
				}
			}
		}
		highScores[i] = scores
	}

	// Combine result sets.
	// We add the scores together for the same card.
	combine := make(map[string]CardInfo)
	for _, list := range highScores {
		for _, card := range list {
			if info, ok := combine[card.RawCard.CardId]; ok {
				// Combine scores.
				info.Score += card.Score
				info.Seen++
				info.WordsFound = append(info.WordsFound, card.WordsFound...)
				combine[card.RawCard.CardId] = info
			} else {
				card.Seen++
				combine[card.RawCard.CardId] = card
			}
		}
	}

	res := make(CardInfoArray, len(combine))
	i := 0
	for _, info := range combine {
		res[i] = info
		// Use the "seen" field to reduce the score of any cards which
		// appeared in one search but not another.
		if info.Seen < len(in) {
			missing := len(in) - info.Seen
			info.Score *= pow(0.5, missing) // half for each missing term.
		}
		i++
	}
	// Sort in order of score, for displaying.
	sort.Sort(res)
	return res
}

func pow(x float32, pow int) float32 {
	res := x
	for i := 1; i < pow; i++ {
		res *= x
	}
	return res
}

func prettyPrintJSON(toPrint interface{}) {
	b, err := json.MarshalIndent(toPrint, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
}
