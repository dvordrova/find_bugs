# Scanner Err Через Vettool

Этот пример показывает line-oriented import customer IDs:

1. `ImportCustomerIDs` создает `bufio.Scanner`.
2. Он задает маленький maximum token size для input validation.
3. Он крутит loop через `scanner.Scan()`.
4. Он возвращает собранные IDs без проверки `scanner.Err()`.

Баг находится в [main.go](main.go): `Scan` возвращает `false` и на EOF, и после ошибки. Без `scanner.Err()` слишком длинная строка может выглядеть как чистый пустой import.

## Запуск

```sh
make run
```

Ожидаемый результат:

```text
imported 0 customer ids
```

Этот вывод и есть баг. Во входе есть один customer ID, но он длиннее настроенного scanner limit, поэтому scanner останавливается с ошибкой, которую код игнорирует.

## Как поймать через scannererr vettool

```sh
make lint
```

Ожидаемый вывод:

```text
main.go:13:13: bufio.Scanner "scanner" is used in Scan loop at line 17 without final check of scanner.Err()
```

Отчет лучше читать как warning про control flow:

1. `bufio.NewScanner` создает `scanner`.
2. `scanner.Scan()` используется как condition в loop.
3. После loop нет вызова `scanner.Err()`, который отличает EOF от scanner failure.

Analyzer `scannererr` доступен в `golang.org/x/tools/go/analysis/passes/scannererr`. На момент добавления примера [Go issue #17747](https://github.com/golang/go/issues/17747) был accepted с milestone Go 1.28 для будущего добавления в `cmd/vet`, поэтому здесь используется маленькая local wrapper через `go vet -vettool`.

`make tool-update` - maintainer-команда для осознанного обновления pinned analyzer dependency `golang.org/x/tools`.

## Один из вариантов исправления

Проверить `scanner.Err()` после loop:

```go
for scanner.Scan() {
	ids = append(ids, scanner.Text())
}
if err := scanner.Err(); err != nil {
	return nil, err
}

return ids, nil
```

Используй `_ = scanner.Err()` только когда игнорирование scanner errors - явное и проверенное решение.
