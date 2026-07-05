# Unbuffered Send After Timeout

Этот пример показывает request timeout, после которого worker остается заблокированным на unbuffered result channel.

Production shape часто встречается вокруг legacy clients, которые не принимают `context.Context`. Service запускает worker, ждет worker result или timeout, и возвращает `DeadlineExceeded`, когда timeout выигрывает. Баг в том, что worker позже отправляет результат в unbuffered channel, хотя caller уже ушел.

## Запуск

```sh
make run
```

Ожидаемый результат:

```text
lookup timed out: true
run make lint to expose the blocked worker send without waiting on wall-clock time
```

## Обычный Тест

```sh
make test
```

Обычный test проверяет только то, что caller получает `DeadlineExceeded`. Он проходит, хотя background worker может остаться жить дальше.

Это слабый assertion:

```text
timeout is returned to the caller
```

Такой assertion не доказывает, что worker завершился.

## Как Поймать Через Synctest

```sh
make lint
```

`lint` target запускает bug-revealing test с `testing/synctest`. Внутри synctest bubble fake time доходит до service timeout, а потом до завершения legacy client. Поздней отправке worker уже некому принять result.

Ожидаемый вывод:

```text
--- FAIL: TestSynctestFindsBlockedSendAfterTimeout (0.00s)
    main_test.go:38: synctest detected blocked worker send: deadlock: main bubble goroutine has exited but blocked goroutines remain
FAIL
FAIL	github.com/dvordrova/find_bugs/synctest/unbuffered_send_after_timeout Xs
?   	github.com/dvordrova/find_bugs/synctest/unbuffered_send_after_timeout/internal/pricing	[no test files]
FAIL
```

## Исправление

Использовать buffered result channel, чтобы поздний result не блокировал worker:

```go
results := make(chan lookupResult, 1)
```

Для clients, которые поддерживают cancellation, еще нужно передавать cancellable context в client и делать worker select по `ctx.Done()`.

## Почему Это Важно

Без synctest timeout-race tests часто зависят от real sleeps или проверяют только return value. `testing/synctest` позволяет мгновенно пройти timeout и late worker completion.

Этот пример требует Go 1.25 или новее. Module фиксирует `go 1.26`, потому что repository CI runs Go 1.26.
