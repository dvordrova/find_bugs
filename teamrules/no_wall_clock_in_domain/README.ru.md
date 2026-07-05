# Wall Clock Не Должен Жить В Domain Logic

Этот пример показывает team architecture rule:

1. Domain code может работать со значениями `time.Time`.
2. Domain code не должен сам спрашивать у операционной системы, сколько сейчас времени.
3. Текущее время должно попадать в domain logic через параметр или маленькую clock abstraction.

Нарушение находится в [internal/billing/domain/subscription.go](internal/billing/domain/subscription.go): `Subscription.NeedsRenewalNotice` напрямую вызывает `time.Now()`. Такой же wall-clock вызов разрешен в [internal/billing/clock/system.go](internal/billing/clock/system.go), потому что это adapter boundary.

## Запуск

```sh
make run
```

Ожидаемый результат:

```text
renewal notice due: true
```

Программа запускается, но domain decision зависит от реального текущего времени.

## Как Поймать Через ruleguard

```sh
make lint
```

Ожидаемый отчет:

```text
internal/billing/domain/subscription.go:11:30: noWallClockInDomain: domain logic must not call time.Now directly (no_wall_clock.go:6)
```

Читай report как нарушение design boundary:

1. `subscription.go` находится в package, чей import path заканчивается на `/domain`.
2. Выражение - `time.Now()`.
3. Team rule держит wall-clock reads в adapters, handlers, jobs или composition roots.

Правило лежит в [rules/no_wall_clock.go](rules/no_wall_clock.go). Оно специально узкое: оно не запрещает `time.Time` или `time.Duration` в domain code, только прямые wall-clock reads.

`make tool-update` - maintainer command для осознанного обновления pinned `ruleguard` dependencies.

## Одно Исправление

Передавать текущее время в domain method:

```go
func (s Subscription) NeedsRenewalNotice(now time.Time, window time.Duration) bool {
	remaining := s.RenewsAt.Sub(now)
	return remaining > 0 && remaining <= window
}
```

Application service, handler, worker или clock adapter может решить, что такое `now`. Domain logic остается deterministic и ее проще тестировать.
