// Package gpool provides "advanced" version of sync.Pool.
// Supports generics (which is its biggest advantage) and 2 handlers that let you decide what to do if pointer is acquired or returned.
//
// There is also Discard function that gives you control when pointer should be rejected when calling Pool.Put, so it won't be used again. It's useful if you want to remove (for example) buffers that have grown over the limit.
//
// IMPORTANT: By default, Pool will set pointer to zero value of generic provided during initialization. It can be disabled by declaring OnInit handler. It's also important when using pooling for buffers!
package gpool

import "sync"

type Pool[V any] struct {
	pool sync.Pool
	opts *options[V]
}

func New[V any](opts ...Option[V]) *Pool[V] {
	_options := parseOpts(opts)
	return &Pool[V]{
		opts: _options,
		pool: sync.Pool{New: func() any {
			_val := new(V)
			if _options.onInit != nil {
				_options.onInit(_val)
			}
			return _val
		}},
	}
}

func (p *Pool[V]) Get() *V {
	v := p.pool.Get().(*V)
	if p.opts.onInit == nil {
		var x V
		*v = x
	}
	return v
}

func (p *Pool[V]) Put(x *V) {
	if p.opts.onInit == nil {
		var x2 V
		*x = x2
	}
	if p.opts.discard != nil {
		if p.opts.discard(x) {
			x = nil
			return
		}
	}
	if p.opts.onPut != nil {
		p.opts.onPut(x)
	}
	p.pool.Put(x)
}
