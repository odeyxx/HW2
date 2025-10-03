package locks

import (
	"runtime"
	"sync/atomic"
)

// CASLock is a basic compare-and-swap spin lock.
type CASLock struct {
	flag int32
	YieldWhileSpinning bool
}

func (l *CASLock) Lock() {
	for {
		if atomic.CompareAndSwapInt32(&l.flag, 0, 1) { return }
		if l.YieldWhileSpinning { runtime.Gosched() }
	}
}

func (l *CASLock) Unlock() { atomic.StoreInt32(&l.flag, 0) }
