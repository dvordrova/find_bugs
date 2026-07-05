# Goroutine Leak При Timeout На Channel

Этот пример показывает сервис, который запускает асинхронную работу и ждет либо результат, либо timeout запроса:

1. `SendWelcomeEmail` просит `EmailGateway` отправить welcome email.
2. Gateway запускает goroutine и возвращает channel с результатом.
3. Request context истекает раньше, чем goroutine отправляет результат.
4. Caller возвращается, channel больше никто не читает, и goroutine навсегда блокируется на send.

Баг находится в [main.go](main.go): `EmailGateway.Deliver` ничего не знает про cancellation и отправляет в unbuffered channel после того, как caller уже мог перестать ждать.

## Запуск

```sh
make run
```

Ожидаемый результат:

```text
welcome email timed out
```

Самостоятельная программа сразу завершается, поэтому leaked goroutine умирает вместе с процессом. В server или test process такая goroutine осталась бы заблокированной.

## Как поймать через goleak

`go.uber.org/goleak` проверяет неожиданные goroutine после завершения тестов. В этом примере leak check спрятан за build tag `leakcheck`: обычные тесты показывают business behavior, а `make lint` демонстрирует detector.

```sh
make lint
```

Форма ожидаемого отчета:

```text
PASS
goleak: Errors on successful test run: found unexpected goroutines:
[Goroutine N in state chan send, with github.com/dvordrova/find_bugs/goleak/channel_timeout_leak.EmailGateway.Deliver.func1 on top of the stack:
github.com/dvordrova/find_bugs/goleak/channel_timeout_leak.EmailGateway.Deliver.func1()
	github.com/dvordrova/find_bugs/goleak/channel_timeout_leak/main.go:36 +0xADDR
created by github.com/dvordrova/find_bugs/goleak/channel_timeout_leak.EmailGateway.Deliver in goroutine N
	github.com/dvordrova/find_bugs/goleak/channel_timeout_leak/main.go:34 +0xADDR
]
FAIL	github.com/dvordrova/find_bugs/goleak/channel_timeout_leak Xs
FAIL
```

Отчет лучше читать с состояния goroutine:

1. `state chan send` означает, что goroutine заблокирована при попытке отправить значение в channel.
2. `EmailGateway.Deliver.func1` указывает на worker goroutine, которую запустил gateway.
3. Строка stack в `main.go` указывает на send в `receipts`.
4. Блок `created by` показывает, где эта goroutine была создана.

`make test` запускает тот же timeout test без tag `leakcheck`. Он проходит, потому что функция возвращает ожидаемую timeout error; только leak detector замечает, что background work остался жить.

Committed snapshot [lint.logs](lint.logs) нормализует goroutine IDs, runtime addresses и package duration, чтобы diffs были стабильными на разных машинах.

`make tool-update` - maintainer-команда для осознанного обновления pinned dependency `goleak`.

## Один из вариантов исправления

Передать context в worker и дать send проиграть cancellation:

```go
func (g EmailGateway) Deliver(ctx context.Context, address string) <-chan Receipt {
	receipts := make(chan Receipt)

	go func() {
		time.Sleep(g.Latency)
		receipt := Receipt{Address: address, DeliveryID: "welcome-001"}

		select {
		case receipts <- receipt:
		case <-ctx.Done():
		}
	}()

	return receipts
}
```

Buffered channel размера 1 тоже может убрать конкретно этот blocked send, но cancellation обычно лучше как API boundary: worker может остановиться, когда request уже никому не нужен.
