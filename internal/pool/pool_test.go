package pool

import (
	"testing"
	"time"
)

func TestPool_RegularCleanUp(t *testing.T) {
	p := NewPool(10, time.Second)
	p.applied = 10
	tCh := make(chan fnType, 1)
	for i := 0; i < 10; i++ {
		if i < 5 {
			p.idling = append(p.idling, &Worker{
				task:         tCh,
				lastIdleTime: time.Now().Add(-time.Second),
			})
		} else {
			p.idling = append(p.idling, &Worker{
				task:         tCh,
				lastIdleTime: time.Now().Add(time.Second),
			})
		}
	}
	go p.RegularCleanUp()
	if done := <-tCh; done == nil {
		p.release <- 0
		if len(p.idling) > 5 {
			t.Fatal("Unwanted result")
		}
	}
}

func TestPool_Employ(t *testing.T) {
	p := NewPool(10, time.Second)
	w := p.Employ()
	if p.applied != 1 {
		t.Fatal("Apply fail")
	}
	if w == nil {
		t.Fatal("Get Worker fail")
	}
	if w.plumbing == nil {
		t.Fatal("Worker's pipe is not initialize")
	}
	if w.task == nil {
		t.Fatal("Worker's task channel is not initialize")
	}
}

func TestPool_RevertWorker(t *testing.T) {
	p := NewPool(10, time.Second)
	go p.RevertWorker()
	p.plumbing <- &Worker{}
	p.release <- 0
	if len(p.idling) == 0 {
		t.Fatal("Revert worker failed")
	}

}

func TestPool_DetachWorker(t *testing.T) {
	var i int
	testFn := func() {
		i++
	}
	p := NewPool(10, time.Second)
	p.Submit(testFn)
	<-p.plumbing
	if i != 1 {
		t.Fatal("Do job failed")
	}
	// reset
	i = 0

	go p.RevertWorker()
	for j := 0; j < 100; j++ {
		p.Submit(testFn)
	}

	if i < 99 {
		t.Fatal("Do jobs failed")
	}

}
