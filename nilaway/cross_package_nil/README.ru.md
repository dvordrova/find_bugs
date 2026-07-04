# Cross-Package Nil Result

Этот пример показывает частый баг в коде с repository/service:

1. Функция поиска возвращает `(*Profile, error)`.
2. Вызывающий код проверяет `err`.
3. Репозиторий возвращает `nil, nil`, если запись не найдена.
4. Вызывающий код разыменовывает profile и получает panic.

Баг находится в [internal/profile/repository.go](internal/profile/repository.go): `FindByEmail` не должен возвращать `nil, nil` для отсутствующего пользователя.

## Запуск

```sh
make run
```

Ожидаемый результат: программа падает, потому что `p` равен nil, а `main.go` читает `p.Email`.

## Как поймать через custom build golangci-lint

NilAway не является встроенным линтером golangci-lint. Он запускается через module plugin system.

`golangci-lint` хранится в `go.mod` через Go tool dependencies.

Собрать и запустить custom linter:

```sh
make lint
```

Ожидаемый вывод:

```text
main.go:18:46: Potential nil panic detected. Observed nil flow from source to dereference point:
	- profile/repository.go:28:9: literal `nil` returned from `FindByEmail()` in position 0
	- cross_package_nil/main.go:18:46: result 0 of `FindByEmail()` accessed field `Email` via the assignment(s):
		- `repo.FindByEmail(...)` to `p` at cross_package_nil/main.go:13:2 (nilaway)
	fmt.Printf("sending welcome email to %s\n", p.Email)
	                                            ^
```

Отчет NilAway лучше читать как путь от источника nil к месту panic. Основной flow - это верхнеуровневые bullets, а вложенная строка `via the assignment(s)` объясняет, через какое присваивание значение попало в локальную переменную:

1. `profile/repository.go:28:9` - место, где появляется плохое значение: `FindByEmail` возвращает literal `nil` как result 0.
2. `main.go:18:46` - место, где result 0 от `FindByEmail` используется как `p.Email`.
3. Вложенная строка `main.go:13:2` - это не еще одно место panic. Она указывает на присваивание `p, err := repo.FindByEmail(...)`, где nil-result был сохранен в `p`.

Каретка (`^`) показывает точное место dereference. Первая строка отчета показывает финальную точку проблемы, но bullets объясняют, как nil туда попал.

Конфигурация находится в:

- [.custom-gcl.yml](.custom-gcl.yml)
- [.golangci.yaml](.golangci.yaml)

Сгенерированный бинарник `custom-gcl` игнорируется git. `.custom-gcl.yml` фиксирует версию NilAway plugin, чтобы отчет не поменялся молча после выхода новой версии NilAway.

`include-pkgs` содержит только модуль этого примера. Так CI проверяет код, которым владеет пример, а не пытается анализировать каждый загруженный package.

`make tool-update` - maintainer-команда для осознанного обновления tool dependency в `go.mod`; для запуска примера она не нужна.

NilAway должен показать potential nil panic: поток идет от `nil` в `FindByEmail` к dereference в `main.go`.

## Один из вариантов исправления

Возвращать явную ошибку not found и обрабатывать ее в вызывающем коде:

```go
var ErrNotFound = errors.New("profile not found")

func (r *Repository) FindByEmail(email string) (*Profile, error) {
	if p, ok := r.byEmail[email]; ok {
		return p, nil
	}
	return nil, ErrNotFound
}
```

Другой нормальный вариант дизайна: вернуть значение без pointer плюс boolean, например `Lookup(email) (Profile, bool)`.
