# Context Первым Аргументом

Этот пример показывает team API rule:

1. I/O-facing functions должны принимать `context.Context`.
2. Если function принимает context, он должен быть первым argument.
3. Consistent context position упрощает чтение cancellation/deadlines и снижает шанс забыть context.

Нарушение находится в [internal/invoices/service/invoice.go](internal/invoices/service/invoice.go): `RebuildInvoice` принимает сначала `invoiceID`, а `context.Context` вторым argument.

## Запуск

```sh
make run
```

Ожидаемый результат:

```text
invoice rebuild requested
run make lint to see the revive report
```

Программа запускается, но API shape идет против convention из standard library и service code.

## Как Поймать Через revive

```sh
make lint
```

Ожидаемый отчет:

```text
internal/invoices/service/invoice.go:11:58: context-as-argument: context.Context should be the first parameter of a function (revive)
func (s InvoiceService) RebuildInvoice(invoiceID string, ctx context.Context) error {
                                                         ^
1 issues:
* revive: 1
```

Читай report как нарушение API consistency:

1. Function принимает `context.Context`.
2. Перед ним стоит другой argument.
3. Team rule ожидает, что `ctx context.Context` будет первым.

Правило лежит в [.golangci.yaml](.golangci.yaml). Оно включает только revive rule `context-as-argument`, чтобы пример оставался сфокусированным.

`make tool-update` - maintainer command для осознанного обновления pinned `golangci-lint` dependencies.

## Одно Исправление

Перенести context на первое место:

```go
func (s InvoiceService) RebuildInvoice(ctx context.Context, invoiceID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return nil
}
```
