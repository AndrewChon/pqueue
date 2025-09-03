package skewbinomial

import (
	"cmp"
)

type Forest[K cmp.Ordered, V any] struct {
	trees []*Tree[K, V]
}

func NewForest[K cmp.Ordered, V any]() *Forest[K, V] {
	return new(Forest[K, V])
}

func (f *Forest[K, V]) Insert(newKey K, newValue V) {
	var newTree *Tree[K, V]

	if len(f.trees) >= 2 && f.trees[0].rank == f.trees[1].rank {
		newTree = skewLink(newKey, newValue, f.trees[0], f.trees[1])
		f.trees = prepend(f.trees[2:], newTree)
	} else {
		newTree = &Tree[K, V]{
			key:   newKey,
			value: newValue,
			rank:  0,
		}
		f.trees = prepend(f.trees, newTree)
	}
}

func (f *Forest[K, V]) Merge(other *Forest[K, V]) {
	f.trees = Merge(f.trees, other.trees)
}

func (f *Forest[K, V]) FindMin() (*Tree[K, V], int) {
	if len(f.trees) == 0 {
		return nil, 0
	}

	minTree := f.trees[0]
	minI := 0
	for i := 1; i < len(f.trees); i++ {
		if f.trees[i].key < minTree.key {
			minTree, minI = f.trees[i], i
		}
	}

	return minTree, minI
}

func (f *Forest[K, V]) RemoveMin() {
	minTree, minI := f.FindMin()
	if minTree == nil {
		return
	}

	f.Remove(minTree, minI)
}

func (f *Forest[K, V]) Remove(tree *Tree[K, V], i int) {
	f.trees = append(f.trees[:i], f.trees[i+1:]...)

	children := tree.children
	if len(children) == 0 {
		return
	}

	type kvp struct {
		key   K
		value V
	}

	zeroRankKVPs := make([]*kvp, 0)
	nonZeroRanked := children[:0]

	// Separate zero-rank children from non-zero-rank children.
	for _, child := range children {
		if child.rank == 0 {
			zeroRankKVPs = prepend(zeroRankKVPs, &kvp{child.key, child.value})
		} else {
			nonZeroRanked = prepend(nonZeroRanked, child)
		}
	}

	// Merge non-zero-rank children into the forest.
	if len(nonZeroRanked) > 0 {
		f.trees = Merge(f.trees, nonZeroRanked)
	}

	// Push zero-rank children back into the forest.
	if len(zeroRankKVPs) == 0 {
		return
	}

	for _, z := range zeroRankKVPs {
		f.Insert(z.key, z.value)
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
	var parent *Tree[K, V]
	var child *Tree[K, V]

	if a.key <= b.key {
		parent, child = a, b
	} else {
		parent, child = b, a
	}

	parent.children = prepend(parent.children, child)
	parent.rank++

	return parent
}

// skewLink links together three trees, one tree, a, having a rank of 0, and two trees, b and c, having the same rank as
// each other.
func skewLink[K cmp.Ordered, V any](aKey K, aValue V, b, c *Tree[K, V]) *Tree[K, V] {
	a := &Tree[K, V]{
		key:   aKey,
		value: aValue,
		rank:  0,
	}

	// Type A
	if aKey <= b.key && aKey <= c.key {
		a.rank = b.rank + 1
		a.children = append(a.children, b, c)
		return a
	}

	// Type B
	if b.key <= c.key {
		b.rank++
		b.children = prepend(b.children, a, c)
		return b
	} else {
		c.rank++
		c.children = prepend(c.children, a, b)
		return c
	}
}

func uniquify[K cmp.Ordered, V any](trees []*Tree[K, V]) []*Tree[K, V] {
	if len(trees) < 2 {
		return trees
	}

	if trees[0].rank == trees[1].rank {
		return uniquify[K, V](prepend(trees[2:], simpleLink(trees[0], trees[1])))
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
		return uniquify[K, V](prepend(mergeUnique[K, V](a[1:], b[1:]), simpleLink(a[0], b[0])))
	}
}

// prepend prepends elements to the beginning of a slice. If the slice does not have sufficient capacity, a new
// underlying array is allocated by appending zero values to the slice with the built-in function append. Otherwise,
// the contents of the slice are simply shifted to the right, and the provided elements are added to the front. Prepend
// tends to be more performant than simpler methods, such as creating a new slice of elements x and appending the
// original slice to it, especially if only one element is to be prepended.
func prepend[V any](slice []V, elems ...V) []V {
	n := len(slice) + len(elems)

	var zero V
	for cap(slice) < n {
		slice = append(slice[:cap(slice)], zero)
	}

	slice = slice[:n]
	copy(slice[len(elems):], slice)

	for i, v := range elems {
		slice[i] = v
	}

	return slice
}
