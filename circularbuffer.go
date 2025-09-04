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

	size int
	root *node[T]
}

func NewCircularBuffer[T any]() *CircularBuffer[T] {
	return &CircularBuffer[T]{
		root: nil,
		id:   circularBufferIDCounter.Add(1),
	}
}

func (cb *CircularBuffer[T]) Size() int {
	cb.l.RLock()
	defer cb.l.RUnlock()

	return cb.size
}

func (cb *CircularBuffer[T]) Clear() {
	cb.l.Lock()
	defer cb.l.Unlock()

	cb.root = nil
	cb.size = 0
}

func (cb *CircularBuffer[T]) Push(v T) {
	cb.l.Lock()
	defer cb.l.Unlock()

	newNode := &node[T]{
		value: v,
	}

	// If the buffer is empty, simply set cb.root to newNode.
	if cb.root == nil {
		newNode.left = newNode
		newNode.right = newNode
		cb.root = newNode
		cb.size = 1
		return
	}

	tail := cb.root.left

	// Detach the tail of the buffer.
	// tail.right = nil
	// cb.root.left = nil

	// Link newNode to the right of the old tail.
	tail.right = newNode
	cb.root.left = newNode
	newNode.left = tail
	newNode.right = cb.root

	cb.size++
}

func (cb *CircularBuffer[T]) Pop() (v T, ok bool) {
	cb.l.Lock()
	defer cb.l.Unlock()

	if cb.root == nil {
		var zero T
		return zero, false
	}

	minNode := cb.root

	// If cb is a singleton, simply set cb.root to nil.
	if minNode.left == minNode && minNode.right == minNode {
		cb.root = nil
		cb.size = 0
		return minNode.value, true
	}

	next := minNode.right
	last := minNode.left

	// Detach minNode from the rest of the buffer.
	// minNode.left = nil
	// minNode.right = nil

	// Connect the next node in line to the last node.
	next.left = last
	last.right = next

	// Update cb.root accordingly.
	cb.root = next

	cb.size--

	return minNode.value, true
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

	defer func() {
		cb.l.Unlock()
		other.l.Unlock()
	}()

	if other.root == nil {
		return
	}

	if cb.root == nil {
		cb.root = other.root
		cb.size = other.size

		other.root = nil
		other.size = 0
		return
	}

	cbLast := cb.root.left
	otherMinNode := other.root
	otherLast := otherMinNode.left

	// Detach otherMinNode from its tail.
	// otherMinNode.left = nil

	// Detach cbLast from its head.
	// cbLast.right = nil

	// Connect cbLast to otherMinNode.
	cbLast.right = otherMinNode
	otherMinNode.left = cbLast

	// Connect otherLast to minNode.
	otherLast.right = cb.root
	cb.root.left = otherLast

	cb.size += other.size

	// Clear other.
	other.root = nil
	other.size = 0
}
