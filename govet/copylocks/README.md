# Copied Lock Value

This example models a ledger type that owns a mutex and a map:

1. `AccountLedger` contains `sync.Mutex`.
2. `Balance` has a value receiver.
3. Calling `Balance` copies the whole `AccountLedger`, including the mutex.
4. The copied mutex does not protect the original object the way the method author intended.

The bug is in [main.go](main.go): `Balance` should not copy a type that contains a lock.

## Run

```sh
make run
```

Expected result:

```text
alice balance: 120
```

The program can appear to work because the example is small and single-threaded.

## Catch With govet Through golangci-lint

```sh
make lint
```

Expected report:

```text
main.go:21:9: copylocks: Balance passes lock by value: github.com/dvordrova/find_bugs/govet/copylocks.AccountLedger contains sync.Mutex (govet)
func (l AccountLedger) Balance(accountID string) int {
        ^
```

Read the report as an API-shape warning:

1. `Balance passes lock by value` means the method receiver copies the lock.
2. `AccountLedger contains sync.Mutex` explains why copying this value is suspicious.
3. The caret points to the value receiver `l AccountLedger`.

`make tool-update` is a maintainer command for intentionally updating the pinned `golangci-lint` dependency.

## One Fix

Use a pointer receiver:

```go
func (l *AccountLedger) Balance(accountID string) int {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.balances[accountID]
}
```

In general, types that contain locks should not be copied after first use.
