# DDD Repository Boundary

Этот пример показывает team architecture rule, а не runtime bug:

1. Application/service code должен вызывать repository interfaces.
2. Repository packages могут использовать `database/sql`.
3. Service code не должен напрямую ходить в `*sql.DB`.

Нарушение находится в [internal/customer/service.go](internal/customer/service.go): `ExportService.ActiveCustomers` вызывает `s.db.QueryContext` из service package. Такой же database call разрешен в [internal/customer/repository/repository.go](internal/customer/repository/repository.go).

## Запуск

```sh
make run
```

Ожидаемый результат:

```text
team rule: database/sql belongs in repository packages
```

Реальное подключение к базе не нужно, потому что правило проверяется статически.

## Как поймать через ruleguard

```sh
make lint
```

Ожидаемый вывод:

```text
internal/customer/service.go:22:15: dddRepository: database/sql calls belong in repository packages (ddd_repository.go:10)
```

Отчет лучше читать как нарушение architecture boundary:

1. `service.go` находится не в package, import path которого заканчивается на `/repository`.
2. `s.db` имеет тип `*sql.DB`.
3. `QueryContext` - database operation, которую это team rule разрешает только в repository packages.

Правило лежит в [rules/ddd_repository.go](rules/ddd_repository.go). Оно type-aware: проверяет `*sql.DB` и `*sql.Tx`, а не просто текст `QueryContext`.

`make tool-update` - maintainer-команда для осознанного обновления pinned `ruleguard` dependencies.

## Один из вариантов исправления

Спрятать database access за repository dependency:

```go
type CustomerRepository interface {
	ActiveCustomers(context.Context) ([]Customer, error)
}

type ExportService struct {
	customers CustomerRepository
}
```

Service все еще выражает use-case logic, но SQL остается в repository packages, где команда ожидает его видеть.
