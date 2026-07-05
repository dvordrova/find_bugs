# Копирование Значения С Lock

Этот пример показывает ledger type, который владеет mutex и map:

1. `AccountLedger` содержит `sync.Mutex`.
2. У `Balance` value receiver.
3. Вызов `Balance` копирует весь `AccountLedger`, включая mutex.
4. Скопированный mutex не защищает original object так, как автор метода ожидал.

Баг находится в [main.go](main.go): `Balance` не должен копировать type, который содержит lock.

## Запуск

```sh
make run
```

Ожидаемый результат:

```text
alice balance: 120
```

Программа может выглядеть рабочей, потому что пример маленький и single-threaded.

## Как поймать через govet в golangci-lint

```sh
make lint
```

Ожидаемый вывод:

```text
main.go:21:9: copylocks: Balance passes lock by value: github.com/dvordrova/find_bugs/govet/copylocks.AccountLedger contains sync.Mutex (govet)
func (l AccountLedger) Balance(accountID string) int {
        ^
```

Отчет лучше читать как warning про форму API:

1. `Balance passes lock by value` означает, что method receiver копирует lock.
2. `AccountLedger contains sync.Mutex` объясняет, почему копирование этого value подозрительно.
3. Caret указывает на value receiver `l AccountLedger`.

`make tool-update` - maintainer-команда для осознанного обновления pinned dependency `golangci-lint`.

## Один из вариантов исправления

Использовать pointer receiver:

```go
func (l *AccountLedger) Balance(accountID string) int {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.balances[accountID]
}
```

В целом types, которые содержат locks, не стоит копировать после первого использования.
