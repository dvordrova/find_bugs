# Goroutine Leak Из-За Неиспользованного Context

Этот пример показывает cache warmer, который принимает request или application context:

1. `SessionCache.Warm` сохраняет initial value.
2. Он запускает goroutine, которая периодически обновляет cache для tenant.
3. Caller отменяет context.
4. Goroutine продолжает работать, потому что worker слушает `context.Background().Done()` вместо context, который передал caller.

Баг находится в [main.go](main.go): у `Warm` есть параметр `ctx`, но goroutine не использует `ctx.Done()`.

## Запуск

```sh
make run
```

Ожидаемый результат: cache один раз прогрет, программа завершается.

## Как поймать через goleak

```sh
make lint
```

Форма ожидаемого отчета:

```text
PASS
goleak: Errors on successful test run: found unexpected goroutines:
[Goroutine N in state select, with github.com/dvordrova/find_bugs/goleak/context_not_cancelled.(*SessionCache).Warm.func1 on top of the stack:
github.com/dvordrova/find_bugs/goleak/context_not_cancelled.(*SessionCache).Warm.func1()
	github.com/dvordrova/find_bugs/goleak/context_not_cancelled/main.go:33 +0xADDR
created by github.com/dvordrova/find_bugs/goleak/context_not_cancelled.(*SessionCache).Warm in goroutine N
	github.com/dvordrova/find_bugs/goleak/context_not_cancelled/main.go:28 +0xADDR
]
FAIL	github.com/dvordrova/find_bugs/goleak/context_not_cancelled Xs
FAIL
```

Отчет лучше читать с состояния goroutine и stack:

1. `state select` означает, что goroutine все еще ждет внутри `select`.
2. `SessionCache.Warm.func1` указывает на background warmer.
3. Строка stack в `main.go` указывает на `select`, который должен был слушать context caller.
4. Блок `created by` показывает API call, который запустил worker.

`make test` проходит, потому что initial cache warm работает. Leak виден только после завершения теста, когда `goleak` проверяет оставшиеся goroutine.

Committed snapshot [lint.logs](lint.logs) нормализует goroutine IDs, runtime addresses и package duration, чтобы diffs были стабильными на разных машинах.

## Один из вариантов исправления

Использовать context caller внутри worker:

```go
case <-ctx.Done():
	return
```

Для долгоживущих application components еще один хороший дизайн - явный lifecycle `Start(ctx)` / `Stop()`, которым владеет application, а не короткий request context.
