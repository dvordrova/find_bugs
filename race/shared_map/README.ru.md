# Data Race На Shared Map

Этот пример показывает request metrics collector:

1. `Metrics.Record` обновляет counter, который хранится в map.
2. `Metrics.Snapshot` копирует counters для reporting.
3. Writer goroutine записывает requests, пока другая goroutine читает snapshot.
4. Map владеет shared mutable counters, и они не защищены mutex, поэтому read и write race.

Баг находится в [main.go](main.go): `counts` владеет shared mutable `Counter` values, и оба метода ходят в эти counters без synchronization.

## Запуск

```sh
make run
```

Ожидаемый результат:

```text
checkout requests: 1
```

Standalone program однопоточный, поэтому race в нем не проявляется.

## Как поймать через race detector

```sh
make lint
```

Форма ожидаемого отчета:

```text
WARNING: DATA RACE
Write at 0xADDR by goroutine N:
  github.com/dvordrova/find_bugs/race/shared_map.(*Metrics).Record()
      github.com/dvordrova/find_bugs/race/shared_map/main.go:23 +0xADDR

Previous read at 0xADDR by goroutine N:
  github.com/dvordrova/find_bugs/race/shared_map.(*Metrics).Snapshot()
      github.com/dvordrova/find_bugs/race/shared_map/main.go:29 +0xADDR
```

Race report лучше читать как два конфликтующих memory access:

1. `Write at` указывает на `Record`, который increment `Counter.Value`.
2. `Previous read at` указывает на `Snapshot`, который читает тот же `Counter.Value`.
3. Stack создания goroutine показывает, где concurrent writer был запущен в test.

`make test` запускает тот же test без `-race`; он может пройти, потому что data race не является обычным test assertion.

Committed snapshot [race.logs](race.logs) нормализует addresses, goroutine IDs и package duration, чтобы diffs были стабильными на разных машинах.

## Один из вариантов исправления

Защитить map и counters через mutex:

```go
type Metrics struct {
	mu     sync.RWMutex
	counts map[string]int
}
```

Использовать `mu.Lock` в `Record` и `mu.RLock` в `Snapshot`. Другой нормальный дизайн - отправлять все metric updates в одну owner goroutine.
