//go:build !solution

package cond

// A Locker represents an object that can be locked and unlocked.
type Locker interface {
	Lock()
	Unlock()
}

// Cond implements a condition variable, a rendezvous point
// for goroutines waiting for or announcing the occurrence
// of an event.
//
// Each Cond has an associated Locker L (often a *sync.Mutex or *sync.RWMutex),
// which must be held when changing the condition and
// when calling the Wait method.
type Cond struct {
	L Locker
    is_goroutine chan chan struct{}
}

// New returns a new Cond with Locker l.
func New(l Locker) *Cond {
	return &Cond{L: l, is_goroutine: make(chan chan struct{}, 500)}
}

// Wait atomically unlocks c.L and suspends execution
// of the calling goroutine. After later resuming execution,
// Wait locks c.L before returning. Unlike in other systems,
// Wait cannot return unless awoken by Broadcast or Signal.
//
// Because c.L is not locked when Wait first resumes, the caller
// typically cannot assume that the condition is true when
// Wait returns. Instead, the caller should Wait in a loop:
//
//    c.L.Lock()
//    for !condition() {
//        c.Wait()
//    }
//    ... make use of condition ...
//    c.L.Unlock()
//
func (c *Cond) Wait() {
	if c.L == nil {
		panic("L is nil")
	}
	new_ch := make(chan struct{}, 1)
	defer func() {
		c.L.Lock() // блокирует текущую горутину
	}()	
	c.is_goroutine <- new_ch
	c.L.Unlock() // снимает блокировку `L`
	<-new_ch
}

// Signal wakes one goroutine waiting on c, if there is any.
//
// It is allowed but not required for the caller to hold c.L
// during the call.
func (c *Cond) helpSignal() bool {
	select {
	case ch := <- c.is_goroutine:
		select {
		case ch <- struct{}{}:
		default:
		}
		return true
	default:
	}
	return false

}
func (c *Cond) Signal() {
	c.helpSignal()

}

// Broadcast wakes all goroutines waiting on c.
//
// It is allowed but not required for the caller to hold c.L
// during the call.
func (c *Cond) Broadcast() {
	for c.helpSignal() {

	}
}
