# Data Race На Config Pointer

Этот пример показывает config cache, которым пользуются request handlers:

1. `ConfigCache.Refresh` заменяет текущий config pointer.
2. `ConfigCache.APIHost` читает этот pointer на request path.
3. Refresh и reads могут происходить одновременно.
4. Pointer field не защищен mutex или atomic operation, поэтому read и write race.

Баг находится в [main.go](main.go): `current` - shared mutable state.

## Запуск

```sh
make run
```

Ожидаемый результат:

```text
api host: api.internal
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
  github.com/dvordrova/find_bugs/race/config_pointer.(*ConfigCache).Refresh()
      github.com/dvordrova/find_bugs/race/config_pointer/main.go:18 +0xADDR

Previous read at 0xADDR by goroutine N:
  github.com/dvordrova/find_bugs/race/config_pointer.(*ConfigCache).APIHost()
      github.com/dvordrova/find_bugs/race/config_pointer/main.go:22 +0xADDR
```

Отчет лучше читать как два конфликтующих access к одному pointer field:

1. `Write at` указывает на `Refresh`, который заменяет `current`.
2. `Previous read at` указывает на `APIHost`, который читает `current`.
3. Stack создания goroutine показывает, где refresh worker был запущен в test.

`make test` может пройти, потому что race detector там не включен.

Committed snapshot [race.logs](race.logs) нормализует addresses, goroutine IDs и package duration, чтобы diffs были стабильными на разных машинах.

## Один из вариантов исправления

Использовать `atomic.Pointer[Config]` для immutable config snapshots:

```go
type ConfigCache struct {
	current atomic.Pointer[Config]
}
```

Храни только полностью построенные immutable configs. `sync.RWMutex` тоже подходит, когда refresh должен обновить несколько связанных fields вместе.
