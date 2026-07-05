# Channel Timeout Goroutine Leak

This example models a service that starts async work and waits for either a result or a request timeout:

1. `SendWelcomeEmail` asks `EmailGateway` to deliver a welcome email.
2. The gateway starts a goroutine and returns a result channel.
3. The request context times out before the goroutine sends the result.
4. The caller returns, nobody receives from the channel, and the goroutine blocks forever on send.

The bug is in [main.go](main.go): `EmailGateway.Deliver` does not know about cancellation and sends to an unbuffered channel after the caller may have stopped waiting.

## Run

```sh
make run
```

Expected result:

```text
welcome email timed out
```

The standalone program exits immediately, so the leaked goroutine dies with the process. In a server or test process, the same goroutine would stay blocked.

## Catch With goleak

`go.uber.org/goleak` checks for unexpected goroutines after tests finish. This example keeps the leak check behind the `leakcheck` build tag so normal tests can still show the business behavior, while `make lint` demonstrates the detector.

```sh
make lint
```

Expected report shape:

```text
PASS
goleak: Errors on successful test run: found unexpected goroutines:
[Goroutine N in state chan send, with github.com/dvordrova/find_bugs/goleak/channel_timeout_leak.EmailGateway.Deliver.func1 on top of the stack:
github.com/dvordrova/find_bugs/goleak/channel_timeout_leak.EmailGateway.Deliver.func1()
	github.com/dvordrova/find_bugs/goleak/channel_timeout_leak/main.go:36 +0xADDR
created by github.com/dvordrova/find_bugs/goleak/channel_timeout_leak.EmailGateway.Deliver in goroutine N
	github.com/dvordrova/find_bugs/goleak/channel_timeout_leak/main.go:34 +0xADDR
]
FAIL	github.com/dvordrova/find_bugs/goleak/channel_timeout_leak Xs
FAIL
```

Read the report from the goroutine state first:

1. `state chan send` means the goroutine is blocked while trying to send to a channel.
2. `EmailGateway.Deliver.func1` points to the worker goroutine started by the gateway.
3. The stack line in `main.go` points to the send into `receipts`.
4. The `created by` section shows where that goroutine was started.

`make test` runs the same timeout test without the `leakcheck` tag. It passes because the function returns the expected timeout error; only the leak detector notices that background work was left behind.

The committed [lint.logs](lint.logs) snapshot normalizes goroutine IDs, runtime addresses, and package duration so diffs stay stable across machines.

`make tool-update` is a maintainer command for intentionally updating the pinned `goleak` dependency.

## One Fix

Pass the context into the worker and let the send lose to cancellation:

```go
func (g EmailGateway) Deliver(ctx context.Context, address string) <-chan Receipt {
	receipts := make(chan Receipt)

	go func() {
		time.Sleep(g.Latency)
		receipt := Receipt{Address: address, DeliveryID: "welcome-001"}

		select {
		case receipts <- receipt:
		case <-ctx.Done():
		}
	}()

	return receipts
}
```

A buffered channel of size 1 can also prevent this exact blocked send, but cancellation is usually the stronger API boundary: it lets the worker stop when the request is gone.
