package test

import (
	"math"
	"math/rand"
	"testing"

	"github.com/AndrewChon/pqueue"
)

func BenchmarkSkewBinomialPush(b *testing.B) {
	q := pqueue.NewSkewBinomial[int, int]()

	for b.Loop() {
		q.Push(rand.Intn(math.MaxInt64), rand.Intn(math.MaxInt64))
	}
}

func BenchmarkSkewBinomialMeld(b *testing.B) {
	qa := pqueue.NewSkewBinomial[int, int]()
	qb := pqueue.NewSkewBinomial[int, int]()

	for b.Loop() {
		qa.Push(rand.Intn(math.MaxInt64), rand.Intn(math.MaxInt64))
		qb.Push(rand.Intn(math.MaxInt64), rand.Intn(math.MaxInt64))
	}

	b.StartTimer()
	qa.Meld(qb)
	b.StopTimer()

	b.ReportMetric(b.Elapsed().Seconds(), "s/total")
}

func BenchmarkSkewBinomialPop(b *testing.B) {
	q := pqueue.NewSkewBinomial[int, int]()

	for i := 0; i < b.N; i++ {
		q.Push(rand.Intn(math.MaxInt64), rand.Intn(math.MaxInt64))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Pop()
	}
}
