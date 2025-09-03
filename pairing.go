package pqueue

import (
	"cmp"
	"sync"
	"sync/atomic"

	"github.com/AndrewChon/pqueue/pairing"
)

var pairingIDCounter atomic.Uint64

// Pairing is a concurrency-safe, min-priority queue built on a pairing heap.
type Pairing[K cmp.Ordered, V any] struct {
	// A locking order needs to be defined and strictly followed for safety; thus, we do not want to expose the mutex.
	l  sync.RWMutex
	id uint64

	root *pairing.Tree[K, V]
	size int
}

func NewPairing[K cmp.Ordered, V any]() *Pairing[K, V] {
	return &Pairing[K, V]{
		id:   pairingIDCounter.Add(1),
		root: nil,
		size: 0,
	}
}

func (p *Pairing[K, V]) Size() int {
	p.l.RLock()
	defer p.l.RUnlock()

	return p.size
}

func (p *Pairing[K, V]) Clear() {
	p.l.Lock()
	defer p.l.Unlock()

	p.root = nil
	p.size = 0
}

func (p *Pairing[K, V]) Peek() V {
	p.l.RLock()
	defer p.l.RUnlock()

	minNode := pairing.FindMin(p.root)
	if minNode == nil {
		var zero V
		return zero
	}

	return minNode.Value()
}

func (p *Pairing[K, V]) Pop() (v V, ok bool) {
	p.l.Lock()
	defer p.l.Unlock()

	t := pairing.FindMin(p.root)
	if t == nil {
		return
	}

	v = t.Value()
	p.root = pairing.RemoveMin(p.root)
	p.size--

	return v, true
}

func (p *Pairing[K, V]) Push(v V, priority K) {
	p.l.Lock()
	defer p.l.Unlock()

	newNode := pairing.NewTree(priority, v)
	p.root = pairing.Insert(p.root, newNode)

	p.size++
}

// Meld merges another Pairing queue into this one and clears it.
func (p *Pairing[K, V]) Meld(other *Pairing[K, V]) {
	if p.id < other.id {
		p.l.Lock()
		other.l.Lock()
	} else if p.id > other.id {
		other.l.Lock()
		p.l.Lock()
	} else {
		panic(ConcurrencySafetyError)
	}

	defer p.l.Unlock()
	defer other.l.Unlock()

	p.root = pairing.Meld(p.root, other.root)
	p.size += other.size

	other.root = nil
	other.size = 0
}
