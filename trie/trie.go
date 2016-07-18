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
	infos    []interface{}
}

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
	for _, r := range runes {
		if node, ok := curNode.children[r]; ok {
			curNode = node
		} else {
			t.nodes++
			curNode = curNode.addChild(r)
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

func (n *Node) collect(prefix string, res map[string][]interface{}) {
	res[prefix] = n.infos
	for r, n := range n.children {
		word := prefix + string(r)
		n.collect(word, res)
	}
}

// addChild adds a child node to node. It sets the childs parent to this node,
// and adds the new node as a child.
func (n *Node) addChild(r rune) *Node {
	// Probably unintended call to addChild.
	if _, ok := n.children[r]; ok {
		panic("Shouldn't add children when one exists.")
	}
	child := &Node{
		r:        r,
		parent:   n,
		children: make(map[rune]*Node),
		infos:    []interface{}{},
	}
	n.children[r] = child
	return child
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
