# Data Race На Shutdown Flag

Этот пример показывает worker со stop flag:

1. `Worker.Run` проверяет `stopping` в loop.
2. `Worker.Stop` пишет `stopping = true` из другой goroutine.
3. Нет mutex, atomic value или channel close, который связывает read и write.
4. Worker обычно останавливается, но доступ к flag все равно остается data race.

Баг находится в [main.go](main.go): `stopping` shared между goroutines без synchronization.

## Запуск

```sh
make run
```

Ожидаемый результат:

```text
worker stopped
```

Программа часто выглядит рабочей, поэтому такой баг легко пропустить.

## Как поймать через race detector

```sh
make lint
```

Форма ожидаемого отчета:

```text
WARNING: DATA RACE
Write at 0xADDR by goroutine N:
  github.com/dvordrova/find_bugs/race/shutdown_flag.(*Worker).Stop()
      github.com/dvordrova/find_bugs/race/shutdown_flag/main.go:26 +0xADDR

Previous read at 0xADDR by goroutine N:
  github.com/dvordrova/find_bugs/race/shutdown_flag.(*Worker).Run()
      github.com/dvordrova/find_bugs/race/shutdown_flag/main.go:20 +0xADDR
```

Отчет лучше читать как два конфликтующих access:

1. `Write at` указывает на `Stop`, который пишет shutdown flag.
2. `Previous read at` указывает на `Run`, который проверяет тот же flag.
3. Stack создания goroutine показывает, где worker был запущен.

`make test` может пройти, потому что worker eventually останавливается; `-race` проверяет именно memory synchronization.

Committed snapshot [race.logs](race.logs) нормализует addresses, goroutine IDs и package duration, чтобы diffs были стабильными на разных машинах.

## Один из вариантов исправления

Использовать context cancellation или закрытый channel вместо shared boolean:

```go
func (w *Worker) Run(ctx context.Context, pollEvery time.Duration) {
	for {
		select {
		case <-time.After(pollEvery):
		case <-ctx.Done():
			return
		}
	}
}
```

Если flag действительно подходит лучше, используй `atomic.Bool` или защищай его mutex.
