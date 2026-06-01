// Package once provides a generic once type. It allows you to lazily initialize a variable in a concurrently safe way.
package once

import (
	"sync"
)

// Once is a wrapper around an initialization function, ensuring it happens only once, even in a concurrent environment.
type OnceNew[T any] struct {
	// The initialization function to create the object of type T.
	init func() T
	// The initialized object of type T.
	inst T
	// A mutex that ensures the initialization only happens once.
	once sync.Once
}

// Do returns the initialized object, creating it if necessary.
func (s *OnceNew[T]) Do() T {
	s.once.Do(func() {
		s.inst = s.init()
	})
	return s.inst
}

// NewOnceNew creates a new OnceNew wrapper around an initialization function.
func NewOnceNew[T any](f func() T) *OnceNew[T] {
	return &OnceNew[T]{
		init: f,
		once: sync.Once{},
	}
}

// OnceErr is an object that will only store an error once.
type OnceErr struct {
	mux *sync.Mutex // Guards following
	err error
	sig chan struct{}
}

// Get an error from OnceErr.
func (e *OnceErr) Get() error {
	e.mux.Lock()
	defer e.mux.Unlock()
	return e.err
}

// Put an error into OnceErr.
func (e *OnceErr) Put(err error) {
	e.mux.Lock()
	defer e.mux.Unlock()
	if e.err != nil {
		return
	}
	e.err = err
	close(e.sig)
}

// When any error puts, the sig will be sent.
func (e *OnceErr) Sig() <-chan struct{} {
	return e.sig
}

// NewOnceErr creates a new OnceErr.
func NewOnceErr() *OnceErr {
	return &OnceErr{
		mux: &sync.Mutex{},
		err: nil,
		sig: make(chan struct{}),
	}
}
