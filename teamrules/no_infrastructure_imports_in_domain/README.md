# No Infrastructure Imports In Domain

This example shows a team architecture rule:

1. Domain packages should model business state and behavior.
2. Database, HTTP, queue, and logging packages belong in adapters or application services.
3. A domain type should not need `database/sql`, `net/http`, or `log/slog` just to compile.

The violation is in [internal/billing/domain/invoice.go](internal/billing/domain/invoice.go): the domain entity uses `sql.NullString` and `http.Header`. The allowed shape is shown in [internal/billing/postgres/invoice_repository.go](internal/billing/postgres/invoice_repository.go), where `database/sql` belongs in an adapter package.

## Run

```sh
make run
```

Expected result:

```text
invoice note: ""
run make lint to see the depguard report
```

The program runs, but infrastructure concerns have leaked into the domain model.

## Catch With depguard

```sh
make lint
```

Expected report:

```text
internal/billing/domain/invoice.go:4:2: import 'database/sql' is not allowed from list 'domain': domain packages must not depend on database/sql (depguard)
	"database/sql"
	^
internal/billing/domain/invoice.go:5:2: import 'net/http' is not allowed from list 'domain': domain packages must not depend on HTTP transport types (depguard)
	"net/http"
	^
2 issues:
* depguard: 2
```

Read the report as an import boundary violation:

1. `invoice.go` is inside a `domain` package path.
2. `database/sql` and `net/http` are denied for domain files.
3. The same imports can still be allowed in adapter packages.

The rule lives in [.golangci.yaml](.golangci.yaml). It uses `depguard`, which is a better fit than call-level matching when the rule is about package dependencies.

`make tool-update` is a maintainer command for intentionally updating pinned `golangci-lint` dependencies.

## One Fix

Translate infrastructure types at the boundary:

```go
type Invoice struct {
	ID           string
	CustomerNote string
}
```

The repository can convert `sql.NullString` into a domain value before returning it.
