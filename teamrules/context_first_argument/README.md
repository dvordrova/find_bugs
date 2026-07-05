# Context First Argument

This example shows a team API rule:

1. I/O-facing functions should accept `context.Context`.
2. When a function accepts a context, it should be the first argument.
3. Consistent context position makes cancellation and deadlines easy to scan and hard to forget.

The violation is in [internal/invoices/service/invoice.go](internal/invoices/service/invoice.go): `RebuildInvoice` accepts `invoiceID` first and `context.Context` second.

## Run

```sh
make run
```

Expected result:

```text
invoice rebuild requested
run make lint to see the revive report
```

The program runs, but the API shape fights the convention used by Go standard-library and service code.

## Catch With revive

```sh
make lint
```

Expected report:

```text
internal/invoices/service/invoice.go:11:58: context-as-argument: context.Context should be the first parameter of a function (revive)
func (s InvoiceService) RebuildInvoice(invoiceID string, ctx context.Context) error {
                                                         ^
1 issues:
* revive: 1
```

Read the report as an API consistency violation:

1. The function accepts `context.Context`.
2. Another argument appears before it.
3. The team rule expects `ctx context.Context` to come first.

The rule lives in [.golangci.yaml](.golangci.yaml). It enables only revive's `context-as-argument` rule so the example stays focused.

`make tool-update` is a maintainer command for intentionally updating pinned `golangci-lint` dependencies.

## One Fix

Move context to the first argument:

```go
func (s InvoiceService) RebuildInvoice(ctx context.Context, invoiceID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return nil
}
```
