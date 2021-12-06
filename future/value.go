package future

var completed = make(chan struct{})

func init() {
	close(completed)
}

// Create a Future that returns a fixed value
func Value(value interface{}) *Future {
	return &Future{
		value: value,
		c:     completed,
	}
}

// Create an f that returns a fixed error
func Error(err error) *Future {
	return &Future{
		err: err,
		c:   completed,
	}
}
