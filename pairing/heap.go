package pairing

import (
	"cmp"
)

type Tree[K cmp.Ordered, V any] struct {
	key   K
	value V

	parent           *Tree[K, V]
	nextOlderSibling *Tree[K, V]
	youngestChild    *Tree[K, V]
}

func NewTree[K cmp.Ordered, V any](key K, value V) *Tree[K, V] {
	return &Tree[K, V]{
		key:              key,
		value:            value,
		parent:           nil,
		nextOlderSibling: nil,
		youngestChild:    nil,
	}
}

func (t *Tree[K, V]) Key() K {
	return t.key
}

func (t *Tree[K, V]) Value() V {
	return t.value
}

// FindMin returns the root node of the Tree, or nil if the Tree is empty.
func FindMin[K cmp.Ordered, V any](t *Tree[K, V]) *Tree[K, V] {
	if t == nil {
		return nil
	}
	return t
}

// Meld forms a new Tree from two other trees, with the largest becoming parent to the smallest.
func Meld[K cmp.Ordered, V any](a, b *Tree[K, V]) *Tree[K, V] {
	if a == nil {
		return b
	} else if b == nil {
		return a
	}

	if a.key < b.key {
		a.addChild(b)
		return a
	}

	b.addChild(a)
	return b
}

// Insert inserts a new node into a Tree.
func Insert[K cmp.Ordered, V any](t *Tree[K, V], new *Tree[K, V]) *Tree[K, V] {
	return Meld(t, new)
}

// RemoveMin removes the root node (in other words, the smallest node) from the provided Tree, rebuilds the Tree, and
// returns the new root node.
func RemoveMin[K cmp.Ordered, V any](t *Tree[K, V]) *Tree[K, V] {
	if t == nil {
		return nil
	}

	// We explicitly "disown" all its children. While not strictly necessary, it can free up the original root node for
	// GC earlier and proactively prevent any strange edge cases that may occur.
	youngestChild := t.youngestChild
	current := youngestChild
	for current != nil {
		next := current.nextOlderSibling
		current.parent = nil
		current = next
	}

	return twoPassMerge(youngestChild)
}

// DecreaseKey decreases the target node's key to the provided new key. The new key must be less than the target node's
// current key.
func DecreaseKey[K cmp.Ordered, V any](t *Tree[K, V], targetNode *Tree[K, V], newKey K) *Tree[K, V] {
	if t == nil || targetNode == nil {
		return t
	}

	if newKey >= targetNode.key {
		return t
	}
	targetNode.key = newKey

	if targetNode.parent == nil {
		return t
	}
	emancipate(targetNode)

	return Meld(t, targetNode)
}

// emancipate is a helper function that detaches a node from its parent.
func emancipate[K cmp.Ordered, V any](t *Tree[K, V]) {
	defer func() {
		t.parent = nil
		t.nextOlderSibling = nil
	}()

	parent := t.parent
	if parent.youngestChild == t {
		parent.youngestChild = t.nextOlderSibling
		return
	}

	ys := parent.youngestChild
	for ys != nil && ys.nextOlderSibling != t {
		ys = ys.nextOlderSibling
	}
	if ys != nil {
		ys.nextOlderSibling = t.nextOlderSibling
	}

}

// twoPassMerge reconstitutes a rootless Tree given its youngest child (and by extension, all of its children) into a
// new Tree, with its smallest member as its root.
func twoPassMerge[K cmp.Ordered, V any](yc *Tree[K, V]) *Tree[K, V] {
	if yc == nil || yc.nextOlderSibling == nil {
		return yc
	}

	var firstPassPairs []*Tree[K, V]

	// Meld the siblings in pairs, pairing the youngest sibling with the next older sibling.
	cur := yc
	for cur != nil {
		a := cur
		b := cur.nextOlderSibling

		if b != nil {
			cur = b.nextOlderSibling
			a.nextOlderSibling = nil
			b.nextOlderSibling = nil
			firstPassPairs = append(firstPassPairs, Meld(a, b))
		} else {
			cur = nil
			a.nextOlderSibling = nil
			firstPassPairs = append(firstPassPairs, a)
		}
	}

	// Meld together the first-pass pairs, but in the opposite direction to prevent the overall Tree from becoming
	// lopsided. The resulting Tree will now have the smallest as its Root.
	for i := len(firstPassPairs) - 2; i >= 0; i-- {
		firstPassPairs[i] = Meld(firstPassPairs[i], firstPassPairs[i+1])
	}

	return firstPassPairs[0]
}

func (t *Tree[K, V]) addChild(ct *Tree[K, V]) {
	ct.parent = t

	if t.youngestChild == nil {
		t.youngestChild = ct
		ct.nextOlderSibling = nil
	} else {
		ct.nextOlderSibling = t.youngestChild
		t.youngestChild = ct
	}
}
