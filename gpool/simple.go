package gpool

import (
	"context"
	"sync"
)

type SimplePool struct {
	maxWorkers int

	tasks     chan func()
	closeCh   chan struct{}
	closeOnce *sync.Once
}

func NewSimplePool(maxWorkers int) *SimplePool {
	if maxWorkers <= 0 {
		maxWorkers = 1
	}
	sp := &SimplePool{
		maxWorkers: maxWorkers,
	}
	sp.run()
	return sp
}

func (sp *SimplePool) run() {
	sp.tasks = make(chan func(), sp.maxWorkers)
	sp.closeCh = make(chan struct{})
	sp.closeOnce = &sync.Once{}
	for range sp.maxWorkers {
		go sp.worker()
	}
}

func (sp *SimplePool) worker() {
	for {
		select {
		case task := <-sp.tasks:
			if task != nil {
				task()
			}
		case <-sp.closeCh:
			return
		}
	}
}

func (sp *SimplePool) Submit(ctx context.Context, task func()) {
	if task == nil {
		return
	}
	select {
	case sp.tasks <- task:
	case <-ctx.Done():
	case <-sp.closeCh:
	}
}

func (sp *SimplePool) Close() {
	sp.closeOnce.Do(func() {
		close(sp.closeCh)
		close(sp.tasks)
	})
}
