package runonce

import "sync"

// New creates a Runner that runs the given function once.
func New[T any](f func() (T, error)) func() (T, error) {
	r := runner[T]{
		f: f,
	}
	return r.Run
}

type runner[T any] struct {
	// f is:
	//  - called once (protected by sync.Once)
	//  - writes the results to t/err
	f   func() (T, error)
	t   T
	err error
	sync.Once
}

func (r *runner[T]) Run() (T, error) {
	r.Once.Do(func() {
		r.t, r.err = r.f()
	})

	// sync.Once ensures that all values have been successfully written before
	// returning, and all go routines will see the same values.
	return r.t, r.err
}
