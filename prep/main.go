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

type SearchInfo struct {
	Cards []Card
	Trie  *trie.Trie
}

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
					t.Add(word, CardInfo{RawCard: &cards[c], Score: thisScore, WordsFound: []string{word}})
				}
			}
		}
	}
	return t, nil
}
