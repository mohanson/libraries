package balloc

import (
	"log"
	"math/bits"
	"sync"
	"unsafe"
)

// A hand written buddy allocator implementation for high performance requirements environments.
//
// This module implements a buddy memory allocator that manages a fixed-size memory pool.
// The buddy allocator works by recursively splitting memory blocks into pairs of equal-sized "buddies" until it finds
// a block of the appropriate size. When memory is freed, it attempts to merge adjacent buddy blocks back together to
// reduce fragmentation.

// Buddy allocation algorithm implementation.
type Algorithm struct {
	FreeList []int
	MaxOrder int
	MaxTotal int
	MinBlock int
	PreAlloc []byte
}

// Information about allocated memory blocks.
type Blockinfo struct {
	// Offset of the allocated block within the memory pool.
	Offset int
	// Length of the allocated block.
	Length int
}

func (b *Algorithm) Alloc(order int) Blockinfo {
	if order > b.MaxOrder {
		return Blockinfo{Offset: -1, Length: 0}
	}
	blockSize := b.MinBlock << order
	if b.FreeList[order] != -1 {
		blockOffset := b.FreeList[order]
		b.FreeList[order] = ildr(b.PreAlloc, blockOffset)
		return Blockinfo{
			Offset: blockOffset,
			Length: blockSize,
		}
	}
	block := b.Alloc(order + 1)
	blockOffset := block.Offset
	if blockOffset == -1 {
		return block
	}
	buddyOffset := blockOffset + blockSize
	istr(b.PreAlloc, buddyOffset, -1)
	b.FreeList[order] = buddyOffset
	return Blockinfo{
		Offset: block.Offset,
		Length: blockSize,
	}
}

func (b *Algorithm) Close(block Blockinfo) {
	order := log2(b.MinBlock, block.Length)
	blockOffset := block.Offset
	blockIdx := blockOffset / block.Length
	buddyIdx := blockIdx ^ 1
	buddyOffset := buddyIdx * block.Length
	buddy := Blockinfo{Offset: buddyOffset, Length: block.Length}
	upper := Blockinfo{Offset: min(blockOffset, buddyOffset), Length: block.Length * 2}
	n := b.FreeList[order]
	m := 0
	for {
		if n == -1 {
			istr(b.PreAlloc, blockOffset, b.FreeList[order])
			b.FreeList[order] = blockOffset
			break
		}
		m = ildr(b.PreAlloc, n)
		if n == buddy.Offset {
			b.FreeList[order] = m
			b.Close(upper)
			break
		}
		if m == buddy.Offset {
			istr(b.PreAlloc, n, ildr(b.PreAlloc, m))
			b.Close(upper)
			break
		}
		n = m
	}
}

// Buddy allocation.
type Allocator struct {
	Inner *Algorithm
	Mutex *sync.Mutex
}

func (b *Allocator) Alloc(size int) []byte {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()
	order := log2(b.Inner.MinBlock, max(b.Inner.MinBlock, npo2(size)))
	block := b.Inner.Alloc(order)
	if block.Offset == -1 {
		log.Println("balloc: out of memory")
		return make([]byte, size)
	}
	return b.Inner.PreAlloc[block.Offset : block.Offset+size]
}

func (b *Allocator) Close(data []byte) {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()
	blockOffset := int(uintptr(unsafe.Pointer(&data[0])) - uintptr(unsafe.Pointer(&b.Inner.PreAlloc[0])))
	if blockOffset < 0 || blockOffset >= b.Inner.MaxTotal {
		return
	}
	b.Inner.Close(Blockinfo{
		Offset: blockOffset,
		Length: max(b.Inner.MinBlock, npo2(len(data))),
	})
}

// Function ildr reads an int value from byte slice m at offset o.
func ildr(m []byte, o int) int {
	return *(*int)(unsafe.Pointer(&m[o]))
}

// Function istr saves an int value into byte slice m at offset o.
func istr(m []byte, o int, v int) {
	b := (*(*[unsafe.Sizeof(v)]byte)(unsafe.Pointer(&v)))[:]
	copy(m[o:], b)
}

// Function isp2 checks if n is a power of 2.
func isp2(n int) bool {
	return n > 0 && (n&(n-1)) == 0
}

// Function log2 computes log2 of n/m where m and n are powers of 2 and m <= n.
func log2(m int, n int) int {
	return bits.Len(uint(n)) - bits.Len(uint(m))
}

// Function npo2 computes the next power of 2 greater than or equal to n.
func npo2(n int) int {
	return 1 << (bits.Len(uint(n - 1)))
}

// New creates a new buddy allocator.
func New(minBlock int, maxTotal int) *Allocator {
	if !isp2(minBlock) {
		log.Panicln("balloc: max block is not a power of 2")
	}
	if !isp2(maxTotal) {
		log.Panicln("balloc: max total is not a power of 2")
	}
	order := log2(minBlock, maxTotal)
	inner := &Algorithm{
		FreeList: make([]int, order+1),
		MaxOrder: order,
		MaxTotal: maxTotal,
		MinBlock: minBlock,
		PreAlloc: make([]byte, maxTotal),
	}
	for i := range inner.FreeList {
		inner.FreeList[i] = -1
	}
	inner.FreeList[order] = 0
	istr(inner.PreAlloc, 0, -1)
	return &Allocator{
		Inner: inner,
		Mutex: &sync.Mutex{},
	}
}
