# DDD Repository Ruleguard Example Decision

Commit: `18c3117 Add DDD repository boundary ruleguard example`

## Options

1. Put the example under `govet` or `golangci`.
   - Score: 2/5
   - The implementation uses a linter, but the thing being taught is a team architecture rule rather than a built-in Go analyzer.

2. Create a top-level `teamrules/ddd_repository_boundary` example.
   - Score: 5/5
   - It separates project conventions from language/runtime bugs and leaves room for rules like `force_sqlc` and transaction boundaries.

3. Use `depguard` instead of `ruleguard`.
   - Score: 3/5
   - Good for broad import boundaries, but less precise for call-level rules such as allowing `database/sql` inside repositories.

4. Ban all `database/sql` imports outside repositories.
   - Score: 3/5
   - Simpler, but it cannot express narrow call-level allowances and would be less useful for examples involving shared types.

## Selected

Option 2 with `ruleguard`.

The example uses a type-aware rule: `*sql.DB` and `*sql.Tx` calls are reported outside packages whose import path ends in `/repository`. This demonstrates how a team can enforce DDD boundaries locally and in CI with one `make lint`/`make test` workflow, while still allowing repository code to own SQL access.
