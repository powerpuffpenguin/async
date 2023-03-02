package future

import "sync"

// Multiple concurrent requests with the same request (the return value is also considered consistent),
// they will share the result of the same request
type Once[T any] struct {
	future *Future[T]
	f      func() (value T, e error)
	rw     sync.RWMutex
}

func NewOnce[T any](f func() (value T, e error)) *Once[T] {
	return &Once[T]{
		f: f,
	}
}
func (o *Once[T]) Do() *Future[T] {
	// executing
	o.rw.RLock()
	future := o.future
	o.rw.RUnlock()
	if future != nil {
		return future
	}

	// executing
	o.rw.Lock()
	defer o.rw.Unlock()
	future = o.future
	if future != nil {
		return future
	}

	// begin execute
	c := New[T]()
	o.future = c.future
	go o.do(c)
	return c.future
}
func (o *Once[T]) do(c *Completer[T]) {
	val, e := o.f()
	if e == nil {
		c.Complete(val)
	} else {
		c.CompleteError(e)
	}
}
