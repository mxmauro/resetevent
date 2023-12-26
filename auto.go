package resetevent

import (
	"context"
	"sync"
)

//------------------------------------------------------------------------------

// AutoResetEvent represents an automatic reset event
type AutoResetEvent struct {
	mtx sync.Mutex
	ch  chan struct{}
}

//------------------------------------------------------------------------------

// NewAutoResetEvent creates a new automatic reset event
func NewAutoResetEvent() *AutoResetEvent {
	e := &AutoResetEvent{
		mtx: sync.Mutex{},
		ch:  make(chan struct{}, 1),
	}
	return e
}

// Reset resets the event
func (e *AutoResetEvent) Reset() {
	e.mtx.Lock()
	// Try to dequeue pending data in channel, if any
	select {
	case <-e.ch:
	default:
	}
	e.mtx.Unlock()
}

// Set signals the event
func (e *AutoResetEvent) Set() {
	e.mtx.Lock()
	// Queue data if channel is empty
	if len(e.ch) == 0 {
		e.ch <- struct{}{}
	}
	e.mtx.Unlock()
}

// WaitCh returns a channel that receives an empty data when set
func (e *AutoResetEvent) WaitCh() <-chan struct{} {
	return e.ch
}

// Wait waits until the event is signalled
func (e *AutoResetEvent) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-e.ch:
		return nil
	}
}
