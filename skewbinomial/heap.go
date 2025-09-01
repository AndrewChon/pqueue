package skewbinomial

import (
	"cmp"
)

type Heap[K cmp.Ordered, V any] struct {
	trees []*Tree[K, V]
}

func NewHeap[K cmp.Ordered, V any]() *Heap[K, V] {
	return &Heap[K, V]{
		trees: make([]*Tree[K, V], 0),
	}
}

func (h *Heap[K, V]) Insert(newKey K, newValue V) {
	if len(h.trees) >= 2 && h.trees[0].rank == h.trees[1].rank {
		newTree := skewLink(newKey, newValue, h.trees[0], h.trees[1])
		newTrees := []*Tree[K, V]{newTree}
		newTrees = append(newTrees, h.trees[2:]...)
		h.trees = newTrees
	} else {
		newTree := &Tree[K, V]{
			key:      newKey,
			value:    newValue,
			rank:     0,
			children: make([]*Tree[K, V], 0),
		}
		newTrees := []*Tree[K, V]{newTree}
		newTrees = append(newTrees, h.trees...)
		h.trees = newTrees
	}
}

func (h *Heap[K, V]) Merge(other *Heap[K, V]) {
	h.trees = Merge(h.trees, other.trees)
}

func (h *Heap[K, V]) FindMin() (*Tree[K, V], int) {
	if len(h.trees) == 0 {
		return nil, 0
	}

	minTree := h.trees[0]
	minI := 0
	for i := 1; i < len(h.trees); i++ {
		if h.trees[i].key < minTree.key {
			minTree = h.trees[i]
			minI = i
		}
	}

	return minTree, minI
}

func (h *Heap[K, V]) RemoveMin() {
	minTree, minI := h.FindMin()
	if minTree == nil {
		return
	}

	h.trees = append(h.trees[:minI], h.trees[minI+1:]...)

	rankZeroChildren := make([]*Tree[K, V], 0)
	rankNonZeroChildren := make([]*Tree[K, V], 0)

	for _, child := range minTree.children {
		if child.rank == 0 {
			rankZeroChildren = append(rankZeroChildren, child)
		} else {
			rankNonZeroChildren = append(rankNonZeroChildren, child)
		}
	}

	// Merge each non-zero rank child into the heap.
	h.trees = Merge(h.trees, rankNonZeroChildren)

	// Then Insert each zero-rank child into the heap.
	for _, zeroChild := range rankZeroChildren {
		h.Insert(zeroChild.key, zeroChild.value)
	}
}

type Tree[K cmp.Ordered, V any] struct {
	key      K
	value    V
	rank     int
	children []*Tree[K, V]
}

func (t *Tree[K, V]) Key() K {
	return t.key
}

func (t *Tree[K, V]) Value() V {
	return t.value
}

func Merge[K cmp.Ordered, V any](a, b []*Tree[K, V]) []*Tree[K, V] {
	return mergeUnique[K, V](uniquify(a), uniquify(b))
}

// simpleLink links together two trees of the same rank, with one becoming the leftmost child of the other. The
// resulting tree will have a rank of one greater than the rank of the two trees.
func simpleLink[K cmp.Ordered, V any](a, b *Tree[K, V]) *Tree[K, V] {
	if a.rank != b.rank {
		panic("cannot use simpleLink on trees of different ranks")
	}
	rank := a.rank

	if a.key <= b.key {
		newTree := &Tree[K, V]{
			key:   a.key,
			value: a.value,
			rank:  rank + 1,
		}

		newTree.children = []*Tree[K, V]{b}
		newTree.children = append(newTree.children, a.children...)

		return newTree
	} else {
		newTree := &Tree[K, V]{
			key:   b.key,
			value: b.value,
			rank:  rank + 1,
		}

		newTree.children = []*Tree[K, V]{a}
		newTree.children = append(newTree.children, b.children...)

		return newTree
	}
}

// skewLink links together three trees, one tree, a, having a rank of 0, and two trees, b and c, having the same rank as
// each other.
func skewLink[K cmp.Ordered, V any](aKey K, aValue V, b, c *Tree[K, V]) *Tree[K, V] {
	if b.rank != c.rank {
		panic("")
	}

	bcRank := b.rank

	// Type A
	if aKey <= b.key && aKey <= c.key {
		newTree := &Tree[K, V]{
			key:      aKey,
			value:    aValue,
			rank:     bcRank + 1,
			children: []*Tree[K, V]{b, c},
		}
		return newTree
	}

	a := &Tree[K, V]{
		key:      aKey,
		value:    aValue,
		rank:     0,
		children: make([]*Tree[K, V], 0),
	}

	// Type B
	if b.key <= c.key {
		newTree := &Tree[K, V]{
			key:      b.key,
			value:    b.value,
			rank:     bcRank + 1,
			children: []*Tree[K, V]{a, c},
		}
		newTree.children = append(newTree.children, b.children...)

		return newTree
	} else {
		newTree := &Tree[K, V]{
			key:      c.key,
			value:    c.value,
			rank:     bcRank + 1,
			children: []*Tree[K, V]{a, b},
		}
		newTree.children = append(newTree.children, c.children...)

		return newTree
	}
}

func uniquify[K cmp.Ordered, V any](trees []*Tree[K, V]) []*Tree[K, V] {
	if len(trees) < 2 {
		return trees
	}

	if trees[0].rank == trees[1].rank {
		simpleLinked := simpleLink(trees[0], trees[1])
		return uniquify[K, V](append([]*Tree[K, V]{simpleLinked}, trees[2:]...))
	} else {
		return append(trees[0:1], uniquify[K, V](trees[1:])...)
	}
}

func mergeUnique[K cmp.Ordered, V any](a, b []*Tree[K, V]) []*Tree[K, V] {
	if len(a) == 0 {
		return b
	} else if len(b) == 0 {
		return a
	}

	if a[0].rank < b[0].rank {
		return append(a[0:1], mergeUnique[K, V](a[1:], b)...)
	} else if a[0].rank > b[0].rank {
		return append(b[0:1], mergeUnique[K, V](a, b[1:])...)
	} else {
		simpleLinked := simpleLink(a[0], b[0])
		uniqueMerged := mergeUnique[K, V](a[1:], b[1:])
		array := append([]*Tree[K, V]{simpleLinked}, uniqueMerged...)
		return uniquify[K, V](array)
	}
}
