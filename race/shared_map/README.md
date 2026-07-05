# Shared Map Data Race

This example models a request metrics collector:

1. `Metrics.Record` updates a counter stored in a map.
2. `Metrics.Snapshot` copies those counters for reporting.
3. A writer goroutine records requests while another goroutine reads a snapshot.
4. The map owns shared mutable counters and they are not protected by a mutex, so reads and writes race.

The bug is in [main.go](main.go): `counts` owns shared mutable `Counter` values and both methods access those counters without synchronization.

## Run

```sh
make run
```

Expected result:

```text
checkout requests: 1
```

The standalone program is single-threaded, so it does not show the race.

## Catch With The Race Detector

```sh
make lint
```

Expected report shape:

```text
WARNING: DATA RACE
Write at 0xADDR by goroutine N:
  github.com/dvordrova/find_bugs/race/shared_map.(*Metrics).Record()
      github.com/dvordrova/find_bugs/race/shared_map/main.go:23 +0xADDR

Previous read at 0xADDR by goroutine N:
  github.com/dvordrova/find_bugs/race/shared_map.(*Metrics).Snapshot()
      github.com/dvordrova/find_bugs/race/shared_map/main.go:29 +0xADDR
```

Read the race report as two conflicting memory accesses:

1. `Write at` points to `Record`, which increments `Counter.Value`.
2. `Previous read at` points to `Snapshot`, which reads the same `Counter.Value`.
3. The goroutine creation stack shows where the concurrent writer was started in the test.

`make test` runs the same test without `-race`; it can pass because data races are not normal test assertions.

The committed [race.logs](race.logs) snapshot normalizes addresses, goroutine IDs, and package duration so diffs stay stable across machines.

## One Fix

Protect the map and its counters with a mutex:

```go
type Metrics struct {
	mu     sync.RWMutex
	counts map[string]int
}
```

Use `mu.Lock` in `Record` and `mu.RLock` in `Snapshot`. Another valid design is to send all metric updates through one owner goroutine.
