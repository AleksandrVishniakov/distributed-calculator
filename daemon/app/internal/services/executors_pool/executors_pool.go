package executors_pool

import "sync"

type Executor interface {
	Task()
}

type ExecutorsPool struct {
	tasks chan Executor
	wg    sync.WaitGroup
}

func NewExecutorsPool(maxGoroutines int) *ExecutorsPool {
	p := ExecutorsPool{
		tasks: make(chan Executor, 5*maxGoroutines),
	}

	p.wg.Add(maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		go func() {
			for w := range p.tasks {
				w.Task()
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
