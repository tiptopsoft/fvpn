package main

import (
	"sync"
)

type MemoryPool struct {
	lock sync.Mutex
	pool sync.Pool
	//cond sync.Cond
}

func NewPool(new func() any) *MemoryPool {
	return &MemoryPool{
		pool: sync.Pool{New: new},
	}
}

func (p *MemoryPool) Get() any {
	return p.pool.Get()
}

func (p *MemoryPool) Put(x any) {
	//defer p.lock.Unlock()
	p.pool.Put(x)
}
