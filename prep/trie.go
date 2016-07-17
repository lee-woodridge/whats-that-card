package prep

import (
	"github.com/derekparker/trie"
)

func CreateCardTrie(cardInfos []CardInfo) *trie.Trie {
	t := trie.New()
	for _, cardInfo := range cardInfos {
		for _, str := range cardInfo.AllStrings {
			t.Add(str, cardInfo.RawCard)
		}
	}
	return t
}
