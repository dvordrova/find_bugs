# find_bugs

Небольшие самодостаточные примеры багов на Go и инструменты, которые могут их поймать.

Каждый пример должен легко запускаться, быстро читаться и быть похожим на обычный production-код. Баги показаны явно, но без искусственно запутанного кода.

## Зачем

Это executable documentation для Go bug patterns и поведения инструментов.

Цель не в benchmark инструментов и не в коллекции хитрых сломанных snippets. Цель - сохранить маленькие реалистичные примеры с точными отчетами инструментов: true positives, false positives и конфигурация, которая нужна, чтобы с ними ответственно работать.

## Быстрый старт

Запустить проверку всего репозитория:

```sh
make test
```

Запустить один пример:

```sh
cd nilaway/cross_package_nil
make run
make lint
```

Некоторые `make lint` targets ожидаемо падают, потому что показывают баг. CI использует `ci-test` targets и committed snapshot logs, чтобы проверять стабильность ожидаемых отчетов.

## Структура

```text
README.md
README.ru.md
BUGS.md
<problem_name>/<category>/
  go.mod
  Makefile
  main.go
  README.md
  README.ru.md
```

## Примеры

- [nilaway/cross_package_nil](nilaway/cross_package_nil/README.md): функция репозитория возвращает `nil, nil`; вызывающий код доверяет пустой ошибке и разыменовывает nil-результат. NilAway может показать этот nil-flow через custom build `golangci-lint` до runtime panic.
- [nilaway/dependency_contract_false_positive](nilaway/dependency_contract_false_positive/README.ru.md): dependency module экспортирует pointer, который инициализируется в `init`; runtime безопасен, но NilAway репортит global pointer как nilable.
- [goleak/channel_timeout_leak](goleak/channel_timeout_leak/README.ru.md): request истекает по timeout, а background worker позже отправляет результат в unbuffered channel. Обычный test проходит, но `go.uber.org/goleak` репортит leaked goroutine.
- [goleak/context_not_cancelled](goleak/context_not_cancelled/README.ru.md): background cache warmer принимает context, но не использует его внутри worker. Обычный test проходит, но `go.uber.org/goleak` репортит goroutine, оставшуюся в `select`.
- [race/shared_map](race/shared_map/README.ru.md): metrics collector хранит mutable counters в map и читает их, пока другая goroutine пишет. `go test -race` репортит конфликтующие accesses.
- [race/config_pointer](race/config_pointer/README.ru.md): config cache обновляет shared `*Config`, пока request handlers читают его. `go test -race` репортит unsynchronized pointer access.
- [race/shutdown_flag](race/shutdown_flag/README.ru.md): worker читает обычный shutdown boolean, пока другая goroutine пишет его. `go test -race` репортит unsynchronized flag access.

## Инструменты

- `golangci-lint`: общий драйвер для Go-линтеров. Сейчас NilAway подключается к нему как custom module plugin.
- `nilaway`: статический анализатор Uber для потенциальных nil panic, здесь используется через custom build `golangci-lint`.
- `go test -race`: runtime race detector для части багов с shared memory.
- Uber leak detector (`go.uber.org/goleak`): test-time detector для утечек goroutine.

## Проверки репозитория

- `make test-update`: перегенерирует pinned tool files и lint snapshot logs для всех примеров.
- `make test`: запускает `make test-update`, а потом падает, если tracked files изменились.

## Зачем snapshots

Сгенерированные `*.logs` файлы коммитятся как snapshots. Когда меняется версия Go, golangci-lint или NilAway, `make test` показывает изменение отчетов через `git diff`.

Так репозиторий честно фиксирует поведение инструментов. Если diagnostic изменился, diff показывает, что именно поменялось, и заставляет человека решить, ожидаемое это изменение или нет.

Полный список планируемых багов находится в [BUGS.md](BUGS.md).

Как добавить новый пример: [CONTRIBUTING.ru.md](CONTRIBUTING.ru.md).
