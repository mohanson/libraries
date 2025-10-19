// Package gool offers a high-level API for running tasks asynchronously, restricting concurrent executions to the
// number of cpu cores or a specified limit.
package gool

import (
	"runtime"
	"sync"
)

var cCpu = make(chan struct{}, runtime.NumCPU())

// Gool manages a pool of goroutines for asynchronous task execution.
type Gool struct {
	c chan struct{}
	m *sync.Mutex
	w *sync.WaitGroup
}

// Call submits a function f for asynchronous execution in a new goroutine, respecting the concurrency limit.
func (g *Gool) Call(f func()) {
	g.c <- struct{}{}
	g.w.Add(1)
	go func() {
		f()
		g.w.Done()
		<-g.c
	}()
}

// Lock executes function f with exclusive access, synchronizing via the mutex, typically for aggregating results.
func (g *Gool) Lock(f func()) {
	g.m.Lock()
	defer g.m.Unlock()
	f()
}

// Wait blocks until all submitted tasks have completed.
func (g *Gool) Wait() {
	g.w.Wait()
}

// Cpu initializes a Gool instance with a global concurrency limit specified by cpu cores.
func Cpu() *Gool {
	return &Gool{
		c: cCpu,
		m: &sync.Mutex{},
		w: &sync.WaitGroup{},
	}
}

// New initializes a Gool instance with a custom concurrency limit specified by n.
func New(n int) *Gool {
	return &Gool{
		c: make(chan struct{}, n),
		m: &sync.Mutex{},
		w: &sync.WaitGroup{},
	}
}
