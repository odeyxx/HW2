# Lock Implementations in Go

This repository contains two lock algorithms and a small benchmark program.

## Files
- `locks/ticketlock.go` — ticket (FIFO) spin lock
- `locks/casspin.go` — compare-and-swap spin lock
- `cmd/bench/main.go` — benchmark driver

## Build and Run
```bash
go run ./cmd/bench -lock=ticket -threads=8 -iters=200000 -critus=5
go run ./cmd/bench -lock=cas    -threads=8 -iters=200000 -critus=5
```

### Flags
- `-lock` : `ticket` or `cas` (default `ticket`)
- `-threads` : number of goroutines
- `-iters` : lock acquisitions per goroutine
- `-critus` : microseconds of work inside the critical section
- `-outsideus` : microseconds of work outside the critical section
- `-yield` : call `runtime.Gosched()` while spinning

### Notes
Implementations spin by design to match the algorithms discussed in class.
