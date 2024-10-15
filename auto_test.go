// See the LICENSE file for license details.

package resetevent_test

import (
	"context"
	"testing"
	"time"

	"github.com/mxmauro/resetevent"
)

//------------------------------------------------------------------------------

func TestAutoResetEvent(t *testing.T) {
	e := resetevent.NewAutoResetEvent()

	go func() {
		<-time.After(1 * time.Second)
		e.Set()
	}()

	ch := make(chan struct{})

	go func() {
		_ = e.Wait(context.Background())
		ch <- struct{}{}
	}()

	<-ch
}
