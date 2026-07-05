# Message Order Assumption

This example shows channel/message code that only works when related events arrive in one order.

The production shape is a read-model projector. `AccountOpened` creates an account in the projection. `CreditReserved` records reserved credit for that account. The bug is that `CreditReserved` is silently dropped if it arrives before `AccountOpened`.

Real brokers, channel pipelines, retries, and backfills can reorder related messages unless the protocol explicitly prevents it.

## Run It

```sh
make run
```

Expected result:

```text
broker order reserved cents: 500
reordered delivery reserved cents: 0
```

## Ordinary Test

```sh
make test
```

The ordinary test only checks the order the developer expected:

```text
AccountOpened, then CreditReserved
```

That passes, but it does not prove the projector handles reordered delivery.

## Catch It With Order Permutations

```sh
make lint
```

The `lint` target runs a bug-revealing test that applies the same logical messages in both orders. This is the small version of the GFuzz idea: change message order before adding a heavy scheduler/fuzzer.

Expected output:

```text
--- FAIL: TestMessageOrderPermutations (0.00s)
    --- FAIL: TestMessageOrderPermutations/credit_before_account (0.00s)
        main_test.go:48: reserved cents after credit_before_account = 0, want 500
FAIL
FAIL	github.com/dvordrova/find_bugs/concurrency/message_order_assumption Xs
?   	github.com/dvordrova/find_bugs/concurrency/message_order_assumption/internal/projection	[no test files]
FAIL
```

## Fix

Make the protocol or projection robust to reordering. One simple projection-side mitigation is to keep pending reservations until the account appears:

```go
case CreditReserved:
	if !p.accounts[event.AccountID] {
		p.pending[event.AccountID] += event.Cents
		return
	}
	p.reserved[event.AccountID] += event.Cents
```

The better fix depends on the production contract: partitioning guarantees, idempotency keys, sequence numbers, or a state machine that can represent out-of-order facts.
