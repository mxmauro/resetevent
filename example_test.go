// See the LICENSE file for license details.

package resetevent_test

import (
	"context"
	"fmt"

	"github.com/mxmauro/resetevent"
)

func ExampleAutoResetEvent() {
	e := resetevent.NewAutoResetEventSignaled()

	fmt.Println(e.TryWait())
	fmt.Println(e.TryWait())

	// Output:
	// true
	// false
}

func ExampleManualResetEvent() {
	e := resetevent.NewManualResetEvent()
	e.Set()

	_ = e.Wait(context.Background())
	fmt.Println("released")

	e.Reset()

	// Output:
	// released
}
