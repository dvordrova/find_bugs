# Message Order Assumption

Этот пример показывает channel/message code, который работает только когда связанные events приходят в одном порядке.

Production shape - read-model projector. `AccountOpened` создает account в projection. `CreditReserved` записывает reserved credit для этого account. Баг в том, что `CreditReserved` молча выбрасывается, если приходит до `AccountOpened`.

Real brokers, channel pipelines, retries и backfills могут менять порядок связанных messages, если protocol явно это не запрещает.

## Запуск

```sh
make run
```

Ожидаемый результат:

```text
broker order reserved cents: 500
reordered delivery reserved cents: 0
```

## Обычный Тест

```sh
make test
```

Обычный test проверяет только порядок, который developer ожидал:

```text
AccountOpened, then CreditReserved
```

Это проходит, но не доказывает, что projector переживает reordered delivery.

## Как Поймать Через Order Permutations

```sh
make lint
```

`lint` target запускает bug-revealing test, который применяет те же logical messages в обоих порядках. Это маленькая версия идеи GFuzz: менять message order до того, как добавлять тяжелый scheduler/fuzzer.

Ожидаемый вывод:

```text
--- FAIL: TestMessageOrderPermutations (0.00s)
    --- FAIL: TestMessageOrderPermutations/credit_before_account (0.00s)
        main_test.go:48: reserved cents after credit_before_account = 0, want 500
FAIL
FAIL	github.com/dvordrova/find_bugs/concurrency/message_order_assumption Xs
?   	github.com/dvordrova/find_bugs/concurrency/message_order_assumption/internal/projection	[no test files]
FAIL
```

## Исправление

Сделать protocol или projection устойчивыми к reordered delivery. Простая projection-side mitigation - хранить pending reservations, пока account не появится:

```go
case CreditReserved:
	if !p.accounts[event.AccountID] {
		p.pending[event.AccountID] += event.Cents
		return
	}
	p.reserved[event.AccountID] += event.Cents
```

Лучшее исправление зависит от production contract: partitioning guarantees, idempotency keys, sequence numbers или state machine, которая умеет выражать out-of-order facts.
