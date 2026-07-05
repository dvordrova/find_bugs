# Infrastructure Imports Не Должны Жить В Domain

Этот пример показывает team architecture rule:

1. Domain packages должны моделировать business state и behavior.
2. Database, HTTP, queue и logging packages должны жить в adapters или application services.
3. Domain type не должен требовать `database/sql`, `net/http` или `log/slog`, чтобы скомпилироваться.

Нарушение находится в [internal/billing/domain/invoice.go](internal/billing/domain/invoice.go): domain entity использует `sql.NullString` и `http.Header`. Разрешенная форма показана в [internal/billing/postgres/invoice_repository.go](internal/billing/postgres/invoice_repository.go), где `database/sql` находится в adapter package.

## Запуск

```sh
make run
```

Ожидаемый результат:

```text
invoice note: ""
run make lint to see the depguard report
```

Программа запускается, но infrastructure concerns уже протекли в domain model.

## Как Поймать Через depguard

```sh
make lint
```

Ожидаемый отчет:

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

Читай report как нарушение import boundary:

1. `invoice.go` находится внутри `domain` package path.
2. `database/sql` и `net/http` запрещены для domain files.
3. Те же imports все еще могут быть разрешены в adapter packages.

Правило лежит в [.golangci.yaml](.golangci.yaml). Оно использует `depguard`, потому что import-level dependency rules лучше выражаются через import boundary, а не через call-level matching.

`make tool-update` - maintainer command для осознанного обновления pinned `golangci-lint` dependencies.

## Одно Исправление

Переводить infrastructure types на boundary:

```go
type Invoice struct {
	ID           string
	CustomerNote string
}
```

Repository может конвертировать `sql.NullString` в domain value перед возвратом.
