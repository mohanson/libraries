package balloc

import (
	"math/rand/v2"
	"testing"

	"github.com/mohanson/libraries/go/doa"
)

func TestBalloc(t *testing.T) {
	maxAlive := 1024
	maxAlloc := 1024
	maxRange := 1024 * 1024
	maxTotal := 1024 * 1024
	minBlock := 64
	balloc := New(maxTotal, minBlock)
	record := []Block{}
	for range maxRange {
		i := rand.Int() % maxAlive
		if i < len(record) {
			balloc.Free(record[i])
			record = append(record[:i], record[i+1:]...)
		}
		if i > len(record) {
			record = append(record, balloc.Malloc(max(1, rand.Int()%maxAlloc)))
		}
	}
	for _, e := range record {
		balloc.Free(e)
	}
	doa.Doa(balloc.Idle() == maxTotal)
}
