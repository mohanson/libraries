// Package once provides concurrency-safe lazy initialization helpers.
package once

import (
	"sync"
)

// OnceNew wraps an initialization function and guarantees it is called at most once.
type OnceNew[T any] struct {
	init func() T
	inst T
	once sync.Once
}

// Do returns the value produced by the initialization function, calling it on the first invocation.
func (s *OnceNew[T]) Do() T {
	s.once.Do(func() {
		s.inst = s.init()
	})
	return s.inst
}

// NewOnceNew returns an OnceNew that will call f to initialize its value on first use.
func NewOnceNew[T any](f func() T) *OnceNew[T] {
	return &OnceNew[T]{init: f}
}

// OnceErr stores the first non-nil error written to it and signals all waiters via a channel.
type OnceErr struct {
	mux sync.Mutex
	err error
	sig chan struct{}
}

// Get returns the stored error, or nil if no error has been set.
func (e *OnceErr) Get() error {
	e.mux.Lock()
	defer e.mux.Unlock()
	return e.err
}

// Put stores err if no error has been stored yet, then closes the signal channel.
func (e *OnceErr) Put(err error) {
	e.mux.Lock()
	defer e.mux.Unlock()
	if e.err != nil {
		return
	}
	e.err = err
	close(e.sig)
}

// Sig returns a channel that is closed when the first error is stored.
func (e *OnceErr) Sig() <-chan struct{} {
	return e.sig
}

// NewOnceErr returns a new OnceErr.
func NewOnceErr() *OnceErr {
	return &OnceErr{sig: make(chan struct{})}
}
