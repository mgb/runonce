package runonce

import "sync"

// Runner is a function that returns a value and an error.
type RunWithError[T any] func() (T, error)

// WrapWithError wraps a function that returns a value and an error, and
// ensures that the function is only called once and all callers see the same
// values
func WrapWithError[T any](f RunWithError[T]) RunWithError[T] {
	r := runner[T]{
		f: f,
	}
	return r.run
}

type runner[T any] struct {
	// f is:
	//  - called once (protected by sync.Once)
	//  - writes the results to t/err
	f   RunWithError[T]
	t   T
	err error
	sync.Once
}

func (r *runner[T]) run() (T, error) {
	r.Once.Do(func() {
		r.t, r.err = r.f()
	})

	// sync.Once ensures that all values have been successfully written before
	// returning, and all go routines will see the same values.
	return r.t, r.err
}
