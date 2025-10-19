// Package once provides a generic once type. It allows you to lazily initialize a variable in a concurrently safe way.
package once

import (
	"sync"
)

// Once is a wrapper around an initialization function, ensuring it happens only once, even in a concurrent environment.
type Once[T any] struct {
	// The initialization function to create the object of type T.
	init func() T
	// The initialized object of type T.
	inst T
	// A mutex that ensures the initialization only happens once.
	once sync.Once
}

// Do returns the initialized object, creating it if necessary.
func (s *Once[T]) Do() T {
	s.once.Do(func() {
		s.inst = s.init()
	})
	return s.inst
}

// New creates a new Once wrapper around an initialization function.
func New[T any](f func() T) *Once[T] {
	return &Once[T]{
		init: f,
		once: sync.Once{},
	}
}
