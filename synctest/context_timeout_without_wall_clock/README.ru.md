# Context Timeout Без Wall Clock

Этот пример показывает timeout logic, которую обычные тесты часто пропускают, потому что проверка deadline требует ждать реальное время.

Production-сценарий - lease/reservation manager. Lease должен истечь после configured TTL. Баг в том, что manager случайно создает timeout на `ttl * 2`, поэтому lease живет дольше публичного контракта.

## Запуск

```sh
make run
```

Ожидаемый результат:

```text
lease order-42 active immediately: true
deadline check is covered by testing/synctest without waiting on wall-clock time
```

## Обычный Тест

```sh
make test
```

Обычный тест проверяет только то, что новый lease сначала active. Он проходит, хотя deadline у lease неверный.

Слабый assertion:

```text
new lease is active immediately
```

Он не доказывает, что lease истекает ровно на configured TTL.

## Как Поймать Через Synctest

```sh
make lint
```

Target `lint` запускает bug-revealing test через `testing/synctest`. Внутри synctest bubble `time.Sleep` двигает fake time, а не wall-clock time.

Тест проверяет обе стороны deadline:

1. прямо перед `ttl` lease еще active;
2. ровно на `ttl` lease уже expired.

Ожидаемый вывод:

```text
--- FAIL: TestSynctestChecksLeaseDeadline (0.00s)
    main_test.go:38: after ttl: lease is still active, want expired
FAIL
FAIL	github.com/dvordrova/find_bugs/synctest/context_timeout_without_wall_clock Xs
?   	github.com/dvordrova/find_bugs/synctest/context_timeout_without_wall_clock/internal/leases	[no test files]
FAIL
```

## Исправление

Использовать configured TTL напрямую:

```go
leaseCtx, cancel := context.WithTimeout(ctx, m.ttl)
```

## Почему Это Важно

Без synctest deadline tests часто выбирают между медленными real sleeps и неполными assertions. `testing/synctest` позволяет мгновенно пройти через время и детерминированно проверить timeout behavior.

Пример требует Go 1.25 или новее. Модуль фиксирует `go 1.26`, потому что CI репозитория запускается на Go 1.26.
