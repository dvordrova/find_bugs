# Context Timeout Without Wall Clock

This example shows timeout logic that ordinary tests often miss because checking the deadline would require waiting on real time.

The production shape is a lease or reservation manager. A lease should expire after a configured TTL. The bug is that the manager accidentally creates a timeout for `ttl * 2`, so the lease stays active longer than the public contract says.

## Run It

```sh
make run
```

Expected result:

```text
lease order-42 active immediately: true
deadline check is covered by testing/synctest without waiting on wall-clock time
```

## Ordinary Test

```sh
make test
```

The ordinary test only checks that a new lease starts active. It passes even though the lease deadline is wrong.

This is the weak assertion:

```text
new lease is active immediately
```

That assertion does not prove the lease expires at the configured TTL.

## Catch It With Synctest

```sh
make lint
```

The `lint` target runs a bug-revealing test with `testing/synctest`. Inside the synctest bubble, `time.Sleep` advances fake time instead of wall-clock time.

The test checks both sides of the deadline:

1. just before `ttl`, the lease should still be active;
2. exactly at `ttl`, the lease should be expired.

Expected output:

```text
--- FAIL: TestSynctestChecksLeaseDeadline (0.00s)
    main_test.go:38: after ttl: lease is still active, want expired
FAIL
FAIL	github.com/dvordrova/find_bugs/synctest/context_timeout_without_wall_clock Xs
?   	github.com/dvordrova/find_bugs/synctest/context_timeout_without_wall_clock/internal/leases	[no test files]
FAIL
```

## Fix

Use the configured TTL directly:

```go
leaseCtx, cancel := context.WithTimeout(ctx, m.ttl)
```

## Why This Matters

Without synctest, deadline tests often choose between slow real sleeps and incomplete assertions. `testing/synctest` lets the test move through time immediately and verify timeout behavior deterministically.

This example requires Go 1.25 or newer. The module pins `go 1.26` because the repository CI runs Go 1.26.
