package future

// Wrapping a future value is usually used to store asynchronous results
//
// like dart Future<T>
type Future[T any] struct {
	value T
	err   error
	c     chan struct{}
}

// Return a chan and wait until notification when the Future is completed
func (f *Future[T]) C() <-chan struct{} {
	return f.c
}

// Returns the value stored by the Future, and does not guarantee that the value is valid
//
// To ensure that a valid value is obtained, you can call Get or check first to determine the valid value stored in the Future
// For example:
//
// <-f.C()
// f.Value()
func (f *Future[T]) Value() (T, error) {
	return f.value, f.err
}

// Return the set value after the Future is set to a valid value
func (f *Future[T]) Get() (T, error) {
	<-f.c
	return f.value, f.err
}
