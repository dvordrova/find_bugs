# Force sqlc Query Layer

Этот пример показывает team architecture rule:

1. SQL должен жить и ревьюиться в generated query layer.
2. Application и repository code должны вызывать generated query methods.
3. Прямые `database/sql` вызовы вне этого слоя считаются bypass.

Нарушение находится в [internal/orders/repository.go](internal/orders/repository.go): `Repository.ListPending` пишет raw SQL через `*sql.DB.QueryContext`. Разрешенная форма показана в [internal/store/sqlc/orders.sql.go](internal/store/sqlc/orders.sql.go), чей package path заканчивается на `/sqlc`.

## Запуск

```sh
make run
```

Ожидаемый результат:

```text
team rule: application code should use the generated sqlc query layer
run make lint to see the ruleguard report
```

Примеру не нужна настоящая database. Смысл в том, чтобы поймать bypass query API до того, как он расползется по codebase.

## Как Поймать Через ruleguard

```sh
make lint
```

Ожидаемый отчет:

```text
internal/orders/repository.go:20:15: forceSQLCQueryLayer: database/sql calls must go through generated sqlc packages (force_sqlc.go:10)
```

Читай report как нарушение ownership для queries:

1. `repository.go` находится вне package, чей import path заканчивается на `/sqlc`.
2. Receiver имеет type `*sql.DB`.
3. Method - query/exec method, который должен быть wrapped generated query package.

Правило лежит в [rules/force_sqlc.go](rules/force_sqlc.go). Оно специально узкое и type-aware. В большом service его стоит комбинировать с import rules вроде `depguard`, если команда хочет еще и запретить imports `database/sql` вне конкретных packages.

`make tool-update` - maintainer command для осознанного обновления pinned `ruleguard` dependencies.

## Одно Исправление

Инжектить и использовать generated query object:

```go
type Repository struct {
	queries *sqlc.Queries
}

func (r Repository) ListPending(ctx context.Context, limit int) ([]sqlc.Order, error) {
	return r.queries.ListPendingOrders(ctx, limit)
}
```

Так query text остается в одном reviewed layer, а application code зависит от меньшего API.
