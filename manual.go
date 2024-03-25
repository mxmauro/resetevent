package resetevent

import (
	"context"
	"sync"
)

//------------------------------------------------------------------------------

// ManualResetEvent represents a manual reset event
type ManualResetEvent struct {
	mtx    sync.Mutex
	ch     chan struct{}
	closed bool
}

//------------------------------------------------------------------------------

// NewManualResetEvent creates a new manual reset event
func NewManualResetEvent() *ManualResetEvent {
	e := &ManualResetEvent{
		mtx: sync.Mutex{},
		ch:  make(chan struct{}),
	}
	return e
}

// Reset resets the event
func (e *ManualResetEvent) Reset() {
	e.mtx.Lock()
	if e.closed {
		e.ch = make(chan struct{})
		e.closed = false
	}
	e.mtx.Unlock()
}

// Set signals the event
func (e *ManualResetEvent) Set() {
	e.mtx.Lock()
	// Close the channel so all waiting goroutines are awaken
	if !e.closed {
		close(e.ch)
		e.closed = true
	}
	e.mtx.Unlock()
}

// WaitCh returns a channel that is closed when set
func (e *ManualResetEvent) WaitCh() <-chan struct{} {
	e.mtx.Lock()
	ch := e.ch
	e.mtx.Unlock()
	return ch
}

// Wait waits until the event is signalled
func (e *ManualResetEvent) Wait(ctx context.Context) error {
	e.mtx.Lock()
	ch := e.ch
	e.mtx.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ch:
		return nil
	}
}
