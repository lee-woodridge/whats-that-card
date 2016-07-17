package trie

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestSingleInsertion(t *testing.T) {
	trie := New()
	trie.Add("test", "abc")
	assert.Equal(t, 1, trie.keys)
	assert.Equal(t, 4, trie.nodes)

	node := trie.Find("test")
	assert.NotNil(t, node)
	assert.Equal(t, len(node.infos), 1)
	actual, ok := node.infos[0].(string)
	assert.True(t, ok)
	assert.Equal(t, "abc", actual)
}

func TestMultiInsertion(t *testing.T) {
	trie := New()
	trie.Add("test", "abc")
	trie.Add("test2", "def")

	node := trie.Find("test")
	assert.NotNil(t, node)
	assert.Equal(t, 1, len(node.infos))

	actual, ok := node.infos[0].(string)
	assert.True(t, ok)
	assert.Equal(t, "abc", actual)

	node = trie.Find("test2")
	assert.NotNil(t, node)
	assert.Equal(t, 1, len(node.infos))

	actual, ok = node.infos[0].(string)
	assert.True(t, ok)
	assert.Equal(t, "def", actual)
}

func TestDuplicateInsertion(t *testing.T) {
	trie := New()
	trie.Add("test", "abc")
	trie.Add("test", "def")
	assert.Equal(t, 2, trie.keys)
	assert.Equal(t, 4, trie.nodes)

	node := trie.Find("test")
	assert.NotNil(t, node)
	assert.Equal(t, len(node.infos), 2)

	actual, ok := node.infos[0].(string)
	assert.True(t, ok)
	assert.Equal(t, "abc", actual)

	actual, ok = node.infos[1].(string)
	assert.True(t, ok)
	assert.Equal(t, "def", actual)
}

func TestPrefixSearch(t *testing.T) {
	trie := New()
	trie.Add("t", 1)
	trie.Add("te", 2)
	trie.Add("tes", 3)
	assert.Equal(t, 3, trie.keys)
	assert.Equal(t, 3, trie.nodes)

	res := trie.PrefixSearch("t")
	var expecteds = []struct {
		prefix string
		info   int
	}{
		{"t", 1},
		{"te", 2},
		{"tes", 3},
	}
	for _, expected := range expecteds {
		infos, ok := res[expected.prefix]
		assert.True(t, ok)
		assert.Equal(t, 1, len(infos))
		i, ok := infos[0].(int)
		assert.True(t, ok)
		assert.Equal(t, expected.info, i)
	}
}
