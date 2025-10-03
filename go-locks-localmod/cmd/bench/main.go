package main

import (
	"flag"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"go-locks/locks"
)

func busyFor(us int) {
	if us <= 0 { return }
	target := time.Now().Add(time.Duration(us) * time.Microsecond)
	for time.Now().Before(target) {}
}

type simpleLock interface{ Lock(); Unlock() }

func makeLock(kind string, yield bool) simpleLock {
	switch kind {
	case "ticket":
		return &locks.TicketLock{YieldWhileSpinning: yield}
	case "cas":
		return &locks.CASLock{YieldWhileSpinning: yield}
	default:
		panic("unknown -lock: " + kind)
	}
}

func main() {
	lockKind := flag.String("lock", "ticket", "ticket|cas")
	threads := flag.Int("threads", 4, "goroutines contending")
	iters := flag.Int("iters", 100000, "acquisitions per goroutine")
	critUS := flag.Int("critus", 0, "µs work in critical section")
	outsideUS := flag.Int("outsideus", 0, "µs work outside")
	yield := flag.Bool("yield", false, "call Gosched() while spinning")
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())
	l := makeLock(*lockKind, *yield)

	fmt.Printf("lock=%s threads=%d iters/g=%d crit=%dµs outside=%dµs yield=%v\n",
		*lockKind, *threads, *iters, *critUS, *outsideUS, *yield)

	var counter int64
	var maxWait, totalWait int64

	start := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(*threads)
	for g := 0; g < *threads; g++ {
		go func() {
			defer wg.Done()
			var localWait, localMax int64
			for i := 0; i < *iters; i++ {
				t0 := time.Now()
				l.Lock()
				w := time.Since(t0).Nanoseconds()
				localWait += w
				if w > localMax { localMax = w }
				busyFor(*critUS)
				atomic.AddInt64(&counter, 1)
				l.Unlock()
				busyFor(*outsideUS)
			}
			atomic.AddInt64(&totalWait, localWait)
			for {
				old := atomic.LoadInt64(&maxWait)
				if localMax <= old { break }
				if atomic.CompareAndSwapInt64(&maxWait, old, localMax) { break }
			}
		}()
	}
	wg.Wait()

	elapsed := time.Since(start)
	ops := int64(*threads * *iters)
	avgWait := time.Duration(totalWait / ops)

	fmt.Printf("ops=%d totalWait=%v avgWait=%v maxWait=%v\n",
		ops, time.Duration(totalWait), avgWait, time.Duration(maxWait))
	fmt.Printf("elapsed=%v throughput=%.1f ops/s\n", elapsed, float64(ops)/elapsed.Seconds())
	fmt.Printf("counter=%d (should equal ops)\n", counter)
}
