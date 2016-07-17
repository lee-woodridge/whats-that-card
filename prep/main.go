package prep

import (
	. "github.com/lee-woodridge/whats-that-card/card"
)

type CardInfo struct {
	RawCard    Card
	AllStrings []string
}

type SearchInfo struct {
	CardInfos []CardInfo
	Trie      *trie.Trie
}

func CardInfoPrep(cards []Card) SearchInfo {
	si := SearchInfo{}
	si.CardInfos = make([]CardInfo, len(cards))
	for i, card := range cards {
		si.CardInfos[i].RawCard = card
		si.CardInfos[i].AllStrings = card.GetAllStrings()
	}
	si.Trie = CreateCardTrie(si.CardInfos)
	return si
}
