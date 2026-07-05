# Unbuffered Send After Timeout

This example shows a request timeout that leaves a worker blocked on an unbuffered result channel.

The production shape is common around legacy clients that do not accept `context.Context`. The service starts a worker, waits for either the worker result or a timeout, and returns `DeadlineExceeded` when the timeout wins. The bug is that the worker still sends to an unbuffered channel after the caller has left.

## Run It

```sh
make run
```

Expected result:

```text
lookup timed out: true
run make lint to expose the blocked worker send without waiting on wall-clock time
```

## Ordinary Test

```sh
make test
```

The ordinary test only checks that the caller receives `DeadlineExceeded`. It passes even though a background worker can still be left behind.

This is the weak assertion:

```text
timeout is returned to the caller
```

That assertion does not prove the worker exited.

## Catch It With Synctest

```sh
make lint
```

The `lint` target runs a bug-revealing test with `testing/synctest`. Inside the synctest bubble, fake time advances to the service timeout and then to the legacy client completion. The late worker send has no receiver.

Expected output:

```text
--- FAIL: TestSynctestFindsBlockedSendAfterTimeout (0.00s)
    main_test.go:38: synctest detected blocked worker send: deadlock: main bubble goroutine has exited but blocked goroutines remain
FAIL
FAIL	github.com/dvordrova/find_bugs/synctest/unbuffered_send_after_timeout Xs
?   	github.com/dvordrova/find_bugs/synctest/unbuffered_send_after_timeout/internal/pricing	[no test files]
FAIL
```

## Fix

Use a buffered result channel so a late result does not block the worker:

```go
results := make(chan lookupResult, 1)
```

For clients that support cancellation, also pass a cancellable context into the client and make the worker select on `ctx.Done()`.

## Why This Matters

Without synctest, tests for timeout races often rely on real sleeps or only assert the return value. `testing/synctest` lets the test move through the timeout and the late worker completion immediately.

This example requires Go 1.25 or newer. The module pins `go 1.26` because the repository CI runs Go 1.26.
