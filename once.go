package runonce

import "sync"

type F[T any] func() (T, error)

type Runner[T any] interface {
	Run() (T, error)
}

func New[T any](f func() (T, error)) Runner[T] {
	return &runner[T]{
		f:    f,
		done: make(chan struct{}),
	}
}

type runner[T any] struct {
	f    func() (T, error)
	done chan struct{}
	sync.Once

	t   T
	err error
	sync.RWMutex
}

func (r *runner[T]) Run() (T, error) {
	r.Once.Do(func() {
		r.Lock()
		defer r.Unlock()

		r.t, r.err = r.f()
		close(r.done)
	})

	<-r.done
	r.RLock()
	defer r.RUnlock()
	return r.t, r.err
}
