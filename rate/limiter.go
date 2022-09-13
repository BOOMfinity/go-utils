package rate

import (
	"context"
	"sync"
	"time"
)

type Limiter struct {
	last       time.Time
	new        chan chan<- bool
	resetAfter time.Duration
	limit      int
	sent       int
	init       int
	m          sync.Mutex
}

func (v *Limiter) Update(resetAfter time.Duration, limit int) {
	v.m.Lock()
	v.resetAfter = resetAfter
	v.limit = limit
	v.m.Unlock()
}

func (v *Limiter) loop() {
	for ch := range v.new {
		(func() {
			v.m.Lock()
			defer v.m.Unlock()
		backHere:
			if time.Now().After(v.last.Add(v.resetAfter)) {
				v.sent = v.init
				if v.init != 0 {
					v.last = time.Now()
				}
				v.init = 0
			}
			if (v.limit - v.sent) == 0 {
				time.Sleep(v.resetAfter + (50 * time.Millisecond))
				v.sent = 0
				goto backHere
			}
			v.sent++
			if v.sent == 1 {
				v.last = time.Now()
			}
			ch <- true
		})()
	}
}

func (v *Limiter) Wait(ctx context.Context) error {
	ch := make(chan bool)
	v.new <- ch
	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func NewLimiter(resetAfter time.Duration, limit int) *Limiter {
	l := new(Limiter)
	l.resetAfter = resetAfter
	l.limit = limit
	l.new = make(chan chan<- bool, 1)
	go l.loop()
	return l
}

func NewLimiterInit(resetAfter time.Duration, limit int, init int) *Limiter {
	l := NewLimiter(resetAfter, limit)
	l.init = init
	return l
}
