package balloc

import (
	"math/rand/v2"
	"testing"
)

func TestFuzz(t *testing.T) {
	maxAlive := 1024
	maxAlloc := 2048
	maxRange := 1024 * 1024
	maxTotal := 1024 * 1024
	minBlock := 64
	balloc := New(minBlock, maxTotal)
	record := [][]byte{}
	for range maxRange {
		actionRandom := rand.Int() % maxAlive
		action := 0
		if len(record) > actionRandom {
			action = 1
		}
		switch action {
		case 0:
			record = append(record, balloc.Alloc(max(1, rand.Int()%maxAlloc)))
		case 1:
			i := rand.Int() % len(record)
			balloc.Close(record[i])
			record = append(record[:i], record[i+1:]...)
		}
	}
	for _, e := range record {
		balloc.Close(e)
	}
	for i := range balloc.Inner.MaxOrder {
		if balloc.Inner.FreeList[i] != -1 {
			t.FailNow()
		}
	}
	if balloc.Inner.FreeList[balloc.Inner.MaxOrder] != 0 {
		t.FailNow()
	}
}
