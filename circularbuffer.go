package pqueue

type node[T any] struct {
	value T
	left  *node[T]
	right *node[T]
}

type CircularBuffer[T any] struct {
	root *node[T]
}

func NewCircularBuffer[T any]() *CircularBuffer[T] {
	return &CircularBuffer[T]{
		root: nil,
	}
}

func (cb *CircularBuffer[T]) Push(v T) {
	newNode := &node[T]{
		value: v,
	}
	newNode.left = newNode
	newNode.right = newNode

	if cb.root == nil {
		cb.root = newNode
		return
	}

	prevLeft := cb.root.left
	cb.root.left = newNode
	prevLeft.right = newNode
}

func (cb *CircularBuffer[T]) Pop() (v T, ok bool) {
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

	root := cb.root
	cb.root.left = cb.root.right
	cb.root = cb.root.right

	return root, true
}

func (cb *CircularBuffer[T]) Peek() T {
	if cb.root == nil {
		var zero T
		return zero
	}

	return cb.root.value
}

func (cb *CircularBuffer[T]) Meld(other *CircularBuffer[T]) {
	otherRoot := other.root
	otherLast := otherRoot.left

	cbRoot := cb.root
	cbLast := cbRoot.left

	cbLast.right = otherRoot
	cbRoot.left = otherLast

	other.root = nil
}
