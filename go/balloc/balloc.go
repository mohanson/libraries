package balloc

import (
	"container/list"
	"sync"

	"github.com/mohanson/libraries/go/doa"
)

// Block represents an allocated or free memory block.
type Block struct {
	// Start index of the block.
	Index int
	// Order of the block (size = minBlock * 2^order).
	Order int
	// Space.
	Space []byte
}

// Buddy manages the buddy memory allocation.
type Buddy struct {
	FreeList []*list.List
	MaxOrder int
	MaxTotal int
	MinBlock int
	PreAlloc []byte
	SyncLock *sync.Mutex
}

// Calloc allocates a block of the requested size and ensure it's initialized to zero.
func (b *Buddy) Calloc(size int) Block {
	block := b.Malloc(size)
	for i := range block.Space {
		block.Space[i] = 0
	}
	return block
}

// Balloc allocates a block of the requested size.
func (b *Buddy) Malloc(size int) Block {
	b.SyncLock.Lock()
	defer b.SyncLock.Unlock()
	doa.Doa(size > 0)
	orderReal := log2(b.MinBlock, max(clp2(size), b.MinBlock))
	orderTemp := orderReal
	for orderTemp <= b.MaxOrder && b.FreeList[orderTemp].Len() == 0 {
		orderTemp++
	}
	doa.Doa(orderTemp <= b.MaxOrder)
	blockElem := b.FreeList[orderTemp].Front()
	b.FreeList[orderTemp].Remove(blockElem)
	block := blockElem.Value.(Block)
	for orderTemp > orderReal {
		orderTemp--
		blockSize := b.MinBlock << orderTemp
		b.FreeList[orderTemp].PushFront(Block{
			Index: block.Index + blockSize,
			Order: orderTemp,
		})
		block.Order = orderTemp
	}
	block.Space = b.PreAlloc[block.Index : block.Index+size]
	return block
}

// Free deallocates a block and merges with its buddy if possible.
func (b *Buddy) Free(block Block) {
	b.SyncLock.Lock()
	defer b.SyncLock.Unlock()
	doa.Doa(block.Index >= 0)
	doa.Doa(block.Index < b.MaxTotal)
	doa.Doa(block.Order >= 0)
	doa.Doa(block.Order < b.MaxOrder)
	index := block.Index
	order := block.Order
	for order < b.MaxOrder {
		blockSize := b.MinBlock << order
		buddy := index ^ blockSize
		buddyElem := (*list.Element)(nil)
		for e := b.FreeList[order].Front(); e != nil; e = e.Next() {
			if b := e.Value.(Block); b.Index == buddy {
				buddyElem = e
				break
			}
		}
		if buddyElem == nil {
			break
		}
		b.FreeList[order].Remove(buddyElem)
		index = min(index, buddy)
		order++
	}
	b.FreeList[order].PushFront(Block{Index: index, Order: order})
}

// Get the size of idle memory.
func (b *Buddy) Idle() int {
	b.SyncLock.Lock()
	defer b.SyncLock.Unlock()
	s := 0
	for i, e := range b.FreeList {
		b := b.MinBlock << i
		s += e.Len() * b
	}
	return s
}

// New creates a new buddy allocator.
func New(maxTotal int, minBlock int) *Buddy {
	doa.Doa(isp2(maxTotal))
	doa.Doa(isp2(minBlock))
	doa.Doa(maxTotal > minBlock)
	order := log2(minBlock, maxTotal)
	buddy := &Buddy{
		FreeList: make([]*list.List, order+1),
		MaxOrder: order,
		MaxTotal: maxTotal,
		MinBlock: minBlock,
		PreAlloc: make([]byte, maxTotal),
		SyncLock: &sync.Mutex{},
	}
	for i := range buddy.FreeList {
		buddy.FreeList[i] = list.New()
	}
	buddy.FreeList[order].PushFront(Block{Index: 0, Order: order})
	return buddy
}

func clp2(n int) int {
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n++
	return n
}

func isp2(n int) bool {
	return n > 0 && (n&(n-1)) == 0
}

func log2(m int, n int) int {
	for i := range 64 {
		doa.Doa(m <= n)
		if m == n {
			return i
		}
		m <<= 1
	}
	panic("unreachable")
}
