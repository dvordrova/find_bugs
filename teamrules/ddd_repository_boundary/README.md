# DDD Repository Boundary

This example shows a team architecture rule, not a runtime bug:

1. Application/service code should call repository interfaces.
2. Repository packages are allowed to use `database/sql`.
3. Service code should not reach into `*sql.DB` directly.

The violation is in [internal/customer/service.go](internal/customer/service.go): `ExportService.ActiveCustomers` calls `s.db.QueryContext` from the service package. The same database call is allowed in [internal/customer/repository/repository.go](internal/customer/repository/repository.go).

## Run

```sh
make run
```

Expected result:

```text
team rule: database/sql belongs in repository packages
```

The program does not need a real database connection because the rule is checked statically.

## Catch With ruleguard

```sh
make lint
```

Expected report:

```text
internal/customer/service.go:22:15: dddRepository: database/sql calls belong in repository packages (ddd_repository.go:10)
```

Read the report as an architecture boundary violation:

1. `service.go` is not in a package whose import path ends with `/repository`.
2. `s.db` has type `*sql.DB`.
3. `QueryContext` is a database operation that this team rule allows only in repository packages.

The rule lives in [rules/ddd_repository.go](rules/ddd_repository.go). It is type-aware: it checks `*sql.DB` and `*sql.Tx`, not just the text `QueryContext`.

`make tool-update` is a maintainer command for intentionally updating pinned `ruleguard` dependencies.

## One Fix

Move database access behind a repository dependency:

```go
type CustomerRepository interface {
	ActiveCustomers(context.Context) ([]Customer, error)
}

type ExportService struct {
	customers CustomerRepository
}
```

The service can still express use-case logic, but SQL stays in repository packages where the team expects it.
