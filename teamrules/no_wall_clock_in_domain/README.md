# No Wall Clock In Domain Logic

This example shows a team architecture rule:

1. Domain code may work with `time.Time` values.
2. Domain code should not ask the operating system what time it is.
3. Current time should enter domain logic through a parameter or a small clock abstraction.

The violation is in [internal/billing/domain/subscription.go](internal/billing/domain/subscription.go): `Subscription.NeedsRenewalNotice` calls `time.Now()` directly. The same wall-clock call is allowed in [internal/billing/clock/system.go](internal/billing/clock/system.go), which is an adapter boundary.

## Run

```sh
make run
```

Expected result:

```text
renewal notice due: true
```

The program runs, but the domain decision depends on the real current time.

## Catch With ruleguard

```sh
make lint
```

Expected report:

```text
internal/billing/domain/subscription.go:11:30: noWallClockInDomain: domain logic must not call time.Now directly (no_wall_clock.go:6)
```

Read the report as a design boundary violation:

1. `subscription.go` is in a package whose import path ends with `/domain`.
2. The expression is `time.Now()`.
3. This team rule keeps wall-clock reads in adapters, handlers, jobs, or composition roots.

The rule lives in [rules/no_wall_clock.go](rules/no_wall_clock.go). It is intentionally narrow: it does not ban `time.Time` or `time.Duration` in domain code, only direct wall-clock reads.

`make tool-update` is a maintainer command for intentionally updating pinned `ruleguard` dependencies.

## One Fix

Pass current time into the domain method:

```go
func (s Subscription) NeedsRenewalNotice(now time.Time, window time.Duration) bool {
	remaining := s.RenewsAt.Sub(now)
	return remaining > 0 && remaining <= window
}
```

An application service, handler, worker, or clock adapter can decide what `now` means. The domain logic stays deterministic and easier to test.
