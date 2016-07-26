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

	assert.Equal(t, 1, trie.root.children[rune('t')].depth)
	assert.Equal(t, 2, trie.root.children[rune('t')].children[rune('e')].depth)

	node := trie.Find("test")
	assert.NotNil(t, node)
	assert.Equal(t, len(node.infos), 1)
	actual, ok := node.infos[0].(string)
	assert.True(t, ok)
	assert.Equal(t, "abc", actual)
}

func TestUTF8(t *testing.T) {
	trie := New()
	trie.Add("世界", "世界")

	node := trie.Find("世界")
	assert.NotNil(t, node)
	assert.Equal(t, len(node.infos), 1)
	actual, ok := node.infos[0].(string)
	assert.True(t, ok)
	assert.Equal(t, "世界", actual)
}

func TestInsertThenFindSubtree(t *testing.T) {
	trie := New()
	trie.Add("test", "abc")

	node := trie.Find("tes")
	assert.NotNil(t, node)
	assert.Equal(t, len(node.infos), 0)
}

func TestInsertThenFindLeafSubtreeFails(t *testing.T) {
	trie := New()
	trie.Add("test", "abc")

	node := trie.FindLeaf("tes")
	assert.Nil(t, node, "If we find nothing, returns nil")
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

	res = trie.PrefixSearch("te")
	expecteds = []struct {
		prefix string
		info   int
	}{
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

	res = trie.PrefixSearch("tes")
	expecteds = []struct {
		prefix string
		info   int
	}{
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

	res = trie.PrefixSearch("test")
	assert.Equal(t, len(res), 0)
}

func TestFuzzySearch(t *testing.T) {
	trie := New()

	expectedOnes := []string{
		"hell",   // 1 removal
		"ello",   // 1 removal front
		"rello",  // 1 replacements
		"hellor", // 1 addition
	}

	expectedTwos := []string{
		"hel",     // 2 removal
		"llo",     // 2 removal front
		"relli",   // 2 replacements
		"hellors", // 2 addition
		"ellor",   // 1 removal, 1 addition
	}

	// Set up the trie with all our strings.
	trie.Add("hello", nil) // 0 edits
	for _, val := range expectedOnes {
		trie.Add(val, nil)
	}
	for _, val := range expectedTwos {
		trie.Add(val, nil)
	}

	// Test we find "hello" in our 0 max cost.
	res := trie.FuzzySearch("hello", 0)
	assert.Equal(t, len(res), 1)
	assert.Contains(t, res, "hello")
	// and make sure we didn't find the one or two edits.
	for _, val := range append(expectedTwos, expectedOnes...) {
		assert.NotContains(t, res, val)
	}

	// Test we find one-edits in max one cost search.
	res = trie.FuzzySearch("hello", 1)
	assert.Equal(t, len(res), len(expectedOnes)+1)
	for _, val := range append(expectedOnes, "hello") {
		assert.Contains(t, res, val)
	}
	// and make sure we didn't find the two edits.
	for _, val := range expectedTwos {
		assert.NotContains(t, res, val)
	}

	// Test we find two-edits in max two cost search.
	res = trie.FuzzySearch("hello", 2)
	assert.Equal(t, len(res), len(expectedTwos)+len(expectedOnes)+1)
	for _, val := range append(expectedTwos, append(expectedOnes, "hello")...) {
		assert.Contains(t, res, val)
	}
}
