package pool

import (
	"sync"
	"sync/atomic"
	"time"
)

type fnType func()

type Pool struct {
	capacity int32
	applied  int32
	lock     sync.RWMutex
	cond     *sync.Cond
	idling   []*Worker
	expiry   time.Duration
	release  chan int32
	plumbing chan *Worker
}

func NewPool(cap int32, expTime time.Duration) *Pool {
	p := &Pool{
		capacity: cap,
		expiry:   expTime,
		release:  make(chan int32),
		plumbing: make(chan *Worker),
	}
	p.cond = sync.NewCond(&p.lock)
	return p
}

func (p *Pool) RegularCleanUp() {
	ticker := time.NewTicker(p.expiry)
	defer ticker.Stop()
	for {
		select {
		case <-p.release:
			return
		case <-ticker.C:
			p.lock.RLock()
			if p.applied == 0 || len(p.idling) == 0 {
				p.lock.RUnlock()
				continue
			}
			var idlingTooLong []*Worker
			now := time.Now()
			var last int //last is the last worker who idling over limit
			for i, w := range p.idling {
				if now.Sub(w.lastIdleTime) < p.expiry {
					break
				}
				idlingTooLong = append(idlingTooLong, w)
				last = i
			}
			p.lock.RUnlock()

			p.lock.Lock()
			if last == len(p.idling)-1 {
				p.idling = p.idling[:0]
			} else {
				// amount is the amount of worker who just idl.
				amount := copy(p.idling, p.idling[last+1:])
				for j := amount; j < len(p.idling); j++ {
					p.idling[j] = nil
				}
				p.idling = p.idling[:amount]
			}
			p.applied -= int32(last + 1)
			// Notify blocking condition
			p.cond.Broadcast()
			p.lock.Unlock()

			// Worker's task channel may be block when they are working
			// So move them to a individual array for quicker unlock mutex
			for _, w := range idlingTooLong {
				w.task <- nil
			}
		}
	}
}

func (p *Pool) RevertWorker() {
	for {
		select {
		case <-p.release:
			return
		case w := <-p.plumbing:
			if w == nil {
				return
			}
			p.lock.Lock()
			p.idling = append(p.idling, w)
			p.lock.Unlock()
		}
	}
}

func (p *Pool) Submit(fn fnType) {
	w := p.DetachWorker()
	w.task <- fn
}

func (p *Pool) Employ() *Worker {
	atomic.AddInt32(&p.applied, 1)
	return &Worker{plumbing: p.plumbing, task: make(chan fnType, 1)}
}

func (p *Pool) DetachWorker() *Worker {
	var w *Worker
	p.lock.RLock()
	if len(p.idling) == 0 {
		if p.applied > p.capacity {
			// Wait for auto clean-up
			p.cond.Wait()
		}
		p.lock.RUnlock()
		w = p.Employ()
		w.Do()
		return w
	}

	last := len(p.idling) - 1
	p.lock.RUnlock()
	p.lock.Lock()
	w = p.idling[last]
	p.idling[last] = nil
	p.idling = p.idling[:last]
	p.lock.Unlock()
	w.Do()
	return w
}
