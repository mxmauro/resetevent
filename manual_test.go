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

func TestManualResetEventImplementsEvent(t *testing.T) {
	var _ resetevent.Event = (*resetevent.ManualResetEvent)(nil)
}

func TestManualResetEventBroadcastsAndAllowsFutureWaiters(t *testing.T) {
	e := resetevent.NewManualResetEvent()
	released := make(chan struct{}, 6)

	for i := 0; i < 5; i++ {
		go func() {
			err := e.Wait(context.Background())
			if err == nil {
				released <- struct{}{}
			}
		}()
	}

	time.Sleep(20 * time.Millisecond)

	e.Set()

	for i := 0; i < 5; i++ {
		select {
		case <-released:
		case <-time.After(200 * time.Millisecond):
			t.Fatal("expected all waiters to be released")
		}
	}

	if err := e.Wait(context.Background()); err != nil {
		t.Fatalf("expected future waiter to pass immediately while signaled, got %v", err)
	}
}

func TestManualResetEventResetBlocksAgainAndChangesWaitCh(t *testing.T) {
	e := resetevent.NewManualResetEvent()
	ch1 := e.WaitCh()

	e.Set()
	e.Reset()

	ch2 := e.WaitCh()
	if ch1 == ch2 {
		t.Fatal("expected Reset to replace the wait channel after a signaled state")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := e.Wait(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected reset event to block again, got %v", err)
	}
}

func TestManualResetEventWaitCanceled(t *testing.T) {
	e := resetevent.NewManualResetEvent()
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := e.Wait(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded, got %v", err)
	}
}

func TestManualResetEventSignaledConstructor(t *testing.T) {
	e := resetevent.NewManualResetEventSignaled()

	if err := e.Wait(context.Background()); err != nil {
		t.Fatalf("expected signaled constructor to allow immediate wait, got %v", err)
	}
}
