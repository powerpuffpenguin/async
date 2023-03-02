package future

var completed = make(chan struct{})

func init() {
	close(completed)
}

// Create a Future that returns a fixed value
func Value[T any](value T) *Future[T] {
	return &Future[T]{
		value: value,
		c:     completed,
	}
}

// Create an f that returns a fixed error
func Error[T any](err error) *Future[T] {
	return &Future[T]{
		err: err,
		c:   completed,
	}
}
