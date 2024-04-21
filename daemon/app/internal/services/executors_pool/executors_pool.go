package executors_pool

import (
	"context"
	"sync"
)

type Executor interface {
	Task(ctx context.Context)
}

type ExecutorsPool struct {
	tasks chan Executor
	wg    sync.WaitGroup
}

func NewExecutorsPool(ctx context.Context, maxGoroutines int) *ExecutorsPool {
	p := ExecutorsPool{
		tasks: make(chan Executor, 5*maxGoroutines),
	}

	p.wg.Add(maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		go func() {
			for w := range p.tasks {
				w.Task(ctx)
			}

			p.wg.Done()
		}()
	}

	return &p
}

func (p *ExecutorsPool) Run(e Executor) {
	p.tasks <- e
}

func (p *ExecutorsPool) Shutdown() {
	close(p.tasks)
	p.wg.Wait()
}
