// Copyright 2019 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pqueue // import "sevki.org/pqueue"

import (
	"container/heap"
	"sync"
)

// PQueue is a thread safe priorityqueue
type PQueue struct {
	q     pq
	c     *sync.Cond
	wg    sync.WaitGroup
	in    chan Item
	out   chan Item
	ready chan interface{}
}

// Item is the interface that is necessary for
// priorityqueue to figure out priorities
type Item interface {
	Priority() float64
}

// New returns a new queue
func New() *PQueue {
	q := &PQueue{
		c:     sync.NewCond(&sync.Mutex{}),
		in:    make(chan Item),
		ready: make(chan interface{}),
		out:   make(chan Item),
	}
	heap.Init(&q.q)

	go func() {
		for {
			select {
			case i := <-q.in:
				q.c.L.Lock()
				heap.Push(&q.q, i)
				q.c.Signal()
				q.c.L.Unlock()
			}
		}
	}()
	go func() {
		for {
			select {
			case <-q.ready:
			}
			q.c.L.Lock()
			if q.q.Len() == 0 {
				q.c.Wait()
			}
			x := heap.Pop(&q.q)
			q.c.L.Unlock()
			q.out <- x.(Item)
		}
	}()
	return q
}

// Push adds a new item to the queue
func (q *PQueue) Push(i Item) {
	q.wg.Add(1)
	q.in <- i
}

// Pop returns an element from the queue
// blocks until it can return an element
func (q *PQueue) Pop() Item {
	defer q.wg.Done()
	q.ready <- nil
	return <-q.out
}

// Wait blocks until no more Items are in the queue
func (q *PQueue) Wait() { q.wg.Wait() }

type pq []Item

func (pq pq) Len() int { return len(pq) }

func (pq pq) Less(i, j int) bool {
	return pq[i].Priority() > pq[j].Priority()
}

func (pq pq) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *pq) Push(x interface{}) {
	if item, ok := x.(Item); ok {
		*pq = append(*pq, item)
	}
}

func (pq *pq) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
