// Copyright 2017 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pqueue

import (
	"context"
	"sync"
	"testing"
	"time"
)

type pint int

func (p pint) Priority() int {
	return int(p)
}

func TestPQueue(t *testing.T) {
	in := []pint{5, 3, 6, 2, 1, 4}
	out := []int{6, 5, 4, 3, 2, 1}
	q := New()

	for _, i := range in {
		q.Push(pint(i))
	}
	for _, outed := range out {
		poped := q.Pop()
		if poped != pint(outed) {
			t.Errorf("%d is not the same as %d", poped, outed)
			t.Fail()
		}
	}
}

func TestEmptyPQueue(t *testing.T) {
	ctx := context.Background()
	tctx, cancel := context.WithTimeout(ctx, time.Second)
	q := New()

	x := func() <-chan interface{} {
		ch := make(chan interface{})
		go func() {
			poped := q.Pop()
			_ = poped
			cancel()
			ch <- nil
		}()
		return ch
	}
	select {
	case <-x():
		t.Fail()
	case <-tctx.Done():

	}

}

func TestPQueuePushAfterPop(t *testing.T) {
	ctx := context.Background()
	tctx, cancel := context.WithTimeout(ctx, time.Second)
	q := New()

	x := func() <-chan interface{} {
		ch := make(chan interface{})
		go func() {
			poped := q.Pop()
			_ = poped
			cancel()
			ch <- nil
		}()
		return ch
	}
	getChan := x()
	q.Push(pint(1))
	select {
	case <-getChan:
	case <-tctx.Done():
		t.Fail()
	}

}

func TestPQueueConcurrent(t *testing.T) {
	in := []pint{5, 3, 6, 2, 1, 4}
	out := []int{6, 5, 4, 3, 2, 1}

	q := New()
	var wg sync.WaitGroup

	for _, i := range in {
		wg.Add(1)
		go func(i pint) {
			q.Push(pint(i))
			wg.Done()
		}(i)
	}
	wg.Wait()

	for _, outed := range out {
		poped := q.Pop()
		if poped != pint(outed) {
			t.Errorf("%d is not the same as %d", poped, outed)
			t.Fail()
		}
	}
}
