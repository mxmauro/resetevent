# resetevent

Manual and automatic reset synchronization event objects.

## Behavior

When an event object becomes signaled (`Set`), waiting Goroutines are selected and can continue working.

For auto-reset events, only one goroutine is selected and the object becomes non-signaled again.
For manual-reset events, any waiting goroutine is selected until the `Reset` method is called.

## LICENSE

See the [license](LICENSE) file for details.

Portions of this code is based on or derived from the [original work](https://github.com/xcdb/syncx) by Chris Burge. 
