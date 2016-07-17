package digest

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestGetCardsFromFile(t *testing.T) {
	cards, err := GetCardsFromFile("../cards.json")
	assert.Nil(t, err)
	assert.Equal(t, "AFK", cards["Basic"][0].Name, "First basic card should be AFK")
}
