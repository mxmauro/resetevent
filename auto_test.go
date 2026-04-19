// See the LICENSE file for license details.

package resetevent_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mxmauro/resetevent"
)

//------------------------------------------------------------------------------

func TestAutoResetEventImplementsEvent(t *testing.T) {
	var _ resetevent.Event = (*resetevent.AutoResetEvent)(nil)
}

func TestAutoResetEventSetReleasesOneWaiter(t *testing.T) {
	e := resetevent.NewAutoResetEvent()
	released := make(chan struct{}, 2)

	for i := 0; i < 2; i++ {
		go func() {
			err := e.Wait(context.Background())
			if err == nil {
				released <- struct{}{}
			}
		}()
	}

	time.Sleep(20 * time.Millisecond)

	e.Set()

	select {
	case <-released:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("expected one waiter to be released")
	}

	select {
	case <-released:
		t.Fatal("expected only one waiter to be released")
	case <-time.After(50 * time.Millisecond):
	}
}

func TestAutoResetEventSetBeforeWaitIsRetainedOnce(t *testing.T) {
	e := resetevent.NewAutoResetEvent()
	e.Set()

	if err := e.Wait(context.Background()); err != nil {
		t.Fatalf("expected retained signal, got %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := e.Wait(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected timeout after retained signal is consumed, got %v", err)
	}
}

func TestAutoResetEventRepeatedSetDoesNotQueueMultipleSignals(t *testing.T) {
	e := resetevent.NewAutoResetEvent()
	e.Set()
	e.Set()

	if err := e.Wait(context.Background()); err != nil {
		t.Fatalf("expected first wait to succeed, got %v", err)
	}

	if e.TryWait() {
		t.Fatal("expected repeated Set calls to coalesce into a single signal")
	}
}

func TestAutoResetEventWaitCanceled(t *testing.T) {
	e := resetevent.NewAutoResetEvent()
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := e.Wait(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded, got %v", err)
	}
}

func TestAutoResetEventTryWait(t *testing.T) {
	e := resetevent.NewAutoResetEvent()

	if e.TryWait() {
		t.Fatal("expected TryWait to fail on an unsignaled event")
	}

	e.Set()

	if !e.TryWait() {
		t.Fatal("expected TryWait to consume a pending signal")
	}

	if e.TryWait() {
		t.Fatal("expected TryWait to consume only one signal")
	}
}

func TestAutoResetEventSignaledConstructor(t *testing.T) {
	e := resetevent.NewAutoResetEventSignaled()

	if !e.TryWait() {
		t.Fatal("expected signaled constructor to retain one signal")
	}
}
