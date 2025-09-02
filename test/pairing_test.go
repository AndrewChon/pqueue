package test

import (
	"math"
	"math/rand"
	"testing"

	"github.com/AndrewChon/pqueue"
)

func BenchmarkPairingPush(b *testing.B) {
	q := pqueue.NewPairing[int, int]()

	for b.Loop() {
		q.Push(rand.Intn(math.MaxInt64), rand.Intn(math.MaxInt64))
	}
}

func BenchmarkPairingMeld(b *testing.B) {
	qa := pqueue.NewPairing[int, int]()
	qb := pqueue.NewPairing[int, int]()

	for b.Loop() {
		qa.Push(rand.Intn(math.MaxInt64), rand.Intn(math.MaxInt64))
		qb.Push(rand.Intn(math.MaxInt64), rand.Intn(math.MaxInt64))
	}

	b.StartTimer()
	qa.Meld(qb)
	b.StopTimer()

	b.ReportMetric(b.Elapsed().Seconds(), "s/total")
}

func BenchmarkPairingPop(b *testing.B) {
	q := pqueue.NewPairing[int, int]()

	for i := 0; i < b.N; i++ {
		q.Push(rand.Intn(math.MaxInt64), rand.Intn(math.MaxInt64))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Pop()
	}
}
