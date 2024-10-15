// See the LICENSE file for license details.

package resetevent_test

import (
	"context"
	"testing"
	"time"

	"github.com/mxmauro/resetevent"
)

//------------------------------------------------------------------------------

func TestManualResetEvent(t *testing.T) {
	e := resetevent.NewManualResetEvent()

	go func() {
		<-time.After(1 * time.Second)
		e.Set()
	}()

	ch := make(chan struct{})

	for i := 1; i <= 5; i++ {
		go func() {
			_ = e.Wait(context.Background())
			ch <- struct{}{}
		}()
	}

	for i := 1; i <= 5; i++ {
		<-ch
	}
}
