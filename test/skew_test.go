package test

import (
	"math"
	"math/rand"
	"testing"

	"github.com/AndrewChon/pqueue"
)

func BenchmarkSkewPush(b *testing.B) {
	q := pqueue.NewSkew[int, int]()

	for b.Loop() {
		q.Push(rand.Intn(math.MaxInt64), rand.Intn(math.MaxInt64))
	}
}

func BenchmarkSkewMeld(b *testing.B) {
	qa := pqueue.NewSkew[int, int]()
	qb := pqueue.NewSkew[int, int]()

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

func BenchmarkSkewPop(b *testing.B) {
	q := pqueue.NewSkew[int, int]()

	for i := 0; i < b.N; i++ {
		q.Push(rand.Intn(math.MaxInt64), rand.Intn(math.MaxInt64))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Pop()
	}
}
