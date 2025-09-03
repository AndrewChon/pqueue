// Package pqueue provides concurrency-safe priority queue implementations.
package pqueue

import (
	"errors"
)

var (
	ConcurrencySafetyError = errors.New("concurrency-safety error: one or more queues share the same underlying" +
		"ID. ensure that all queues are being created via their designated constructors")
)

type CrossMeldable interface {
	CrossMeld(other CrossMeldable)
}
