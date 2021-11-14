package state

import (
	"math/rand"
	"sync"
	"time"
)

type bstNgrams struct {
	mutex  sync.RWMutex
	root   *node
	random *rand.Rand
}

type node struct {
	mutex sync.RWMutex
	left  *node
	right *node
	key   string
	value *wordFreq
}

// NewBstNgrams creates a new Ngrams which implements Ngrams using a Binary Search Tree. The binary search tree can
// potentially have O(log n)search and insert. There is additional overhead due to the pointers for the structure,
// however the separation of the keys allows us to lock small parts of the tree to allow faster concurrency safe writes.
func NewBstNgrams() Ngrams {
	return &bstNgrams{
		mutex:  sync.RWMutex{},
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (b *bstNgrams) Store(words ...string) {
	key, word := getKeyAndWord(words...)
	if key == "" || word == "" {
		return
	}
	if b.root == nil {
		b.mutex.Lock()
		if b.root == nil {
			b.root = &node{
				key:   key,
				value: newWordFreq(),
			}
		}
		b.mutex.Unlock()
	}
	b.root.insert(key, word)
}

func (n *node) insert(key string, word string) {
	if n.key == key {
		n.value.add(word)
		return
	}
	if key < n.key {
		if n.left == nil {
			n.mutex.Lock()
			if n.left == nil {
				n.left = &node{
					value: newWordFreq(),
				}
			}
			n.mutex.Unlock()
		}
		n.left.insert(key, word)
		return
	}
	if key > n.key {
		if n.right == nil {
			n.mutex.Lock()
			if n.right == nil {
				n.right = &node{
					value: newWordFreq(),
				}
			}
			n.mutex.Unlock()
		}
		n.right.insert(key, word)
		return
	}
	return
}

func (b *bstNgrams) Get(words ...string) string {
	if b.root == nil {
		return ""
	}
	key := getKey(words...)
	n := b.root.getNode(key)
	if n == nil {
		return ""
	}
	if n.value == nil {
		return ""
	}
	return n.value.get(b.random)
}

func (n *node) getNode(key string) *node {
	if n.key == key {
		return n
	}
	if key < n.key {
		if n.left != nil {
			return n.left.getNode(key)
		}
	}
	if key > n.key {
		if n.right != nil {
			return n.right.getNode(key)
		}
	}
	return nil
}
