// See the LICENSE file for license details.

// Package resetevent provides manual-reset and auto-reset synchronization events.
//
// Auto-reset events retain at most one signal and release at most one waiter per
// successful Set/TryWait cycle.
//
// Manual-reset events stay signaled until Reset is called, releasing all current
// waiters and allowing future waiters to proceed immediately while signaled.
package resetevent
