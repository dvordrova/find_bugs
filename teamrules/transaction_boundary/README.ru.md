# Transaction Boundary

Этот пример показывает team architecture rule:

1. Service/application code может запрашивать atomic work.
2. Открытие, commit и rollback transactions должны жить в transaction manager или unit-of-work package.
3. Business services не должны выращивать собственный transaction lifecycle вокруг каждой multi-step operation.

Нарушение находится в [internal/orders/service.go](internal/orders/service.go): `Service.Capture` сам начинает transaction, делает rollback в defer и напрямую вызывает commit. Разрешенная форма показана в [internal/orders/transaction/manager.go](internal/orders/transaction/manager.go), чей package path заканчивается на `/transaction`.

## Запуск

```sh
make run
```

Ожидаемый результат:

```text
transaction boundaries should live in transaction manager packages
run make lint to see the ruleguard report
```

Примеру не нужна настоящая database. Баг здесь в architecture boundary и форме кода, а не в runtime failure.

## Как Поймать Через ruleguard

```sh
make lint
```

Ожидаемый отчет:

```text
internal/orders/service.go:18:13: transactionBoundary: transactions belong in transaction manager packages (transaction_boundary.go:8)
internal/orders/service.go:23:7: transactionBoundary: transactions belong in transaction manager packages (transaction_boundary.go:17)
internal/orders/service.go:33:12: transactionBoundary: transactions belong in transaction manager packages (transaction_boundary.go:16)
```

Читай report как нарушение ownership для transactions:

1. `service.go` находится вне package, чей import path заканчивается на `/transaction`.
2. `BeginTx` значит, что этот package владеет стартом transaction.
3. `Rollback` и `Commit` значит, что этот package владеет завершением transaction lifecycle.

Правило лежит в [rules/transaction_boundary.go](rules/transaction_boundary.go). Оно type-aware: ищет `*sql.DB.BeginTx` и `*sql.Tx.Commit`/`Rollback`, а не просто method names.

`make tool-update` - maintainer command для осознанного обновления pinned `ruleguard` dependencies.

## Одно Исправление

Перенести transaction lifecycle в transaction manager и передавать `*sql.Tx` в repositories или query objects:

```go
err := txManager.Within(ctx, func(ctx context.Context, tx *sql.Tx) error {
	if err := orders.MarkCaptured(ctx, tx, orderID); err != nil {
		return err
	}
	return outbox.Enqueue(ctx, tx, orderID, "order.captured")
})
```

Service все еще описывает business operation, но lifecycle policy централизована.
