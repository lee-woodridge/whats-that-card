package prep

import (
	// External
	"github.com/stretchr/testify/assert"

	// Internal
	. "github.com/lee-woodridge/whats-that-card/card"
	"github.com/lee-woodridge/whats-that-card/digest"

	// Built in
	"testing"
)

func TestTrieFindId(t *testing.T) {
	cards, _ := digest.GetCardsFromFile("../cards.json")
	si := CardInfoPrep(cards)
	var expected Card
	for _, c := range si.CardInfos {
		if c.RawCard.CardId == "HERO_09" {
			expected = c.RawCard
		}
	}
	node, ok := si.Trie.Find("HERO_09")
	assert.True(t, ok)
	actual, ok := node.Meta().(Card)
	assert.True(t, ok)
	assert.Equal(t, expected, actual)
}
