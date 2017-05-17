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
	tctx, _ := context.WithTimeout(ctx, time.Millisecond)
	q := New()

	x := func() <-chan interface{} {
		ch := make(chan interface{})
		go func() {
			poped := q.Pop()
			_ = poped
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
	tctx, _ := context.WithTimeout(ctx, time.Second)
	q := New()

	ch := make(chan interface{})

	func() {
		go func() {
			poped := q.Pop()
			_ = poped
			ch <- nil
		}()
	}()
	q.Push(pint(1))

	select {
	case <-ch:
	case <-tctx.Done():
		t.Log(tctx.Err())
		t.Fail()
	}
}

func TestPQueueConcurrent(t *testing.T) {
	in := []pint{2, 123, 12, 1, 11, 90}
	out := []int{123, 90, 12, 11, 2, 1}

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
	time.Sleep(500 * time.Millisecond)
	for _, outed := range out {
		poped := q.Pop()
		if poped != pint(outed) {
			t.Errorf("%d is not the same as %d", poped, outed)
			t.Fail()
		}
	}
}
