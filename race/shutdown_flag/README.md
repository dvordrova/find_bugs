# Shutdown Flag Data Race

This example models a worker with a stop flag:

1. `Worker.Run` checks `stopping` in a loop.
2. `Worker.Stop` writes `stopping = true` from another goroutine.
3. There is no mutex, atomic value, or channel close connecting the read and write.
4. The worker usually stops, but the flag access is still a data race.

The bug is in [main.go](main.go): `stopping` is shared between goroutines without synchronization.

## Run

```sh
make run
```

Expected result:

```text
worker stopped
```

The program often appears to work, which is why this bug is easy to miss.

## Catch With The Race Detector

```sh
make lint
```

Expected report shape:

```text
WARNING: DATA RACE
Write at 0xADDR by goroutine N:
  github.com/dvordrova/find_bugs/race/shutdown_flag.(*Worker).Stop()
      github.com/dvordrova/find_bugs/race/shutdown_flag/main.go:26 +0xADDR

Previous read at 0xADDR by goroutine N:
  github.com/dvordrova/find_bugs/race/shutdown_flag.(*Worker).Run()
      github.com/dvordrova/find_bugs/race/shutdown_flag/main.go:20 +0xADDR
```

Read the report as two conflicting accesses:

1. `Write at` points to `Stop`, which writes the shutdown flag.
2. `Previous read at` points to `Run`, which checks the same flag.
3. The goroutine creation stack shows where the worker was started.

`make test` can pass because the worker eventually stops; `-race` is what checks the memory synchronization.

The committed [race.logs](race.logs) snapshot normalizes addresses, goroutine IDs, and package duration so diffs stay stable across machines.

## One Fix

Use context cancellation or a closed channel instead of a shared boolean:

```go
func (w *Worker) Run(ctx context.Context, pollEvery time.Duration) {
	for {
		select {
		case <-time.After(pollEvery):
		case <-ctx.Done():
			return
		}
	}
}
```

If a flag is the right shape, use `atomic.Bool` or protect it with a mutex.
