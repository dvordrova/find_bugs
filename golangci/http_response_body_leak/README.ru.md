# HTTP Response Body Leak

Этот пример показывает HTTP client helper, который читает status code, но забывает закрыть response body.

Production shape часто встречается в health checks, metadata clients и status-only API calls. Даже если caller не читает body, `resp.Body` нужно закрыть, чтобы transport мог переиспользовать или освободить connection.

## Запуск

```sh
make run
```

Ожидаемый результат:

```text
status=200 err=<nil>
run make lint to see the bodyclose report
```

Программа работает на маленьком request, но repeated calls могут leak connections.

## Как Поймать Через bodyclose

```sh
make lint
```

Ожидаемый отчет:

```text
internal/catalog/client.go:15:24: response body must be closed (bodyclose)
	resp, err := client.Do(req)
	                      ^
1 issues:
* bodyclose: 1
```

Читай report как нарушение resource lifetime:

1. `client.Do` возвращает `*http.Response`.
2. Function возвращается без закрытия `resp.Body`.
3. Connection не может безопасно переиспользоваться или освободиться transport.

Правило лежит в [.golangci.yaml](.golangci.yaml). Оно включает только `bodyclose`, чтобы пример оставался сфокусированным.

`make tool-update` - maintainer command для осознанного обновления pinned `golangci-lint` dependencies.

## Одно Исправление

Закрыть body сразу после проверки request error:

```go
resp, err := client.Do(req)
if err != nil {
	return 0, fmt.Errorf("fetch catalog: %w", err)
}
defer resp.Body.Close()
```
