// See the LICENSE file for license details.

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
	return newAutoResetEvent(false)
}

// NewAutoResetEventSignaled creates a new automatic reset event in the signaled state.
func NewAutoResetEventSignaled() *AutoResetEvent {
	return newAutoResetEvent(true)
}

// Reset resets the event
func (e *AutoResetEvent) Reset() {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	// Try to dequeue pending data in channel, if any
	select {
	case <-e.ch:
	default:
	}
}

// Set signals the event
func (e *AutoResetEvent) Set() {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	// Queue data if channel is empty.
	select {
	case e.ch <- struct{}{}:
	default:
	}
}

// TryWait tries to consume a signal without blocking.
func (e *AutoResetEvent) TryWait() bool {
	select {
	case <-e.ch:
		return true
	default:
		return false
	}
}

// WaitCh returns a channel that receives an empty data when set
func (e *AutoResetEvent) WaitCh() <-chan struct{} {
	return e.ch
}

// Wait waits until the event is signaled
func (e *AutoResetEvent) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-e.ch:
		return nil
	}
}

//------------------------------------------------------------------------------

func newAutoResetEvent(signaled bool) *AutoResetEvent {
	e := &AutoResetEvent{
		mtx: sync.Mutex{},
		ch:  make(chan struct{}, 1),
	}

	if signaled {
		e.ch <- struct{}{}
	}

	return e
}
