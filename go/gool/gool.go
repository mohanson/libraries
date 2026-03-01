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
	i int
	m []*sync.Mutex
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

// Done waits for all previously submitted tasks to finish their result processing, then calls f, and finally signals
// the next task that it may proceed. Together with Next this forms a mutex chain: task i blocks until task i-1 has
// called Done, guaranteeing that f is invoked in the same order as the corresponding Next calls, while the heavy
// computation still runs in parallel.
func (g *Gool) Done(i int, f func()) {
	g.m[i].Lock()
	f()
	g.m[(i+1)%len(g.m)].Unlock()
}

// Next returns a sequential index (token) for the next task. Call Next before launching each goroutine to assign an
// ordering slot. The returned index must be passed to Done so that result callbacks are executed in submission order
// even though the underlying tasks run concurrently.
func (g *Gool) Next() int {
	t := g.i
	g.i += 1
	g.i %= len(g.m)
	return t
}

// Wait blocks until all submitted tasks have completed.
func (g *Gool) Wait() {
	g.w.Wait()
}

// Cpu initializes a Gool instance with a global concurrency limit specified by cpu cores.
func Cpu() *Gool {
	g := New(runtime.NumCPU())
	g.c = cCpu
	return g
}

// New initializes a Gool instance with a custom concurrency limit specified by n.
func New(n int) *Gool {
	m := make([]*sync.Mutex, n)
	for i := range m {
		m[i] = &sync.Mutex{}
		m[i].Lock()
	}
	m[0].Unlock()
	return &Gool{
		c: make(chan struct{}, n),
		i: 0,
		m: m,
		w: &sync.WaitGroup{},
	}
}
