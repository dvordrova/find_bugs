# find_bugs

Небольшие самодостаточные примеры багов на Go и инструменты, которые могут их поймать.

Каждый пример должен легко запускаться, быстро читаться и быть похожим на обычный production-код. Баги показаны явно, но без искусственно запутанного кода.

Первый пример использует Go tool dependencies для custom linter build.

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

## Первый пример

- [nilaway/cross_package_nil](nilaway/cross_package_nil/README.md): функция репозитория возвращает `nil, nil`; вызывающий код доверяет пустой ошибке и разыменовывает nil-результат. NilAway может показать этот nil-flow через custom build `golangci-lint` до runtime panic.
- [nilaway/dependency_contract_false_positive](nilaway/dependency_contract_false_positive/README.ru.md): dependency module экспортирует pointer, который инициализируется в `init`; runtime безопасен, но NilAway репортит global pointer как nilable.

## Инструменты

- `golangci-lint`: общий драйвер для Go-линтеров. Сейчас NilAway подключается к нему как custom module plugin.
- `nilaway`: статический анализатор Uber для потенциальных nil panic, здесь используется через custom build `golangci-lint`.
- `go test -race`: runtime race detector для части багов с shared memory.
- Uber leak detector (`go.uber.org/goleak`): test-time detector для утечек goroutine.

## Проверки репозитория

- `make test-update`: перегенерирует pinned tool files и lint snapshot logs для всех примеров.
- `make test`: запускает `make test-update`, а потом падает, если tracked files изменились.

Сгенерированные `*.logs` файлы коммитятся как snapshots. Когда меняется версия Go, golangci-lint или NilAway, `make test` показывает изменение отчетов через `git diff`.

Полный список планируемых багов находится в [BUGS.md](BUGS.md).

Как добавить новый пример: [CONTRIBUTING.ru.md](CONTRIBUTING.ru.md).
