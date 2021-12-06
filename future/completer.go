package future

import "sync/atomic"

// a completer. Usually used to notify the completion of asynchronous operations
// like dart Completer<T>
type Completer struct {
	future    *Future
	completed int32
}

func New() *Completer {
	return &Completer{
		future: &Future{
			c: make(chan struct{}),
		},
	}
}

// Returns the associated Future
func (c *Completer) Future() *Future {
	return c.future
}

// Completes future with the supplied values.
//
// Complete and CompleteError are only valid at the first call and return true, and subsequent calls will return false
//
// Complete and CompleteError can be safely called in multiple goroutine environments but it is impossible to determine which goroutine call will be executed at first (the value that is executed first will notify the future)
func (c *Completer) Complete(value interface{}) bool {
	if c.completed == 0 && atomic.SwapInt32(&c.completed, 1) == 0 {
		c.future.value = value
		close(c.future.c)
		return true
	}
	return false
}

// Complete future with an error
//
// Complete and CompleteError are only valid at the first call and return true, and subsequent calls will return false
//
// Complete and CompleteError can be safely called in multiple goroutine environments but it is impossible to determine which goroutine call will be executed at first (the value that is executed first will notify the future)
func (c *Completer) CompleteError(err error) bool {
	if c.completed == 0 && atomic.SwapInt32(&c.completed, 1) == 0 {
		c.future.err = err
		close(c.future.c)
		return true
	}
	return false
}
