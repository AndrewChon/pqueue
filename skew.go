package pqueue

import (
	"cmp"
	"sync"
	"sync/atomic"

	"github.com/AndrewChon/pqueue/skew"
)

var skewIDCounter atomic.Uint64

// Skew is a concurrency-safe, min-priority queue built on a skew heap.
type Skew[K cmp.Ordered, V any] struct {
	// A locking order needs to be defined and strictly followed for safety; thus, we do not want to expose the mutex.
	l  sync.RWMutex
	id uint64

	root *skew.Tree[K, V]
	size int
}

func NewSkew[K cmp.Ordered, V any]() *Skew[K, V] {
	return &Skew[K, V]{
		id:   skewIDCounter.Add(1),
		root: nil,
		size: 0,
	}
}

func (s *Skew[K, V]) Size() int {
	s.l.RLock()
	defer s.l.RUnlock()

	return s.size
}

func (s *Skew[K, V]) Clear() {
	s.l.Lock()
	defer s.l.Unlock()

	s.root = nil
	s.size = 0
}

func (s *Skew[K, V]) Peek() V {
	s.l.RLock()
	defer s.l.RUnlock()

	minNode := skew.FindMin(s.root)
	if minNode == nil {
		var zero V
		return zero
	}

	return minNode.Value()
}

func (s *Skew[K, V]) Pop() (v V, ok bool) {
	s.l.Lock()
	defer s.l.Unlock()

	t := skew.FindMin(s.root)
	if t == nil {
		return
	}

	v = t.Value()
	s.root = skew.RemoveMin(s.root)
	s.size--

	return v, true
}

func (s *Skew[K, V]) Push(v V, priority K) {
	s.l.Lock()
	defer s.l.Unlock()

	newNode := skew.NewTree(priority, v)
	s.root = skew.Insert(s.root, newNode)

	s.size++
}

// Meld merges another Skew queue into this one and clears it.
func (s *Skew[K, V]) Meld(other *Skew[K, V]) {
	if s.id < other.id {
		s.l.Lock()
		other.l.Lock()
	} else if s.id > other.id {
		other.l.Lock()
		s.l.Lock()
	} else {
		panic(ConcurrencySafetyError)
	}

	defer s.l.Unlock()
	defer other.l.Unlock()

	s.root = skew.Meld(s.root, other.root)
	s.size += other.size

	other.root = nil
	other.size = 0
}
