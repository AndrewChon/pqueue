package pqueue

import (
	"cmp"
	"sync"
	"sync/atomic"

	"github.com/AndrewChon/pqueue/skewbinomial"
)

var skewBinomialIDCounter atomic.Uint64

// SkewBinomial is a concurrency-safe, min-priority queue built on a skew binomial heap.
type SkewBinomial[K cmp.Ordered, V any] struct {
	l  sync.RWMutex
	id uint64

	heap *skewbinomial.Forest[K, V]
	size int
}

func NewSkewBinomial[K cmp.Ordered, V any]() *SkewBinomial[K, V] {
	return &SkewBinomial[K, V]{
		id:   skewBinomialIDCounter.Add(1),
		heap: skewbinomial.NewForest[K, V](),
		size: 0,
	}
}

func (sb *SkewBinomial[K, V]) Size() int {
	sb.l.RLock()
	defer sb.l.RUnlock()

	return sb.size
}

func (sb *SkewBinomial[K, V]) Clear() {
	sb.l.Lock()
	defer sb.l.Unlock()

	sb.heap = skewbinomial.NewForest[K, V]()
	sb.size = 0
}

func (sb *SkewBinomial[K, V]) Peek() V {
	sb.l.RLock()
	defer sb.l.RUnlock()

	minTree, _ := sb.heap.FindMin()
	if minTree == nil {
		var zero V
		return zero
	}

	return minTree.Value()
}

func (sb *SkewBinomial[K, V]) Pop() (v V, ok bool) {
	sb.l.Lock()
	defer sb.l.Unlock()

	minTree, i := sb.heap.FindMin()
	if minTree == nil {
		return
	}

	v = minTree.Value()

	sb.heap.Remove(minTree, i)

	sb.size--
	return v, true
}

func (sb *SkewBinomial[K, V]) Push(v V, priority K) {
	sb.l.Lock()
	defer sb.l.Unlock()

	sb.heap.Insert(priority, v)
	sb.size++
}

func (sb *SkewBinomial[K, V]) Meld(other *SkewBinomial[K, V]) {
	if sb.id < other.id {
		sb.l.Lock()
		other.l.Lock()
	} else if sb.id > other.id {
		other.l.Lock()
		sb.l.Lock()
	} else {
		panic(ConcurrencySafetyError)
	}

	defer sb.l.Unlock()
	defer other.l.Unlock()

	sb.heap.Merge(other.heap)
	sb.size += other.size

	other.heap = skewbinomial.NewForest[K, V]()
	other.size = 0
}
