# resetevent

Manual and automatic reset synchronization event objects.

## Behavior

When an event object becomes signaled (`Set`), waiting goroutines are selected and can continue working.

For auto-reset events, only one goroutine is selected and the object becomes non-signaled again.
For manual-reset events, any waiting goroutine is selected until the `Reset` method is called.

## API semantics

### `AutoResetEvent`

- `Set` retains at most one pending signal.
- If `Set` is called before `Wait`, the next waiter proceeds immediately.
- Repeated `Set` calls while already signaled are coalesced into a single pending signal.
- `TryWait` consumes a pending signal without blocking and reports whether it succeeded.
- `WaitCh` exposes the internal signal channel and receives one value per successful release.

### `ManualResetEvent`

- `Set` releases all current waiters.
- Once signaled, future calls to `Wait` return immediately until `Reset` is called.
- `Reset` creates a new wait channel. Code that caches a prior `WaitCh()` must refresh it after reset.
- `WaitCh` returns a channel that is closed while the event is signaled.

### Cancellation

- `Wait` returns `context.Canceled` or `context.DeadlineExceeded` when the provided context ends first.

## Usage

```go
ctx := context.Background()

auto := resetevent.NewAutoResetEvent()
auto.Set()
_ = auto.Wait(ctx) // consumes one signal

manual := resetevent.NewManualResetEventSignaled()
_ = manual.Wait(ctx) // returns immediately while signaled
manual.Reset()
```

## LICENSE

See the [license](LICENSE) file for details.

Portions of this code is based on or derived from the [original work](https://github.com/xcdb/syncx) by Chris Burge. 
