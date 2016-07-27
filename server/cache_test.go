package server

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

// Test of minimum functionality, that we can add then get a result.
func TestCacheBasic(t *testing.T) {
	sc := NewSearchCache()

	term := "test"
	res, ok := sc.GetResult(term)
	assert.False(t, ok)
	assert.Nil(t, res)

	sc.AddResult(term, 123)

	// External api.
	res, ok = sc.GetResult(term)
	assert.True(t, ok)
	assert.Equal(t, 123, res)

	// Internal stuff.
	assert.Equal(t, len(sc.queue), 1)
	assert.Equal(t, sc.queue[0], term)
}

// Test we evict over max cache, and it's the oldest.
func TestCacheLimit(t *testing.T) {
	sc := NewSearchCache()
	sc.maxCount = 2 // simulate max cache size.

	var expecteds = []struct {
		term string
		info int
	}{
		{"test1", 1},
		{"test2", 2},
		{"test3", 3},
	}
	for _, val := range expecteds {
		sc.AddResult(val.term, val.info)
	}
	// First result should have been evicted.
	res, ok := sc.GetResult(expecteds[0].term)
	assert.False(t, ok)
	assert.Nil(t, res)
	// Other two should be there.
	res, ok = sc.GetResult(expecteds[1].term)
	assert.True(t, ok)
	assert.Equal(t, 2, res)

	res, ok = sc.GetResult(expecteds[2].term)
	assert.True(t, ok)
	assert.Equal(t, 3, res)

	// Queue should be in order, with first removed.
	assert.Equal(t, []string{"test2", "test3"}, sc.queue)
}

// Test duplicate requests don't fill the queue.
func TestCacheDupes(t *testing.T) {
	sc := NewSearchCache()
	sc.maxCount = 2 // simulate max cache size.

	var expecteds = []struct {
		term string
		info int
	}{
		{"test1", 1},
		{"test1", 1},
		{"test1", 1},
	}
	for _, val := range expecteds {
		sc.AddResult(val.term, val.info)
	}
	// Queue should still be length 1.
	assert.Equal(t, 1, len(sc.queue))
	assert.Equal(t, []string{"test1"}, sc.queue)

	res, ok := sc.GetResult(expecteds[0].term)
	assert.True(t, ok)
	assert.Equal(t, 1, res)
}
