# Context AfterFunc Negative Assertion

This example shows a cancellation hook that writes an audit record too early.

The production shape is common in request cleanup code: a service registers a callback that should run when a context is canceled. The bug is that the callback is also started during registration, so the audit record exists before cancellation.

## Run It

```sh
make run
```

Expected behavior:

```text
before cancel: 0 audit record(s)
after cancel: 1 audit record(s)
```

Actual behavior is different because the callback runs before cancellation.

```text
before cancel: 1 audit record(s)
after cancel: 2 audit record(s)
```

## Ordinary Test

```sh
make test
```

The ordinary test only checks that an audit record eventually exists after cancellation. It passes even when the record was written before cancellation.

This is the weak assertion:

```text
after cancel, at least one audit record exists
```

That assertion does not prove the hook waited for cancellation.

## Catch It With Synctest

```sh
make lint
```

The `lint` target runs a bug-revealing test with `testing/synctest`. The test enters a synctest bubble, registers the callback, then calls `synctest.Wait()` before canceling the context.

`synctest.Wait()` waits until the goroutine started inside the bubble has either finished or is durably blocked. This makes the negative assertion deterministic:

```text
before cancel, no audit record exists
```

Expected output:

```text
--- FAIL: TestSynctestChecksBeforeCancel (0.00s)
    main_test.go:33: before cancel: audit records = 1, want 0
FAIL
FAIL	github.com/dvordrova/find_bugs/synctest/context_afterfunc_negative_assertion Xs
FAIL
```

## Fix

Remove the eager goroutine and only register the `context.AfterFunc` callback:

```go
func NotifyWhenCanceled(ctx context.Context, sink *AuditSink, accountID string) func() bool {
	return context.AfterFunc(ctx, func() {
		sink.Record(accountID)
	})
}
```

## Why This Matters

Without synctest, tests for "nothing happened yet" often use sleeps. Those tests are either slow, flaky, or both. `testing/synctest` lets the test wait for the concurrent code to settle without waiting on wall-clock time.

This example requires Go 1.25 or newer. The module pins `go 1.26` because the repository CI runs Go 1.26.
