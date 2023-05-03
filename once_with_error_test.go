package runonce

import (
	"fmt"
	"sync/atomic"
	"testing"
)

func TestRunWithError_success(t *testing.T) {
	var i atomic.Int64
	f := WrapWithError(func() (int64, error) {
		return i.Add(1), nil
	})

	for j := 0; j < 10; j++ {
		t.Run(fmt.Sprintf("Run %d", j), func(t *testing.T) {
			t.Parallel()
			v, err := f()
			if err != nil {
				t.Error(err)
			}
			if v != 1 {
				t.Errorf("Expected 1, got %d", v)
			}
		})
	}
}

func TestRunWithError_error(t *testing.T) {
	var i atomic.Int64
	f := WrapWithError(func() (int64, error) {
		return i.Add(1), fmt.Errorf("error")
	})

	for j := 0; j < 10; j++ {
		t.Run(fmt.Sprintf("Run %d", j), func(t *testing.T) {
			t.Parallel()
			v, err := f()
			if err == nil {
				t.Error("Expected error")
			}
			if v != 1 {
				t.Errorf("Expected 1, got %d", v)
			}
		})
	}
}
