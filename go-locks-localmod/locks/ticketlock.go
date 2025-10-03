package locks

import (
	"runtime"
	"sync/atomic"
)

// TicketLock is a FIFO spin lock.
type TicketLock struct {
	next  uint32
	owner uint32
	YieldWhileSpinning bool
}

func (l *TicketLock) Lock() {
	my := atomic.AddUint32(&l.next, 1) - 1
	for {
		if atomic.LoadUint32(&l.owner) == my { return }
		if l.YieldWhileSpinning { runtime.Gosched() }
	}
}

func (l *TicketLock) Unlock() { atomic.AddUint32(&l.owner, 1) }
