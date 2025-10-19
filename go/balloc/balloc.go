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
	freeList []*list.List
	maxOrder int
	maxTotal int
	minBlock int
	preAlloc []byte
	syncLock *sync.Mutex
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
	b.syncLock.Lock()
	defer b.syncLock.Unlock()
	doa.Doa(size > 0)
	orderReal := log2(b.minBlock, max(clp2(size), b.minBlock))
	orderTemp := orderReal
	for orderTemp <= b.maxOrder && b.freeList[orderTemp].Len() == 0 {
		orderTemp++
	}
	doa.Doa(orderTemp <= b.maxOrder)
	blockElem := b.freeList[orderTemp].Front()
	b.freeList[orderTemp].Remove(blockElem)
	block := blockElem.Value.(Block)
	for orderTemp > orderReal {
		orderTemp--
		blockSize := b.minBlock << orderTemp
		b.freeList[orderTemp].PushFront(Block{
			Index: block.Index + blockSize,
			Order: orderTemp,
		})
		block.Order = orderTemp
	}
	block.Space = b.preAlloc[block.Index : block.Index+size]
	return block
}

// Free deallocates a block and merges with its buddy if possible.
func (b *Buddy) Free(block Block) {
	b.syncLock.Lock()
	defer b.syncLock.Unlock()
	doa.Doa(block.Index >= 0)
	doa.Doa(block.Index < b.maxTotal)
	doa.Doa(block.Order >= 0)
	doa.Doa(block.Order < b.maxOrder)
	index := block.Index
	order := block.Order
	for order < b.maxOrder {
		blockSize := b.minBlock << order
		buddy := index ^ blockSize
		buddyElem := (*list.Element)(nil)
		for e := b.freeList[order].Front(); e != nil; e = e.Next() {
			if b := e.Value.(Block); b.Index == buddy {
				buddyElem = e
				break
			}
		}
		if buddyElem == nil {
			break
		}
		b.freeList[order].Remove(buddyElem)
		index = min(index, buddy)
		order++
	}
	b.freeList[order].PushFront(Block{Index: index, Order: order})
}

// Get the size of idle memory.
func (b *Buddy) Idle() int {
	b.syncLock.Lock()
	defer b.syncLock.Unlock()
	s := 0
	for i, e := range b.freeList {
		b := b.minBlock << i
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
		freeList: make([]*list.List, order+1),
		maxOrder: order,
		maxTotal: maxTotal,
		minBlock: minBlock,
		preAlloc: make([]byte, maxTotal),
		syncLock: &sync.Mutex{},
	}
	for i := range buddy.freeList {
		buddy.freeList[i] = list.New()
	}
	buddy.freeList[order].PushFront(Block{Index: 0, Order: order})
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
