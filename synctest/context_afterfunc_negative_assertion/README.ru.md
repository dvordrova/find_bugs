# Negative Assertion Для Context AfterFunc

Этот пример показывает cancellation hook, который пишет audit record слишком рано.

Production-сценарий обычный: service регистрирует callback, который должен выполниться при cancel context. Баг в том, что callback дополнительно запускается уже во время регистрации, поэтому audit record появляется до cancel.

## Запуск

```sh
make run
```

Ожидаемое поведение:

```text
before cancel: 0 audit record(s)
after cancel: 1 audit record(s)
```

Фактическое поведение отличается, потому что callback запускается до cancel.

```text
before cancel: 1 audit record(s)
after cancel: 2 audit record(s)
```

## Обычный Тест

```sh
make test
```

Обычный тест проверяет только то, что после cancel audit record когда-нибудь появился. Он проходит даже если record был записан до cancel.

Слабый assertion выглядит так:

```text
after cancel, at least one audit record exists
```

Он не доказывает, что hook дождался cancel.

## Как Поймать Через Synctest

```sh
make lint
```

Target `lint` запускает bug-revealing test через `testing/synctest`. Тест входит в synctest bubble, регистрирует callback, а потом вызывает `synctest.Wait()` до cancel context.

`synctest.Wait()` ждет, пока goroutine внутри bubble завершится или durably blocked. Поэтому negative assertion становится детерминированным:

```text
before cancel, no audit record exists
```

Ожидаемый вывод:

```text
--- FAIL: TestSynctestChecksBeforeCancel (0.00s)
    main_test.go:33: before cancel: audit records = 1, want 0
FAIL
FAIL	github.com/dvordrova/find_bugs/synctest/context_afterfunc_negative_assertion Xs
FAIL
```

## Исправление

Убрать eager goroutine и только зарегистрировать callback через `context.AfterFunc`:

```go
func NotifyWhenCanceled(ctx context.Context, sink *AuditSink, accountID string) func() bool {
	return context.AfterFunc(ctx, func() {
		sink.Record(accountID)
	})
}
```

## Почему Это Важно

Без synctest проверки вида "событие еще не произошло" часто пишут через sleep. Такие тесты либо медленные, либо flaky, либо сразу оба. `testing/synctest` позволяет дождаться, пока concurrent code успокоится, без ожидания wall-clock time.

Пример требует Go 1.25 или новее. Модуль фиксирует `go 1.26`, потому что CI репозитория запускается на Go 1.26.
