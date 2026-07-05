# HTTP Response Body Leak

This example shows an HTTP client helper that reads the status code but forgets to close the response body.

The production shape is common in health checks, metadata clients, and status-only API calls. Even when the caller does not need the body, `resp.Body` must be closed so the transport can reuse or release the connection.

## Run

```sh
make run
```

Expected result:

```text
status=200 err=<nil>
run make lint to see the bodyclose report
```

The program works on a tiny request, but repeated calls can leak connections.

## Catch With bodyclose

```sh
make lint
```

Expected report:

```text
internal/catalog/client.go:15:24: response body must be closed (bodyclose)
	resp, err := client.Do(req)
	                      ^
1 issues:
* bodyclose: 1
```

Read the report as a resource lifetime violation:

1. `client.Do` returns `*http.Response`.
2. The function returns without closing `resp.Body`.
3. The connection cannot be safely reused or released by the transport.

The rule lives in [.golangci.yaml](.golangci.yaml). It enables only `bodyclose` so the example stays focused.

`make tool-update` is a maintainer command for intentionally updating pinned `golangci-lint` dependencies.

## One Fix

Close the body immediately after checking the request error:

```go
resp, err := client.Do(req)
if err != nil {
	return 0, fmt.Errorf("fetch catalog: %w", err)
}
defer resp.Body.Close()
```
