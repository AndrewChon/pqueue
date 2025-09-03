package binary

import (
	"cmp"
)

type Node[K cmp.Ordered, V any] struct {
	key   K
	value V
}

func NewNode[K cmp.Ordered, V any](key K, value V) *Node[K, V] {
	return &Node[K, V]{
		key:   key,
		value: value,
	}
}

func (n *Node[K, V]) Key() K {
	return n.key
}

func (n *Node[K, V]) Value() V {
	return n.value
}

type Heap[K cmp.Ordered, V any] struct {
	array []*Node[K, V]
}

func NewHeap[K cmp.Ordered, V any]() *Heap[K, V] {
	return &Heap[K, V]{
		array: make([]*Node[K, V], 0),
	}
}

func (h *Heap[K, V]) Size() int {
	return len(h.array)
}

func (h *Heap[K, V]) Clear() {
	h.array = make([]*Node[K, V], 0)
}

func (h *Heap[K, V]) FindMin() *Node[K, V] {
	if len(h.array) == 0 {
		return nil
	}
	return h.array[0]
}

func Merge[K cmp.Ordered, V any](a, b *Heap[K, V]) *Heap[K, V] {
	newHeap := NewHeap[K, V]()

	sizeA := len(a.array)
	sizeB := len(b.array)

	if sizeA == 0 && sizeB == 0 {
		return newHeap
	} else if sizeA == 0 {
		newHeap.array = append(newHeap.array, b.array...)
		return newHeap
	} else if sizeB == 0 {
		newHeap.array = append(newHeap.array, a.array...)
		return newHeap
	}

	newHeap.array = append(newHeap.array, a.array...)
	newHeap.array = append(newHeap.array, b.array...)

	// Heapify down from the bottom up.
	size := len(newHeap.array)
	for i := size/2 - 1; i >= 0; i-- {
		newHeap.heapifyDown(i)
	}

	return newHeap
}

func (h *Heap[K, V]) Insert(newNode *Node[K, V]) {
	h.array = append(h.array, newNode)
	newNodeIndex := len(h.array) - 1

	h.heapifyUp(newNodeIndex)
}

func (h *Heap[K, V]) RemoveMin() {
	size := len(h.array)

	if size == 0 {
		return
	}

	h.array[0] = h.array[size-1]
	h.array = h.array[:size-1]
	h.heapifyDown(0)
}

func (h *Heap[K, V]) heapifyUp(i int) {
	if i <= 0 {
		return
	}

	parentIndex := (i - 1) / 2
	if h.array[i].key < h.array[parentIndex].key {
		h.array[i], h.array[parentIndex] = h.array[parentIndex], h.array[i]
		h.heapifyUp(parentIndex)
	}
}

func (h *Heap[K, V]) heapifyDown(i int) {
	size := len(h.array)

	if i >= size || i < 0 {
		return
	}

	leftChild := 2*i + 1
	rightChild := 2*i + 2
	smallest := i

	if leftChild < size && h.array[leftChild].key < h.array[smallest].key {
		smallest = leftChild
	}

	if rightChild < size && h.array[rightChild].key < h.array[smallest].key {
		smallest = rightChild
	}

	if smallest != i {
		h.array[i], h.array[smallest] = h.array[smallest], h.array[i]
		h.heapifyDown(smallest)
	}
}
