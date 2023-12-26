package resetevent

import (
	"context"
	"sync"
)

//------------------------------------------------------------------------------

// ManualResetEvent represents a manual reset event
type ManualResetEvent struct {
	mtx sync.Mutex
	ch  chan struct{}
}

//------------------------------------------------------------------------------

// NewManualResetEvent creates a new manual reset event
func NewManualResetEvent() *ManualResetEvent {
	e := &ManualResetEvent{
		mtx: sync.Mutex{},
		ch:  make(chan struct{}, 1),
	}
	return e
}

// Reset resets the event
func (e *ManualResetEvent) Reset() {
	e.mtx.Lock()
	// Recreate the channel if the is closed
	select {
	case <-e.ch:
		e.ch = make(chan struct{}, 1)
	default:
	}
	e.mtx.Unlock()
}

// Set signals the event
func (e *ManualResetEvent) Set() {
	e.mtx.Lock()
	// Close the channel so all waiting goroutines are awaken
	select {
	case <-e.ch:
	default:
		close(e.ch)
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
