# Явный noCopy Marker

Этот пример показывает type, который не должен копироваться после construction:

1. `StreamConsumer` представляет long-lived consumer handle.
2. Type содержит private `noCopy` marker.
3. У `noCopy` есть методы `Lock` и `Unlock` на pointer type.
4. Analyzer `copylocks` внутри `govet` распознает такую форму и репортит случайные copies.

Баг находится в [main.go](main.go): у `Topic` value receiver, поэтому вызов копирует `StreamConsumer`.

## Запуск

```sh
make run
```

Ожидаемый результат:

```text
consumer topic: payments
```

Программа может выглядеть рабочей, потому что marker не имеет runtime behavior. Он нужен для static analysis.

## Как поймать через govet в golangci-lint

```sh
make lint
```

Ожидаемый вывод:

```text
main.go:19:9: copylocks: Topic passes lock by value: github.com/dvordrova/find_bugs/govet/nocopy_marker.StreamConsumer contains github.com/dvordrova/find_bugs/govet/nocopy_marker.noCopy (govet)
func (c StreamConsumer) Topic() string {
        ^
```

Отчет лучше читать как intentional non-copyable contract:

1. `Topic passes lock by value` означает, что method receiver копирует owner type.
2. `StreamConsumer contains ... noCopy` объясняет, что этот type явно включил copy detection.
3. Caret указывает на value receiver `c StreamConsumer`.

`noCopy` не делает магии в runtime. Это работает потому, что `govet copylocks` считает type lock-like, когда `*T` implements `Lock` и `Unlock`, а сам `T` - нет.

`make tool-update` - maintainer-команда для осознанного обновления pinned dependency `golangci-lint`.

## Один из вариантов исправления

Использовать pointer receivers и передавать `*StreamConsumer` через API:

```go
func (c *StreamConsumer) Topic() string {
	return c.topic
}
```

Этот pattern полезен для handles, которые владеют goroutines, file descriptors, mutex-protected state или lifecycle-sensitive resources.
