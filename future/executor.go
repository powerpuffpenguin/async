package future

import (
	"sync"
)

type OnceExecutor interface {
	Execute() (value any, e error)
}
type OnceExecute[T any] struct {
	value T
	err   error
	f     func() (T, error)
}

// new an executor by func
func NewOnceExecute[T any](f func() (T, error)) *OnceExecute[T] {
	return &OnceExecute[T]{
		f: f,
	}
}

// Implement the OnceExecutor interface
func (o *OnceExecute[T]) Execute() (any, error) {
	o.value, o.err = o.f()
	return o.value, o.err
}

// For the same concurrent executor, the same request result will be shared
type Executor struct {
	keys map[OnceExecutor]*Future[any]
	rw   sync.RWMutex
}

func NewExecutor() *Executor {
	return &Executor{
		keys: make(map[OnceExecutor]*Future[any]),
	}
}

// run an executor
func (m *Executor) Do(executor OnceExecutor) *Future[any] {
	// executing
	m.rw.RLock()
	future := m.keys[executor]
	m.rw.RUnlock()
	if future != nil {
		return future
	}
	// executing
	m.rw.Lock()
	defer m.rw.Unlock()
	future = m.keys[executor]
	if future != nil {
		return future
	}
	// begin execute
	c := New[any]()
	m.keys[executor] = c.future
	go m.do(executor, c)
	return c.future
}
func (m *Executor) do(executor OnceExecutor, c *Completer[any]) {
	value, e := executor.Execute()
	if e == nil {
		c.Complete(value)
	} else {
		c.CompleteError(e)
	}
	// delete cache
	m.rw.Lock()
	delete(m.keys, executor)
	m.rw.Unlock()
}
