# Transaction Boundary

This example shows a team architecture rule:

1. Service/application code may request atomic work.
2. Opening, committing, and rolling back transactions should live in a transaction manager or unit-of-work package.
3. Business services should not grow their own transaction lifecycle around every multi-step operation.

The violation is in [internal/orders/service.go](internal/orders/service.go): `Service.Capture` starts a transaction, rolls it back in a defer, and commits it directly. The allowed shape is shown in [internal/orders/transaction/manager.go](internal/orders/transaction/manager.go), whose package path ends with `/transaction`.

## Run

```sh
make run
```

Expected result:

```text
transaction boundaries should live in transaction manager packages
run make lint to see the ruleguard report
```

The program does not need a real database. The bug is an architecture boundary issue in code shape, not a runtime failure in this tiny example.

## Catch With ruleguard

```sh
make lint
```

Expected report:

```text
internal/orders/service.go:18:13: transactionBoundary: transactions belong in transaction manager packages (transaction_boundary.go:8)
internal/orders/service.go:23:7: transactionBoundary: transactions belong in transaction manager packages (transaction_boundary.go:17)
internal/orders/service.go:33:12: transactionBoundary: transactions belong in transaction manager packages (transaction_boundary.go:16)
```

Read the report as a transaction ownership violation:

1. `service.go` is outside a package whose import path ends with `/transaction`.
2. `BeginTx` means this package owns transaction start.
3. `Rollback` and `Commit` mean this package also owns transaction lifecycle completion.

The rule lives in [rules/transaction_boundary.go](rules/transaction_boundary.go). It is type-aware: it looks for `*sql.DB.BeginTx` and `*sql.Tx.Commit`/`Rollback`, not just method names.

`make tool-update` is a maintainer command for intentionally updating pinned `ruleguard` dependencies.

## One Fix

Move transaction lifecycle into a transaction manager and pass the `*sql.Tx` into repositories or query objects:

```go
err := txManager.Within(ctx, func(ctx context.Context, tx *sql.Tx) error {
	if err := orders.MarkCaptured(ctx, tx, orderID); err != nil {
		return err
	}
	return outbox.Enqueue(ctx, tx, orderID, "order.captured")
})
```

The service still describes the business operation, but the lifecycle policy is centralized.
