package resetevent

import (
	"context"
)

//------------------------------------------------------------------------------

type Event interface {
	Reset()
	Set()

	WaitCh() <-chan struct{}
	Wait(ctx context.Context)
}
