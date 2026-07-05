# WaitGroup Add Внутри Goroutine

Этот пример показывает batch mailer:

1. `SendAll` запускает по одной goroutine на recipient.
2. Он должен дождаться всех send operations.
3. Каждая goroutine вызывает `wg.Add(1)` уже после старта.
4. `wg.Wait()` может выполниться до того, как хоть одна goroutine увеличит counter, поэтому `SendAll` может вернуться слишком рано.

Баг находится в [main.go](main.go): `WaitGroup.Add` должен происходить до запуска goroutine, которую он отслеживает.

## Запуск

```sh
make run
```

Ожидаемый результат:

```text
sent messages: 2
```

Программа выглядит рабочей, потому что `main` ждет после вызова `SendAll`. API contract все равно сломан: сам `SendAll` надежно не дождался завершения работы.

## Как поймать через govet в golangci-lint

```sh
make lint
```

Ожидаемый вывод:

```text
main.go:20:10: waitgroup: WaitGroup.Add called from inside new goroutine (govet)
			wg.Add(1)
			      ^
```

Отчет лучше читать как warning про lifecycle ordering:

1. `WaitGroup.Add called from inside new goroutine` означает, что parent может дойти до `Wait`, пока counter еще равен zero.
2. Caret указывает на `Add`, который должен был случиться до `go func`.
3. Если `Wait` вернулся слишком рано, caller может увидеть incomplete work или закрыть resources, которые еще используют workers.

`make tool-update` - maintainer-команда для осознанного обновления pinned dependency `golangci-lint`.

## Один из вариантов исправления

Вызывать `Add` до старта goroutine:

```go
for _, recipient := range recipients {
	recipient := recipient
	wg.Add(1)
	go func() {
		defer wg.Done()
		m.deliver(recipient)
	}()
}
wg.Wait()
```

На Go versions с `WaitGroup.Go` тот же lifecycle можно выразить одним вызовом.
