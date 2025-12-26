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

// Algorithm implements the core buddy allocation algorithm. It manages a memory pool by maintaining free lists for
// different block sizes.
type Algorithm struct {
	// FreeList maintains linked lists of free blocks for each order (size). Each index represents blocks of size
	// MinBlock * 2^index.
	FreeList []int
	// MaxOrder is the maximum order (size class) available, calculated as log2(MaxTotal/MinBlock).
	MaxOrder int
	// MaxTotal is the total size of the memory pool in bytes.
	MaxTotal int
	// MinBlock is the minimum allocation unit size in bytes.
	MinBlock int
	// PreAlloc is the pre-allocated memory pool that the allocator manages.
	PreAlloc []byte
}

// Blockinfo contains information about allocated memory blocks. It describes the location and size of a block within
// the memory pool.
type Blockinfo struct {
	// Offset is the starting position of the allocated block within the memory pool.
	Offset int
	// Length is the size of the allocated block in bytes.
	Length int
}

// Alloc allocates a memory block of the specified order. The order parameter determines the block size as
// MinBlock * 2^order. Returns a Blockinfo with Offset=-1 if allocation fails. If no block of the requested order is
// available, it recursively splits larger blocks.
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

// Avail calculates the total available memory in the pool by summing the sizes of all free blocks across all orders.
func (b *Algorithm) Avail() int {
	s := 0
	for order := 0; order <= b.MaxOrder; order++ {
		n := b.FreeList[order]
		for {
			if n == -1 {
				break
			}
			s += b.MinBlock << order
			n = ildr(b.PreAlloc, n)
		}
	}
	return s
}

// Close frees a memory block and attempts to merge it with its buddy. This implements the buddy merging logic: when
// both buddies are free, they are merged into a larger block at the next order level.
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

// Allocator is a thread-safe wrapper around the buddy allocation algorithm. It provides synchronized access to the
// underlying Algorithm for concurrent use.
type Allocator struct {
	Inner *Algorithm
	Mutex *sync.Mutex
}

// Alloc allocates a byte slice of the requested size from the memory pool. If the pool is exhausted, it falls back to
// heap allocation. This method is thread-safe.
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

// Avail returns the total available memory in the pool. This method is thread-safe.
func (b *Allocator) Avail() int {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()
	return b.Inner.Avail()
}

// Close returns a previously allocated byte slice back to the memory pool. The slice must have been allocated by this
// Allocator's Alloc method. If the slice is not from this pool (e.g., heap-allocated fallback), it's safely ignored.
// This method is thread-safe.
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

// New creates a new buddy allocator with the specified parameters. Parameter minBlock is the minimum allocation unit
// size (must be a power of 2). Parameter maxTotal is the total memory pool size (must be a power of 2).
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
