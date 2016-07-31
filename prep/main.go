package prep

import (
	. "github.com/lee-woodridge/whats-that-card/card"
	"github.com/lee-woodridge/whats-that-card/trie"

	"encoding/json"
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"
)

// SearchInfo is the struct which holds the card information for querying.
type SearchInfo struct {
	Cards []Card
	Trie  *trie.Trie
}

// CardInfoPrep is the over-seer for all preparation of card metadata.
func CardInfoPrep(cards []Card) (SearchInfo, error) {
	si := SearchInfo{}
	si.Cards = cards
	var err error
	si.Trie, err = CreateCardTrie(cards)
	if err != nil {
		return si, err
	}
	return si, nil
}

// CreateCardTrie takes our list of cards and enters the relevant metadata
// into a Trie along with scores for importance, for later searching.
func CreateCardTrie(cards []Card) (*trie.Trie, error) {
	// Read in the scoring json.
	scoreFile, e := ioutil.ReadFile("./rankings.json")
	if e != nil {
		return nil, e
	}
	var scores map[string]int
	if err := json.Unmarshal(scoreFile, &scores); err != nil {
		return nil, err
	}
	// Create the trie.
	t := trie.New()
	for c, card := range cards {
		if !card.Collectible {
			continue
		}
		// Get the card structs info.
		v := reflect.ValueOf(card)
		// For each field, check the type is a string.
		for i := 0; i < v.NumField(); i++ {
			switch s := v.Field(i).Interface().(type) {
			case string:
				// Skip this field if its score is 0.
				thisScore := scores[v.Type().Field(i).Name]
				if thisScore == 0 {
					continue
				}
				// Clean the string, and split into words.
				s = strings.ToLower(s)
				r := regexp.MustCompile("\\w+")
				for _, word := range r.FindAllString(s, -1) {
					// Enter it into the trie, with the score associated with this field.
					ci := NewCardInfo()
					ci.RawCard = &cards[c]
					ci.Score = float32(thisScore)
					ci.WordsFound = append(ci.WordsFound, word)
					t.Add(word, ci)
				}
			}
		}
	}
	return t, nil
}
