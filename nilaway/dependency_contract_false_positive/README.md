# SDK Init Global False Positive

This example shows a NilAway false positive that crosses into a dependency module.

The app depends on a small local SDK module:

```go
require github.com/acme/contractsdk v0.0.0

replace github.com/acme/contractsdk => ./contractsdk
```

The SDK exposes a package-level pointer:

```go
var DefaultPlan *Plan

func init() {
	DefaultPlan = &Plan{
		Name: "enterprise",
	}
}
```

At runtime this is safe because Go runs package `init` before `main`, so `contractsdk.DefaultPlan` is initialized before the app reads it.

## Run

```sh
make run
```

Expected result:

```text
tenant uses enterprise plan
```

## Catch With NilAway

```sh
make lint
```

Expected report:

```text
main.go:16:38: Potential nil panic detected. Observed nil flow from source to dereference point:
	- contractsdk/tenant.go:8:5: nilable value assigned into global variable `DefaultPlan`
	- dependency_contract_false_positive/main.go:12:12: global variable `DefaultPlan` passed as arg `plan` to `printPlan()` via the assignment(s):
		- `contractsdk.DefaultPlan` to `plan` at dependency_contract_false_positive/main.go:10:2
	- dependency_contract_false_positive/main.go:16:38: function parameter `plan` accessed field `Name` (nilaway)
	fmt.Printf("tenant uses %s plan\n", plan.Name)
	                                    ^
```

Read the report as a flow from the possible nil source to the dereference:

1. `contractsdk/tenant.go:8:5`: `DefaultPlan` is a package-level pointer. Its zero value is nil.
2. `main.go:10:2`: the app assigns `contractsdk.DefaultPlan` to local variable `plan`.
3. `main.go:12:12`: the app passes `plan` into `printPlan`.
4. `main.go:16:38`: `printPlan` dereferences `plan.Name`.

The false-positive part is that NilAway does not prove the SDK `init` function always assigns `DefaultPlan` before `main` uses it. Runtime does guarantee package initialization order, but static analysis is conservative around mutable global pointers.

## Handle The Confirmed False Positive

In production, do not hide this with a broad dependency exclusion. A global pointer is still a risky API shape: another SDK version, test, or future mutation could make it nil.

Both lint configs keep `include-pkgs` explicit:

```yaml
include-pkgs: github.com/dvordrova/find_bugs/nilaway/dependency_contract_false_positive,github.com/acme/contractsdk
```

That means NilAway analyzes the app and the SDK contract involved in this example. Do not put every SDK there by default. Add packages you own, or packages whose nil contracts you intentionally want NilAway to reason about. Leaving `include-pkgs` empty is useful for exploration, but it is usually too broad for CI because third-party and generated packages can add noise and slow the run.

This example uses a narrow golangci-lint exclusion in [.golangci.fixed.yaml](.golangci.fixed.yaml):

```yaml
linters:
  exclusions:
    warn-unused: true
    rules:
      # Known false positive: contractsdk.DefaultPlan is an SDK global pointer
      # initialized from contractsdk.init before main runs. Keep this narrow:
      # only NilAway, only this app file, only the DefaultPlan global-flow report.
      # In production, back this kind of suppression with an SDK contract or test.
      - linters:
          - nilaway
        path: ^main\.go$
        text: nilable value assigned into global variable `DefaultPlan`
```

This is intentionally specific:

1. `linters: [nilaway]` keeps every other linter active.
2. `path: ^main\.go$` limits the rule to this app entrypoint.
3. `text: ...DefaultPlan` matches the known SDK global-pointer false positive, not every nil dereference.
4. `warn-unused: true` makes golangci-lint warn if the exclusion stops matching after the code changes.

Run:

```sh
make lint-fixed
```

or directly:

```sh
./custom-gcl run --config .golangci.fixed.yaml
```

Expected lint result:

```text
0 issues.
```

The custom linter config pins the NilAway module version in [.custom-gcl.yml](.custom-gcl.yml), so `make lint` and `make lint-fixed` use the same analyzer version until it is intentionally updated.

`make tool-update` is a maintainer command for intentionally updating the tool dependency in `go.mod`; it is not needed to run the example.

## Production Alternatives

In real code, a better fix is often to avoid spreading a nil-capable SDK global through business logic. Check the SDK boundary once:

```go
func defaultPlan() *contractsdk.Plan {
	plan := contractsdk.DefaultPlan
	if plan == nil {
		log.Fatal("SDK default plan is not initialized")
	}

	return plan
}
```

If the SDK is under your control, an even stronger API is to avoid exported mutable pointer globals:

```go
func DefaultPlan() *Plan
```

or:

```go
func DefaultPlan() (*Plan, error)
```
