package test

import (
	"math"
	"math/rand"
	"testing"

	"github.com/AndrewChon/pqueue"
)

func BenchmarkBinaryPush(b *testing.B) {
	q := pqueue.NewBinary[int, int]()

	for b.Loop() {
		q.Push(rand.Intn(math.MaxInt64), rand.Intn(math.MaxInt64))
	}
}

func BenchmarkBinaryMeld(b *testing.B) {
	qa := pqueue.NewBinary[int, int]()
	qb := pqueue.NewBinary[int, int]()

	for b.Loop() {
		qa.Push(rand.Intn(math.MaxInt64), rand.Intn(math.MaxInt64))
		qb.Push(rand.Intn(math.MaxInt64), rand.Intn(math.MaxInt64))
	}

	b.ResetTimer()
	b.StartTimer()
	qa.Meld(qb)
	b.StopTimer()

	b.ReportMetric(b.Elapsed().Seconds(), "s/total")
}

func BenchmarkBinaryPop(b *testing.B) {
	q := pqueue.NewBinary[int, int]()

	for i := 0; i < b.N; i++ {
		q.Push(rand.Intn(math.MaxInt64), rand.Intn(math.MaxInt64))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Pop()
	}
}
