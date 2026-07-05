# find_bugs

Небольшие самодостаточные примеры багов на Go и инструменты, которые могут их поймать.

Каждый пример должен легко запускаться, быстро читаться и быть похожим на обычный production-код. Баги показаны явно, но без искусственно запутанного кода.

## Зачем

Это executable documentation для Go bug patterns и поведения инструментов.

Цель не в benchmark инструментов и не в коллекции хитрых сломанных snippets. Цель - сохранить маленькие реалистичные примеры с точными отчетами инструментов: true positives, false positives и конфигурация, которая нужна, чтобы с ними ответственно работать.

## Базовая статья

Concurrency-часть каталога опирается на ["Understanding Real-World Concurrency Bugs in Go"](https://songlh.github.io/paper/go-study.pdf) Tu, Liu, Song и Zhang. Авторы изучили 171 баг из production Go-проектов и используют полезную таксономию: поведение `blocking` vs `non-blocking`, пересеченное с причинами `shared memory` vs `message passing`.

Используй эту статью как карту, а примеры в репозитории - как запускаемые checkpoints.

## Быстрый старт

Запустить проверку всего репозитория:

```sh
make
```

То же самое явным target:

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

## Tooling Contract

Каждый пример устроен так, чтобы читатель мог превратить урок в локальную или CI-проверку:

- `make run` показывает поведение программы.
- `make lint` запускает detector, который должен поймать баг.
- `make test` запускает обычный test path для этого примера.
- `make ci-test` - repository check, который перегенерирует committed logs и проверяет, что ожидаемый detector signal все еще есть.

Версии инструментов живут внутри module примера, обычно через Go tool dependencies и `go tool`. Сгенерированные helper binaries, например custom build `golangci-lint` для NilAway или scannererr vettool, игнорируются git.

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
- [govet/copylocks](govet/copylocks/README.ru.md): method копирует struct, который содержит `sync.Mutex`. `govet` через `golangci-lint` репортит copied lock value.
- [govet/nocopy_marker](govet/nocopy_marker/README.ru.md): type явно включает copy detection через private `noCopy` marker. `govet` через `golangci-lint` репортит accidental value copies.
- [govet/lostcancel](govet/lostcancel/README.ru.md): timeout context создается, но cancel function выбрасывается. `govet` через `golangci-lint` репортит context leak.
- [govet/waitgroup_add_inside_goroutine](govet/waitgroup_add_inside_goroutine/README.ru.md): `WaitGroup.Add` вызывается внутри goroutine, которую должен отслеживать. `govet` через `golangci-lint` репортит lifecycle ordering bug.
- [govet/scannererr_vettool](govet/scannererr_vettool/README.ru.md): line importer использует `bufio.Scanner` с маленьким token limit и забывает `scanner.Err`. Local wrapper через `go vet -vettool` запускает analyzer `scannererr` из `golang.org/x/tools`.
- [golangci/sql_rows_not_closed](golangci/sql_rows_not_closed/README.ru.md): repository method сканирует database rows и проверяет iteration errors, но забывает `rows.Close`. `sqlclosecheck` через `golangci-lint` репортит resource leak.
- [teamrules/ddd_repository_boundary](teamrules/ddd_repository_boundary/README.ru.md): service code напрямую вызывает `*sql.DB`. Type-aware правило `ruleguard` держит database calls внутри repository packages.
- [teamrules/no_wall_clock_in_domain](teamrules/no_wall_clock_in_domain/README.ru.md): domain code напрямую вызывает `time.Now`. Узкое правило `ruleguard` держит wall-clock reads в adapters или composition roots.
- [synctest/context_afterfunc_negative_assertion](synctest/context_afterfunc_negative_assertion/README.ru.md): cancellation hook пишет audit record до cancel context. `testing/synctest` делает assertion "еще ничего не произошло" детерминированным.

## Инструменты

- `golangci-lint`: общий драйвер для Go-линтеров. Сейчас NilAway подключается к нему как custom module plugin.
- `sqlclosecheck` и `rowserrcheck`: точечные golangci-lint линтеры для lifetime SQL rows и iteration errors.
- `scannererr`: Go analysis pass из `golang.org/x/tools`, здесь запускается через маленький local binary для `go vet -vettool`, пока он не доступен через стандартный `go vet`.
- `ruleguard`: custom team rules для architecture boundaries и project-specific conventions.
- `testing/synctest`: стандартная библиотека для детерминированных тестов concurrent code, fake time и negative assertions без wall-clock sleeps.
- `nilaway`: статический анализатор Uber для потенциальных nil panic, здесь используется через custom build `golangci-lint`.
- `go test -race`: runtime race detector для части багов с shared memory.
- Uber leak detector (`go.uber.org/goleak`): test-time detector для утечек goroutine.

## Проверки репозитория

- `make test-update`: перегенерирует pinned tool files и lint snapshot logs для всех примеров.
- `make` или `make test`: запускает `make test-update`, а потом падает, если tracked files изменились.

Попробовать другой golangci-lint config на каталоге:

```sh
make test LINT_CONFIG=/Users/me/project/.golangci.yaml
```

`config=/Users/me/project/.golangci.yaml` работает как короткий alias. В этом режиме golangci-lint примеры пишут в свои `lint.logs` то, что выдал custom config; примеры на `go test -race` и `goleak` продолжают использовать свои инструменты. Финальный `git diff` показывает, ловит ли custom config ожидаемые проблемы.

## Зачем snapshots

Сгенерированные `*.logs` файлы коммитятся как snapshots. Когда меняется версия Go, golangci-lint или NilAway, `make test` показывает изменение отчетов через `git diff`.

Так репозиторий честно фиксирует поведение инструментов. Если diagnostic изменился, diff показывает, что именно поменялось, и заставляет человека решить, ожидаемое это изменение или нет.

Полный список планируемых багов находится в [BUGS.md](BUGS.md).

Как добавить новый пример: [CONTRIBUTING.ru.md](CONTRIBUTING.ru.md).

Рабочий backlog для следующих примеров и tooling лежит в [docs/backlog.md](docs/backlog.md).

Чтобы воспроизвести стиль этого репозитория в другом проекте, используй prompt для агента: [docs/agent-bootstrap.ru.md](docs/agent-bootstrap.ru.md).
