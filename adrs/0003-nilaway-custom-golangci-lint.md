# ADR 0003: NilAway Through Custom golangci-lint

## Status

Accepted

## Context

NilAway is not a built-in golangci-lint linter. The examples should still use the familiar `golangci-lint` entrypoint and be easy to run with `make lint`.

The user wants to use Go tool dependencies and a custom golangci-lint build.

## Decision

Each NilAway example has:

- `.custom-gcl.yml` for the custom golangci-lint module plugin;
- `.golangci.yaml` for normal lint;
- optionally `.golangci.fixed.yaml` to demonstrate handling a confirmed false positive.

The custom plugin version is pinned:

```yaml
version: v2.12.2
plugins:
  - module: "go.uber.org/nilaway"
    import: "go.uber.org/nilaway/cmd/gclplugin"
    version: v0.0.0-20260702211033-e66cfc93566b
```

Use explicit `include-pkgs`:

- include the app package for simple examples;
- include local SDK/dependency packages only when NilAway should reason across that boundary;
- do not put every SDK in `include-pkgs` by default.

## Consequences

`make lint` builds `custom-gcl` if needed and runs NilAway through golangci-lint.

Leaving `include-pkgs` empty may be useful for exploration, but it is usually too broad/noisy for CI and can analyze packages outside the example's ownership boundary.
