# Cross-Package Nil Result

This example models a common repository/service bug:

1. A lookup function returns `(*Profile, error)`.
2. The caller checks `err`.
3. The repository returns `nil, nil` when the record is missing.
4. The caller dereferences the profile and panics.

The bug is in [internal/profile/repository.go](internal/profile/repository.go): `FindByEmail` should not return `nil, nil` for a missing user.

## Run

```sh
make run
```

Expected result: the program panics because `p` is nil and `main.go` reads `p.Email`.

## Catch With golangci-lint Custom Build

NilAway is not a built-in golangci-lint linter. It is run through golangci-lint's module plugin system.

The `golangci-lint` tool is tracked in `go.mod` with Go's tool dependency mechanism.

Build and run the custom linter:

```sh
make lint
```

Expected report:

```text
main.go:18:46: Potential nil panic detected. Observed nil flow from source to dereference point:
	- profile/repository.go:28:9: literal `nil` returned from `FindByEmail()` in position 0
	- cross_package_nil/main.go:18:46: result 0 of `FindByEmail()` accessed field `Email` via the assignment(s):
		- `repo.FindByEmail(...)` to `p` at cross_package_nil/main.go:13:2 (nilaway)
	fmt.Printf("sending welcome email to %s\n", p.Email)
	                                            ^
```

Read the NilAway report from the nil source to the panic point. The top-level bullets are the main flow; the indented `via the assignment(s)` line explains how the value was stored in a local variable:

1. `profile/repository.go:28:9` is where the bad value is created: `FindByEmail` returns literal `nil` as result 0.
2. `main.go:18:46` is where result 0 of `FindByEmail` is used as `p.Email`.
3. The nested `main.go:13:2` line is not another panic point. It points to the assignment `p, err := repo.FindByEmail(...)`, where the nil result was stored in `p`.

The caret (`^`) points to the exact dereference. The first line of the report is the final place NilAway reports the problem, but the bullets explain how the nil got there.

The configuration lives in:

- [.custom-gcl.yml](.custom-gcl.yml)
- [.golangci.yaml](.golangci.yaml)

The generated `custom-gcl` binary is ignored by git. `.custom-gcl.yml` pins the NilAway plugin version so the reported flow does not silently change when a new NilAway build is released.

`include-pkgs` only contains this example module. That keeps CI focused on the code this example owns, instead of asking NilAway to analyze every loaded package.

`make tool-update` is a maintainer command for intentionally updating the tool dependency in `go.mod`; it is not needed to run the example.

NilAway should report a potential nil panic with a flow from the `nil` return in `FindByEmail` to the dereference in `main.go`.

## One Fix

Return a real not-found error and make the caller handle it:

```go
var ErrNotFound = errors.New("profile not found")

func (r *Repository) FindByEmail(email string) (*Profile, error) {
	if p, ok := r.byEmail[email]; ok {
		return p, nil
	}
	return nil, ErrNotFound
}
```

Another valid design is to return a non-pointer value plus a boolean, like `Lookup(email) (Profile, bool)`.
