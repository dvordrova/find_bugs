# SQL Rows Not Closed

This example models a repository method that loads open invoices from a database:

1. `OpenInvoices` calls `QueryContext`.
2. It scans all rows.
3. It checks `rows.Err()`.
4. It forgets to close `rows`.

The bug is in [main.go](main.go): `*sql.Rows` keeps database resources until it is closed.

## Run

```sh
make run
```

Expected result:

```text
load open invoices
```

The program does not need a real database connection because the bug is demonstrated statically.

## Catch With sqlclosecheck Through golangci-lint

```sh
make lint
```

Expected report:

```text
main.go:20:32: Rows/Stmt/NamedStmt was not closed (sqlclosecheck)
	rows, err := s.db.QueryContext(ctx, `
	                              ^
```

Read the report as a database resource lifetime warning:

1. `QueryContext` returns `*sql.Rows`.
2. The caller owns that rows object.
3. The function returns without calling `rows.Close()`.

This example still checks `rows.Err()` on purpose. A missing iteration error check is a different bug, and `rowserrcheck` is the focused linter for that case.

`make tool-update` is a maintainer command for intentionally updating the pinned `golangci-lint` dependency.

## One Fix

Close rows immediately after the query succeeds:

```go
rows, err := s.db.QueryContext(ctx, query)
if err != nil {
	return nil, err
}
defer rows.Close()
```

Keep the `rows.Err()` check after the loop, because `Close` and iteration errors cover different failure modes.
