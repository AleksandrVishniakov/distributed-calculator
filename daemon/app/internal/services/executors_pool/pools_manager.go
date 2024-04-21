package executors_pool

import (
	"context"
	"sync"
)

type PoolManager struct {
	mu    *sync.RWMutex
	pools map[uint64]*ExecutorsPool

	maxGoroutines int
}

func NewManager(maxGoroutines int) *PoolManager {
	return &PoolManager{
		mu:            &sync.RWMutex{},
		pools:         make(map[uint64]*ExecutorsPool),
		maxGoroutines: maxGoroutines,
	}
}

func (p *PoolManager) Pool(userID uint64) *ExecutorsPool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if pool, ok := p.pools[userID]; ok {
		return pool
	}

	pool := NewExecutorsPool(context.Background(), p.maxGoroutines)

	p.pools[userID] = pool

	return pool
}

func (p *PoolManager) Shutdown() {
	wg := &sync.WaitGroup{}

	p.mu.Lock()
	defer p.mu.Unlock()
	for _, pool := range p.pools {
		wg.Add(1)
		go func(pool *ExecutorsPool) {
			defer wg.Done()

			pool.Shutdown()
		}(pool)
	}

	wg.Wait()
}
