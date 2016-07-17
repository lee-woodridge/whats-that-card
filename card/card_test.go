package card

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestFillAllStrings(t *testing.T) {
	card := Card{Locale: "test", Text: "a", HowToGet: "b", ImgGold: "c"}
	card.FillAllStrings()
	assert.Equal(t, []string{"test", "a", "b", "c"}, card.AllStrings())
}
