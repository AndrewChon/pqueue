package skew

import (
	"cmp"
)

type Tree[K cmp.Ordered, V any] struct {
	key   K
	value V
	left  *Tree[K, V]
	right *Tree[K, V]
}

func NewTree[K cmp.Ordered, V any](key K, value V) *Tree[K, V] {
	return &Tree[K, V]{
		key:   key,
		value: value,
		left:  nil,
		right: nil,
	}
}

func (t *Tree[K, V]) Key() K {
	return t.key
}

func (t *Tree[K, V]) Value() V {
	return t.value
}

func FindMin[K cmp.Ordered, V any](t *Tree[K, V]) *Tree[K, V] {
	if t == nil {
		return nil
	}
	return t
}

func Meld[K cmp.Ordered, V any](a, b *Tree[K, V]) *Tree[K, V] {
	if a == nil {
		return b
	} else if b == nil {
		return a
	}

	if b.key < a.key {
		a, b = b, a
	}

	a.right, a.left = a.left, Meld(b, a.right)
	return a
}

func Insert[K cmp.Ordered, V any](t *Tree[K, V], new *Tree[K, V]) *Tree[K, V] {
	return Meld(t, new)
}

func RemoveMin[K cmp.Ordered, V any](t *Tree[K, V]) *Tree[K, V] {
	if t == nil {
		return nil
	}
	return Meld(t.left, t.right)
}
