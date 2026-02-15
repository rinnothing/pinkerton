package healthcheck

import (
	"context"
	"sync"
	"time"

	"github.com/rinnothing/pinkerton/pkg/heap"
)

type timeAndPeriod struct {
	time   time.Time
	period time.Duration
}

type emitter struct {
	mx        sync.Mutex
	closest   *time.Time
	container heap.Heap[timeAndPeriod, string]
	update    chan struct{}
}

func newEmitter() *emitter {
	return &emitter{
		container: heap.New[timeAndPeriod, string](func(a, b timeAndPeriod) bool {
			return a.time.Before(b.time)
		}),
		closest: nil,
		update:  make(chan struct{}),
	}
}

func (e *emitter) updateClosest() {
	tm, _, _ := e.container.Top()
	e.closest = &tm.time
	e.update <- struct{}{}
}

func (e *emitter) AddEvent(at timeAndPeriod, url string) {
	e.mx.Lock()
	defer e.mx.Unlock()

	e.container.Push(url, at)
	e.updateClosest()
}

func (e *emitter) UpdateEvent(at timeAndPeriod, url string) {
	e.mx.Lock()
	defer e.mx.Unlock()

	e.container.Remove(url)
	e.container.Push(url, at)
	e.updateClosest()
}

func (e *emitter) RemoveEvent(url string) {
	e.mx.Lock()
	defer e.mx.Unlock()

	e.container.Remove(url)
	e.updateClosest()
}

// Start should only be called once (need to fix this with once)
func (e *emitter) Start(ctx context.Context) <-chan string {
	events := make(chan string)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-e.update:
			}

			for {
				select {
				case <-ctx.Done():
					return
				case <-e.update:
					continue
				case <-time.After(time.Until(*e.closest)):
				}

				empty := e.afterTimeEvent(ctx, events)
				if empty {
					break
				}
			}
		}
	}()

	return events
}

func (e *emitter) afterTimeEvent(ctx context.Context, events chan string) (empty bool) {
	e.mx.Lock()
	defer e.mx.Unlock()

	fstTm, fstId, ok := e.container.Top()
	if !ok {
		empty = true
		return true
	}

	for fstTm.time.Before(time.Now()) {
		e.container.Pop()
		e.container.Push(fstId, timeAndPeriod{time: time.Now().Add(fstTm.period), period: fstTm.period})

		select {
		case <-ctx.Done():
			return
		case events <- fstId:
		}

		fstTm, fstId, ok = e.container.Top()
		if !ok {
			empty = true
			break
		}
	}

	e.closest = &fstTm.time
	return
}
