package server

import ()

const (
	// NUM_RESULTS is number of results to cache.
	// 1000 because each result ~100KB, so 1000*100KB = ~100MB.
	NUM_RESULTS = 1000
)

// SearchCache is the struct holding the cache information.
type SearchCache struct {
	// map from search term to result.
	cache map[string]interface{}
	// queue to know age of query. (evict should take oldest query)
	// order is queue[0] oldest -> queue[len-1] newest
	queue    []string
	maxCount int
}

func NewSearchCache() *SearchCache {
	return &SearchCache{
		cache:    make(map[string]interface{}),
		queue:    []string{},
		maxCount: NUM_RESULTS,
	}
}

// AddResult inserts a result into the cache.
// It handles when the cache is full, where it evicts the oldest entry.
// When a duplicate is passed to AddResult, it "refreshes" this entry,
// so it is pushed to the back of the age queue.
func (sc *SearchCache) AddResult(term string, res interface{}) {
	// If we've maxed the cache, remove a result.
	if len(sc.queue) >= sc.maxCount {
		var oldest string
		oldest, sc.queue = sc.queue[0], sc.queue[1:] // pop oldest query
		delete(sc.cache, oldest)
	}

	// If we've seen this term, remove it from the queue so we don't end
	// up with a queue full of the same common query.
	if _, ok := sc.cache[term]; ok {
		for i, val := range sc.queue {
			if val == term {
				sc.queue = append(sc.queue[:i], sc.queue[i+1:]...)
			}
		}
	}
	// Set the result in the map.
	sc.cache[term] = res
	// Push it onto the back of the queue.
	sc.queue = append(sc.queue, term)
}

// GetResult fetches a result from the cache.
// Returns the metadata (nil when not found) and a boolean reflecting success/failure.
func (sc *SearchCache) GetResult(term string) (interface{}, bool) {
	if val, ok := sc.cache[term]; ok {
		return val, true
	}
	return nil, false
}
