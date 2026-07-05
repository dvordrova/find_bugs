# Config Pointer Data Race

This example models a config cache used by request handlers:

1. `ConfigCache.Refresh` replaces the current config pointer.
2. `ConfigCache.APIHost` reads that pointer on the request path.
3. Refresh and reads can happen concurrently.
4. The pointer field is not protected by a mutex or atomic operation, so the read and write race.

The bug is in [main.go](main.go): `current` is shared mutable state.

## Run

```sh
make run
```

Expected result:

```text
api host: api.internal
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
  github.com/dvordrova/find_bugs/race/config_pointer.(*ConfigCache).Refresh()
      github.com/dvordrova/find_bugs/race/config_pointer/main.go:18 +0xADDR

Previous read at 0xADDR by goroutine N:
  github.com/dvordrova/find_bugs/race/config_pointer.(*ConfigCache).APIHost()
      github.com/dvordrova/find_bugs/race/config_pointer/main.go:22 +0xADDR
```

Read the report as two conflicting accesses to the same pointer field:

1. `Write at` points to `Refresh`, which replaces `current`.
2. `Previous read at` points to `APIHost`, which reads `current`.
3. The goroutine creation stack shows where the refresh worker was started in the test.

`make test` can pass because the race detector is not enabled there.

The committed [race.logs](race.logs) snapshot normalizes addresses, goroutine IDs, and package duration so diffs stay stable across machines.

## One Fix

Use `atomic.Pointer[Config]` for immutable config snapshots:

```go
type ConfigCache struct {
	current atomic.Pointer[Config]
}
```

Store only fully-built immutable configs. A `sync.RWMutex` is also fine when refresh needs to update several related fields together.
