package test

import (
	"math"
	"math/rand"
	"testing"

	"github.com/AndrewChon/pqueue"
)

func BenchmarkCircularPush(b *testing.B) {
	q := pqueue.NewCircularBuffer[int]()

	for b.Loop() {
		q.Push(rand.Intn(math.MaxInt64))
	}
}

func BenchmarkCircularMeld(b *testing.B) {
	qa := pqueue.NewCircularBuffer[int]()
	qb := pqueue.NewCircularBuffer[int]()

	for b.Loop() {
		qa.Push(rand.Intn(math.MaxInt64))
		qb.Push(rand.Intn(math.MaxInt64))
	}

	b.ResetTimer()
	b.StartTimer()
	qa.Meld(qb)
	b.StopTimer()

	b.ReportMetric(b.Elapsed().Seconds(), "s/total")
}

func BenchmarkCircularPop(b *testing.B) {
	q := pqueue.NewCircularBuffer[int]()

	for i := 0; i < b.N; i++ {
		q.Push(rand.Intn(math.MaxInt64))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Pop()
	}
}
