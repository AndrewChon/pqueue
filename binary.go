package pqueue

import (
	"cmp"
	"sync"
	"sync/atomic"

	"github.com/AndrewChon/pqueue/binary"
)

var binaryIDCounter atomic.Uint64

// Binary is a concurrency-safe, min-priority queue built on a binary heap.
type Binary[K cmp.Ordered, V any] struct {
	// A locking order needs to be defined and strictly followed for safety; thus, we do not want to expose the mutex.
	l  sync.RWMutex
	id uint64

	heap *binary.Heap[K, V]
}

func NewBinary[K cmp.Ordered, V any]() *Binary[K, V] {
	return &Binary[K, V]{
		id:   binaryIDCounter.Add(1),
		heap: binary.NewHeap[K, V](),
	}
}

func (b *Binary[K, V]) Size() int {
	b.l.RLock()
	defer b.l.RUnlock()

	return b.heap.Size()
}

func (b *Binary[K, V]) Clear() {
	b.l.Lock()
	defer b.l.Unlock()

	b.heap.Clear()
}

func (b *Binary[K, V]) Peek() V {
	b.l.RLock()
	defer b.l.RUnlock()

	minNode := b.heap.FindMin()
	if minNode == nil {
		var zero V
		return zero
	}

	return minNode.Value()
}

func (b *Binary[K, V]) Pop() (v V, ok bool) {
	b.l.Lock()
	defer b.l.Unlock()

	n := b.heap.FindMin()
	if n == nil {
		return
	}

	v = n.Value()
	b.heap.RemoveMin()
	return v, true
}

func (b *Binary[K, V]) Push(v V, priority K) {
	b.l.Lock()
	defer b.l.Unlock()

	newNode := binary.NewNode(priority, v)
	b.heap.Insert(newNode)
}

func (b *Binary[K, V]) Meld(other *Binary[K, V]) {
	if b.id < other.id {
		b.l.Lock()
		other.l.Lock()
	} else if b.id > other.id {
		other.l.Lock()
		b.l.Lock()
	} else {
		panic(ConcurrencySafetyError)
	}

	defer b.l.Unlock()
	defer other.l.Unlock()

	b.heap = binary.Merge(b.heap, other.heap)
	other.heap.Clear()
}
