// See the LICENSE file for license details.

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
	return newManualResetEvent(false)
}

// NewManualResetEventSignaled creates a new manual reset event in the signaled state.
func NewManualResetEventSignaled() *ManualResetEvent {
	return newManualResetEvent(true)
}

// Reset resets the event
func (e *ManualResetEvent) Reset() {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	if e.closed {
		e.ch = make(chan struct{})
		e.closed = false
	}
}

// Set signals the event
func (e *ManualResetEvent) Set() {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	// Close the channel so all waiting goroutines are awaken
	if !e.closed {
		close(e.ch)
		e.closed = true
	}
}

// WaitCh returns a channel that is closed when set
func (e *ManualResetEvent) WaitCh() <-chan struct{} {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	ch := e.ch
	return ch
}

// Wait waits until the event is signaled
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

//------------------------------------------------------------------------------

func newManualResetEvent(signaled bool) *ManualResetEvent {
	e := &ManualResetEvent{
		mtx: sync.Mutex{},
		ch:  make(chan struct{}),
	}

	if signaled {
		close(e.ch)
		e.closed = true
	}

	return e
}
