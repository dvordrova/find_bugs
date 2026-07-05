# Force sqlc Query Layer

This example shows a team architecture rule:

1. SQL should be written and reviewed in the generated query layer.
2. Application and repository code should call generated query methods.
3. Direct `database/sql` calls outside that layer are treated as a bypass.

The violation is in [internal/orders/repository.go](internal/orders/repository.go): `Repository.ListPending` writes raw SQL through `*sql.DB.QueryContext`. The allowed shape is shown in [internal/store/sqlc/orders.sql.go](internal/store/sqlc/orders.sql.go), whose package path ends with `/sqlc`.

## Run

```sh
make run
```

Expected result:

```text
team rule: application code should use the generated sqlc query layer
run make lint to see the ruleguard report
```

The program does not need a real database. The point is to catch code that bypasses the query API before it spreads.

## Catch With ruleguard

```sh
make lint
```

Expected report:

```text
internal/orders/repository.go:20:15: forceSQLCQueryLayer: database/sql calls must go through generated sqlc packages (force_sqlc.go:10)
```

Read the report as a query ownership violation:

1. `repository.go` is outside a package whose import path ends with `/sqlc`.
2. The receiver has type `*sql.DB`.
3. The method is a query/exec method that should be wrapped by the generated query package.

The rule lives in [rules/force_sqlc.go](rules/force_sqlc.go). It is intentionally narrow and type-aware. In a larger service, combine this with import rules such as `depguard` if the team also wants to ban `database/sql` imports outside specific packages.

`make tool-update` is a maintainer command for intentionally updating pinned `ruleguard` dependencies.

## One Fix

Inject and use the generated query object:

```go
type Repository struct {
	queries *sqlc.Queries
}

func (r Repository) ListPending(ctx context.Context, limit int) ([]sqlc.Order, error) {
	return r.queries.ListPendingOrders(ctx, limit)
}
```

This keeps query text in one reviewed layer and makes application code depend on a smaller API.
