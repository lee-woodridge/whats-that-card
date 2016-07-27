// trie is a package created solely for use in text based searching
// of a static set of strings. For this reason, there is no option
// for removal or editing once the keys are inserted.
//
// I chose this approach using other trie libraries as they didn't
// keep multiples of metadata for duplicate keys, choosing rather to
// overwrite. In my case I require that these be kept.
//
// For example, in my hearthstone card searching application, the
// key "fire" may be used may times over, and I will want to retrieve
// all cards with this word rather than the most recent insert.
package trie

import ()

// Trie is the API level object, which contains a pointer to the
// root of our trie, and some top level metadata.
type Trie struct {
	root  *Node
	keys  int
	nodes int
}

// Node is the internal trie node struct.
type Node struct {
	r        rune
	parent   *Node
	children map[rune]*Node
	depth    int
	infos    []interface{}
}

//
// Trie functions.
//

// New returns a new empty Trie.
func New() *Trie {
	return &Trie{
		root: &Node{
			children: make(map[rune]*Node),
			infos:    []interface{}{},
		},
		keys:  0,
		nodes: 0,
	}
}

// Add inserts a set of info into the tree.
// It will insert all new nodes required if the key isn't present,
// or will insert the info at the correct node if it is.
func (t *Trie) Add(key string, info interface{}) {
	runes := []rune(key)
	t.keys++
	curNode := t.root
	// Iterate over the runes, we'll need to either get to the end
	// if there is already this key, to insert our info, or we'll
	// be creating new nodes all the way down.
	for i, r := range runes {
		if node, ok := curNode.children[r]; ok {
			curNode = node
		} else {
			t.nodes++
			curNode = curNode.addChild(r, i+1)
		}
	}
	curNode.addInfo(info)
}

// Find performs a full key search. It returns the Node in the
// tree at the end of the key, whether or not it's a leaf.
func (t *Trie) Find(key string) *Node {
	runes := []rune(key)
	curNode := t.root
	for _, r := range runes {
		if node, ok := curNode.children[r]; ok {
			curNode = node
		} else {
			return nil
		}
	}
	return curNode
}

// FindLeaf performs a full key search, and returns the node iff
// the string is found, and it is a leaf node.
func (t *Trie) FindLeaf(key string) *Node {
	node := t.Find(key)
	if node != nil {
		if node.isLeaf() {
			return node
		}
	}
	return nil
}

// PrefixSearch finds all nodes which are a child of the node which
// is the prefix. We return this as a map from the word up to the node
// to that node.
//
// For example: if we insert ("hi", 0) and ("hire", 1), we'll get back
// map[string]*Node{"hi": interface{0}, "hire": interface{1}}
func (t *Trie) PrefixSearch(prefix string) map[string][]interface{} {
	res := make(map[string][]interface{})
	// First, we need to find the node which represents the prefix.
	prefixNode := t.Find(prefix)
	if prefixNode == nil {
		return res
	}
	// Now we need to find all the nodes below this, with their associated word.
	prefixNode.collect(prefix, res)
	return res
}

// FuzzySearch performs a search for any nodes which represent a word which has
// a Levenshtein distance less than maxCost.
func (t *Trie) FuzzySearch(word string, maxCost int) map[string][]interface{} {
	res := make(map[string][]interface{})
	currentRow := make([]int, len(word)+1)
	for i, _ := range currentRow {
		currentRow[i] = i
	}
	for _, child := range t.root.children {
		child.fuzzyRecursive([]rune(word), []rune(""), currentRow, res, maxCost, true)
	}
	return res
}

// FuzzyPrefixSearch does a prefix search for all fuzzy paths found for our prefix,
// with a Levenshtein distance of 1 for len(prefix) <= 5 and 2 for len(prefix) > 5.
func (t *Trie) FuzzyPrefixSearch(prefix string) []map[string][]interface{} {
	potentialPrefixes := make(map[string][]interface{})
	prefixRunes := []rune(prefix)
	currentRow := make([]int, len(prefixRunes)+1)
	for i, _ := range currentRow {
		currentRow[i] = i
	}
	maxCost := 1
	if len(prefixRunes) > 5 {
		maxCost = 2
	}
	// Find potential paths which maxCost edits away from our prefix.
	for _, child := range t.root.children {
		child.fuzzyRecursive(prefixRunes, []rune(""), currentRow, potentialPrefixes, maxCost, false)
	}
	// Find prefix matches on these potential paths.
	results := make([]map[string][]interface{}, len(potentialPrefixes))
	i := 0
	for k, _ := range potentialPrefixes {
		results[i] = t.PrefixSearch(k)
		i++
	}
	// TODO: need to score on how far away the prefixes are.
	return results
}

// min is a basic min function for two ints, returning the smaller value.
func min(x, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

// fuzzyRecursive is the recursive function used by FuzzySearch, for recursing
// through the tree nodes. Accumulates the result in the res map.
func (n *Node) fuzzyRecursive(word, prefix []rune, prevRow []int,
	res map[string][]interface{}, maxCost int, needsInfo bool) {
	columns := len(word) + 1
	currentRow := []int{prevRow[0] + 1}
	for i := 1; i < columns; i++ {
		insertCost := currentRow[i-1] + 1
		deleteCost := prevRow[i] + 1

		var replaceCost int
		if word[i-1] != n.r {
			replaceCost = prevRow[i-1] + 1
		} else {
			replaceCost = prevRow[i-1]
		}

		currentRow = append(currentRow, min(min(insertCost, deleteCost), replaceCost))
	}

	if currentRow[len(currentRow)-1] <= maxCost && (!needsInfo || n.hasInfo()) {
		// fmt.Printf("adding %s with score %d\n", string(append(prefix, n.r)), currentRow[len(currentRow)-1])
		res[string(append(prefix, n.r))] = n.infos
	}

	min := currentRow[0]
	for _, val := range currentRow {
		if val < min {
			min = val
		}
	}
	if min <= maxCost {
		for _, child := range n.children {
			child.fuzzyRecursive(word, append(prefix, n.r), currentRow, res, maxCost, needsInfo)
		}
	}
}

//
// Node functions.
//

func (n *Node) collect(prefix string, res map[string][]interface{}) {
	res[prefix] = n.infos
	for r, n := range n.children {
		word := prefix + string(r)
		n.collect(word, res)
	}
}

// addChild adds a child node to node. It sets the childs parent to this node,
// and adds the new node as a child.
func (n *Node) addChild(r rune, depth int) *Node {
	// Probably unintended call to addChild.
	if _, ok := n.children[r]; ok {
		panic("Shouldn't add children when one exists.")
	}
	child := &Node{
		r:        r,
		parent:   n,
		children: make(map[rune]*Node),
		infos:    []interface{}{},
		depth:    depth,
	}
	n.children[r] = child
	return child
}

// hasInfo returns whether the current node has information or not. This is relevent
// as nodes which have info are nodes where a string terminated.
func (n *Node) hasInfo() bool {
	return len(n.infos) > 0
}

// isLeaf returns whether the current node is a leaf or not, by checking if
// the child map is empty.
func (n *Node) isLeaf() bool {
	return len(n.children) == 0
}

// addInfo simply adds the users info to the current node. This is used so that
// we can store multiple infos at a certain level to not lose information.
func (n *Node) addInfo(info interface{}) {
	n.infos = append(n.infos, info)
}
