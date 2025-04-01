package gpool

import (
	"context"
	"sync"
	"testing"
)

func TestSimplePool(t *testing.T) {
	sp := NewSimplePool(5)
	defer sp.Close()

	sw := &sync.WaitGroup{}

	for i := range 10 {
		sw.Add(1)
		// i := i // before go 1.22
		sp.Submit(context.Background(), func() {
			t.Log("task", i)
			sw.Done()
		})
	}

	sw.Wait()
}
