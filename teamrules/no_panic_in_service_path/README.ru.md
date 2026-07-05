# Panic Не Должен Жить В Service Path

Этот пример показывает team architecture rule:

1. Service/application code должен возвращать обычные business и dependency failures как errors.
2. `panic` остается для programmer errors или process startup failures.
3. Service method не должен ронять процесс из-за declined/rejected payment.

Нарушение находится в [internal/payments/service/payment.go](internal/payments/service/payment.go): `PaymentService.Capture` паникует, когда amount превышает capture limit. Happy-path unit test все равно проходит, поэтому team rule ловит проблему до редкого production path.

## Запуск

```sh
make run
```

Ожидаемый результат:

```text
service panicked: payment pay_42 exceeds capture limit
```

Программа делает recover только чтобы пример мог напечатать failure. В реальном service process такой panic обычно потеряет request, а иногда и весь worker, если panic уйдет достаточно далеко.

## Как Поймать Через ruleguard

```sh
make lint
```

Ожидаемый отчет:

```text
internal/payments/service/payment.go:27:3: noPanicInServicePath: service code must return errors instead of panicking (no_panic.go:6)
```

Читай report как нарушение service boundary:

1. `payment.go` находится в package, чей import path заканчивается на `/service`.
2. Выражение - `panic(...)`.
3. Team rule ожидает, что service methods возвращают errors для нормальных failure modes.

Правило лежит в [rules/no_panic.go](rules/no_panic.go). Оно специально узкое: оно не запрещает `panic` во всех packages, только в service packages.

`make tool-update` - maintainer command для осознанного обновления pinned `ruleguard` dependencies.

## Одно Исправление

Вернуть error вместо panic:

```go
if payment.Amount > s.limit {
	return fmt.Errorf("payment %s exceeds capture limit", payment.ID)
}
```

Caller может решить, retry это, reject request, metrics или user-facing error.
