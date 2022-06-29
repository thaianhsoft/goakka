package goakka

import (
	"fmt"
	"sync"
)

type Pool struct {
	sizeThread int
	queue *Queue
	cv_threads *sync.Cond
	locker *sync.Mutex
	threadProcessFunc ThreadProcessFunc
}

func NewPool(sizeThread int, tf ThreadProcessFunc) *Pool {
	p := &Pool{
		queue:      NewQueue(),
		cv_threads: sync.NewCond(&sync.Mutex{}),
		locker:     &sync.Mutex{},
		sizeThread: sizeThread,
		threadProcessFunc: tf,
	}
	return p
}

func (p *Pool) PushBackJob(job any) {
	p.locker.Lock()
	p.queue.PushBack(job)
	p.locker.Unlock()
	p.cv_threads.Signal()
}

func (p *Pool) PushFrontJob(job any) {
	p.locker.Lock()
	p.queue.PushFront(job)
	p.locker.Unlock()
	p.cv_threads.Signal()
}
func (p *Pool) spawnThread() {
	for id := 0; id < p.sizeThread; id++ {
		fmt.Printf("spawn thread id=%v\n", id)
		go p.runThread(uint32(id))
	}
}

func (p *Pool) runThread(threadId uint32) {
	var initBehaviour ThreadBehaviour = &RunnableBehaviour{}
	for {
		initBehaviour = initBehaviour.Process(p.threadProcessFunc)
	}
}

func (p *Pool) GetQueue() *Queue {
	return p.queue
}

func (p *Pool) GetCondVariableLocker() *sync.Cond {
	return p.cv_threads
}
