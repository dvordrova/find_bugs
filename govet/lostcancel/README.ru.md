# Потерянный Context Cancel

Этот пример показывает service call с timeout:

1. `LoadProfile` создает child context через `context.WithTimeout`.
2. Он передает этот context в `ProfileGateway.Lookup`.
3. Возвращенная cancel function выбрасывается.
4. Timer, связанный с timeout, может жить дольше, чем нужно.

Баг находится в [main.go](main.go): у каждого успешного `WithTimeout` должен вызываться cancel function.

## Запуск

```sh
make run
```

Ожидаемый результат:

```text
loaded profile profile-001 for Alice
```

Программа выглядит рабочей, потому что timeout достаточно длинный и call возвращается быстро.

## Как поймать через govet в golangci-lint

```sh
make lint
```

Ожидаемый вывод:

```text
main.go:26:7: lostcancel: the cancel function returned by context.WithTimeout should be called, not discarded, to avoid a context leak (govet)
	ctx, _ = context.WithTimeout(ctx, 500*time.Millisecond)
	     ^
```

Отчет лучше читать как warning про lifetime ресурса:

1. `context.WithTimeout` создает timer.
2. Возвращенная cancel function освобождает этот timer раньше deadline.
3. Если выбросить cancel, cleanup ждет timeout или cancellation parent context.

`make tool-update` - maintainer-команда для осознанного обновления pinned dependency `golangci-lint`.

## Один из вариантов исправления

Вызвать cancel через `defer` после создания context:

```go
ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
defer cancel()

return gateway.Lookup(ctx, id)
```

Вызывай `cancel`, даже если operation успешно завершилась раньше deadline.
