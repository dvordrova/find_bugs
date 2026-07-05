# Scanner Err Vettool

This example models a line-oriented customer import:

1. `ImportCustomerIDs` creates a `bufio.Scanner`.
2. It sets a small maximum token size for input validation.
3. It loops with `scanner.Scan()`.
4. It returns the collected IDs without checking `scanner.Err()`.

The bug is in [main.go](main.go): `Scan` returns `false` both at EOF and after an error. Without `scanner.Err()`, a long line can look like a clean empty import.

## Run

```sh
make run
```

Expected result:

```text
imported 0 customer ids
```

That output is the bug. The input contains one customer ID, but it is longer than the configured scanner limit, so the scanner stops with an error that the code ignores.

## Catch With scannererr Vettool

```sh
make lint
```

Expected report:

```text
main.go:13:13: bufio.Scanner "scanner" is used in Scan loop at line 17 without final check of scanner.Err()
```

Read the report as a control-flow warning:

1. `bufio.NewScanner` creates `scanner`.
2. `scanner.Scan()` is used as the loop condition.
3. No later call to `scanner.Err()` distinguishes EOF from scanner failure.

The `scannererr` analyzer is available in `golang.org/x/tools/go/analysis/passes/scannererr`. At the time this example was added, [Go issue #17747](https://github.com/golang/go/issues/17747) was accepted with a Go 1.28 milestone for future `cmd/vet` integration, so this example uses a tiny local `go vet -vettool` wrapper.

`make tool-update` is a maintainer command for intentionally updating the pinned `golang.org/x/tools` analyzer dependency.

## One Fix

Check `scanner.Err()` after the loop:

```go
for scanner.Scan() {
	ids = append(ids, scanner.Text())
}
if err := scanner.Err(); err != nil {
	return nil, err
}

return ids, nil
```

Use `_ = scanner.Err()` only when ignoring scanner errors is an explicit, reviewed decision.
