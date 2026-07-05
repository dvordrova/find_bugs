# No Panic In Service Path

This example shows a team architecture rule:

1. Service/application code should return ordinary business and dependency failures as errors.
2. `panic` is reserved for programmer errors or process startup failures.
3. A service method should not crash the process because a payment was declined or rejected.

The violation is in [internal/payments/service/payment.go](internal/payments/service/payment.go): `PaymentService.Capture` panics when the amount exceeds a capture limit. A happy-path unit test still passes, so the team rule catches the problem before a rare production path hits it.

## Run

```sh
make run
```

Expected result:

```text
service panicked: payment pay_42 exceeds capture limit
```

The program recovers only so the example can print the failure. A real service process would normally lose the request, and sometimes the whole worker, depending on where the panic escapes.

## Catch With ruleguard

```sh
make lint
```

Expected report:

```text
internal/payments/service/payment.go:27:3: noPanicInServicePath: service code must return errors instead of panicking (no_panic.go:6)
```

Read the report as a service boundary violation:

1. `payment.go` is in a package whose import path ends with `/service`.
2. The expression is `panic(...)`.
3. This team rule expects service methods to return errors for normal failure modes.

The rule lives in [rules/no_panic.go](rules/no_panic.go). It is intentionally narrow: it does not ban `panic` in every package, only in service packages.

`make tool-update` is a maintainer command for intentionally updating pinned `ruleguard` dependencies.

## One Fix

Return an error instead of panicking:

```go
if payment.Amount > s.limit {
	return fmt.Errorf("payment %s exceeds capture limit", payment.ID)
}
```

The caller can decide whether to retry, reject the request, emit metrics, or return a user-facing error.
