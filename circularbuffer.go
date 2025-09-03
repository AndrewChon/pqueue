package pqueue

import (
	"sync"
	"sync/atomic"
)

var circularBufferIDCounter atomic.Uint64

type node[T any] struct {
	value T
	left  *node[T]
	right *node[T]
}

type CircularBuffer[T any] struct {
	l  sync.RWMutex
	id uint64

	root *node[T]
}

func NewCircularBuffer[T any]() *CircularBuffer[T] {
	return &CircularBuffer[T]{
		root: nil,
		id:   circularBufferIDCounter.Add(1),
	}
}

func (cb *CircularBuffer[T]) Push(v T) {
	cb.l.Lock()
	defer cb.l.Unlock()

	newNode := &node[T]{
		value: v,
	}

	if cb.root == nil {
		newNode.left = newNode
		newNode.right = newNode
		cb.root = newNode
		return
	}

	prevLeft := cb.root.left

	newNode.left = prevLeft
	newNode.right = cb.root

	cb.root.left = newNode
	prevLeft.right = newNode

	cb.root = newNode
}

func (cb *CircularBuffer[T]) Pop() (v T, ok bool) {
	cb.l.Lock()
	defer cb.l.Unlock()

	root, ok := cb.pop()
	if !ok {
		return
	}

	return root.value, ok
}

func (cb *CircularBuffer[T]) pop() (*node[T], bool) {
	if cb.root == nil {
		return nil, false
	}

	minNode := cb.root

	// If there's only one node in the circular buffer.
	if cb.root.left == cb.root && cb.root.right == cb.root {
		cb.root = nil
		return minNode, true
	}

	minNodeLeft := cb.root.left
	minNodeRight := cb.root.right

	minNodeLeft.right = minNodeRight
	minNodeRight.left = minNodeLeft

	cb.root = minNodeRight

	return minNode, true
}

func (cb *CircularBuffer[T]) Peek() T {
	cb.l.RLock()
	defer cb.l.RUnlock()

	if cb.root == nil {
		var zero T
		return zero
	}

	return cb.root.value
}

func (cb *CircularBuffer[T]) Meld(other *CircularBuffer[T]) {
	if cb.id < other.id {
		cb.l.Lock()
		other.l.Lock()
	} else if cb.id > other.id {
		other.l.Lock()
		cb.l.Lock()
	} else {
		panic(ConcurrencySafetyError)
	}

	defer cb.l.Unlock()
	defer other.l.Unlock()

	otherRoot := other.root
	otherLast := otherRoot.left

	cbRoot := cb.root
	cbLast := cbRoot.left

	cbLast.right = otherRoot
	cbRoot.left = otherLast

	other.root = nil
}
