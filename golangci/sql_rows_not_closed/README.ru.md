# Незакрытый SQL Rows

Этот пример показывает repository method, который загружает неоплаченные invoices из базы:

1. `OpenInvoices` вызывает `QueryContext`.
2. Он сканирует все rows.
3. Он проверяет `rows.Err()`.
4. Он забывает закрыть `rows`.

Баг находится в [main.go](main.go): `*sql.Rows` держит database resources, пока его не закрыли.

## Запуск

```sh
make run
```

Ожидаемый результат:

```text
load open invoices
```

Реальное подключение к базе не нужно, потому что баг показывается статически.

## Как поймать через sqlclosecheck в golangci-lint

```sh
make lint
```

Ожидаемый вывод:

```text
main.go:20:32: Rows/Stmt/NamedStmt was not closed (sqlclosecheck)
	rows, err := s.db.QueryContext(ctx, `
	                              ^
```

Отчет лучше читать как warning про lifetime database resource:

1. `QueryContext` возвращает `*sql.Rows`.
2. Вызывающий код становится владельцем этого rows object.
3. Функция возвращается без `rows.Close()`.

В этом примере `rows.Err()` проверяется специально. Пропущенная проверка iteration error - другой баг, и для него точечный линтер `rowserrcheck`.

`make tool-update` - maintainer-команда для осознанного обновления pinned dependency `golangci-lint`.

## Один из вариантов исправления

Закрыть rows сразу после успешного query:

```go
rows, err := s.db.QueryContext(ctx, query)
if err != nil {
	return nil, err
}
defer rows.Close()
```

Проверку `rows.Err()` после loop нужно оставить, потому что `Close` и iteration errors закрывают разные failure modes.
