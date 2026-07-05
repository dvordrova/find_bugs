# Prompt Для Агента

Используй этот файл, когда хочешь, чтобы другой coding agent воспроизвел стиль и tooling contract репозитория `find_bugs` в новом Go-репозитории.

## Prompt Для Копирования

```text
Используй репозиторий `find_bugs` как reference style.

Собери каталог небольших запускаемых Go-примеров с багами и инструментами, которые их ловят. Результат должен быть полезен локально и в CI, а не быть просто набором snippets.

Контракт репозитория:

- Каждый пример живет в `<tool_or_problem>/<category>/`.
- Каждый пример самодостаточный и обычно содержит:
  - `go.mod`
  - `Makefile`
  - `main.go`
  - `README.md`
  - `README.ru.md`
- Код примера должен быть похож на production-код, но оставаться достаточно маленьким, чтобы быстро понять идею.
- Не комментируй каждую строку. Объяснение должно быть в README.
- Используй настоящее поведение инструментов, а не выдуманный output.
- Фиксируй версии инструментов через Go tool dependencies или явные module versions.
- Не используй `go get -u` для обычной настройки.

Контракт Makefile для каждого примера:

- Сначала идут user-facing targets:
  - `make run`, если важно показать поведение программы
  - `make lint` или ближайшая команда detector
  - `make test`
- Потом идут maintainer/CI targets:
  - `make ci-test`
  - `make tool-update`
  - targets для generated helper binaries
  - `make clean`
- `make lint` может падать, если пример специально демонстрирует баг.
- `make ci-test` должен быть стабильным в CI: он перегенерирует committed snapshots вроде `lint.logs`, `race.logs` или `lint-fixed.logs`.
- Если module downloads могут попасть в snapshot logs, запускай `go mod download` до команды, которая пишет snapshot.

Контракт root repository:

- `make test-update` проходит по всем example directories и запускает их `ci-test`.
- `make test` запускает `make test-update`, потом `git diff --exit-code`.
- Default target для `make` должен быть `test`.
- CI должен запускать одну команду: `make test`.
- Если примеры принимают custom golangci-lint config, root `make test LINT_CONFIG=/absolute/path/to/.golangci.yaml` должен пробрасывать его в релевантные примеры.

Контракт документации:

- Root `README.md` объясняет purpose, quick start, tool contract, список примеров, tools и зачем нужны snapshots.
- Root `README.ru.md` зеркалит важные user-facing части на русском.
- `BUGS.md` содержит запланированные и реализованные bug patterns.
- `CONTRIBUTING.md` объясняет, как добавить один пример.
- ADR используются только для repository-level decisions, не для обычных примеров.
- README каждого примера объясняет:
  - какой production bug показан;
  - как запустить пример локально;
  - какой инструмент ловит баг;
  - как читать ожидаемый output;
  - один реалистичный fix или mitigation.

Acceptance check:

После изменений новый разработчик или тестировщик должен уметь запустить:

```sh
make
```

и понимать failures через committed snapshot diffs. Для одного примера он должен уметь запустить:

```sh
cd <tool_or_problem>/<category>
make run
make lint
make test
```

и понять, какой баг был пойман, без чтения всего репозитория.

Избегай:

- примеров, которые требуют глобально установленных tools;
- скрытых setup steps вне Makefiles;
- незапиненных tool upgrades;
- debug-only файлов в репозитории;
- snapshots с шумом от cold-cache downloads;
- настолько искусственного кода, что баг уже не похож на production.
```

## Как Использовать

Дай prompt выше агенту вместе со ссылкой или локальным путем к этому репозиторию.

Если target repository уже существует, сначала попроси агента изучить текущую build system и адаптировать контракт, а не заменить все целиком. Цель - сохранить ту же operational shape:

- одна команда для local и CI checks;
- маленькие запускаемые примеры;
- committed snapshots output инструментов;
- pinned tool versions;
- docs, которые учат багу, а не только команде.

## Хороший Результат

Хороший clone этого подхода должен дать читателю ощущение:

- Я могу запустить весь каталог одной командой.
- Я могу запустить один пример без чтения global setup docs.
- Я вижу, какой инструмент ловит баг.
- Я могу сравнить свой lint config с каталогом.
- Я могу ревьюить изменения поведения tools через `git diff`.

## Полезный Follow-Up Prompt

```text
Проверь этот репозиторий по `docs/agent-bootstrap.ru.md`.

Скажи, какие части уже соответствуют контракту, чего не хватает, и внеси минимальные изменения, чтобы разработчик мог запускать полную проверку одной командой и смотреть per-example snapshots output инструментов.
```
