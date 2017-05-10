// Copyright 2017 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pqueue

import (
	"container/heap"
	"sync"
)

// PQueue is a thread safe priorityqueue
type PQueue struct {
	q pq
	c *sync.Cond
	l *sync.RWMutex
}

// Item is the interface that is necessary for
// priorityqueue to figure out priorities
type Item interface {
	Priority() int
}

// New returns a new queue
func New() *PQueue {
	lock := &sync.RWMutex{}
	ch := &PQueue{
		c: sync.NewCond(lock.RLocker()),
		l: lock,
	}
	heap.Init(&ch.q)
	return ch
}

// Push adds a new item to the queue
func (q *PQueue) Push(i interface{}) {
	q.l.Lock()
	heap.Push(&q.q, i)
	q.l.Unlock()
	q.c.Signal()
}

// Pop returns an element from the queue
// blocks until it can return an element
func (q *PQueue) Pop() interface{} {
	q.c.L.Lock()
	for q.q.Len() == 0 {
		q.c.Wait()
	}
	x := heap.Pop(&q.q)
	q.c.L.Unlock()
	return x
}

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
