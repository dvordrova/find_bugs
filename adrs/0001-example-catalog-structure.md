# ADR 0001: Example Catalog Structure

## Status

Accepted

## Context

The repository is intended to collect small Go bug examples and show which tools can catch them. The examples should be easy to run, easy to read, and useful as teaching/debugging material.

The user wants examples that are self-contained and not overly artificial. Small local dependencies are acceptable when they make the bug more production-like.

## Decision

Use one directory per problem/category, with each runnable example owning its own Go module and documentation:

```text
<problem_name>/<maybe_nested_category>/
  go.mod
  Makefile
  main.go
  README.md
  README.ru.md
```

Each example should include:

- a user-facing `run` target when useful;
- a user-facing `lint` or tool target that demonstrates the problem;
- English and Russian READMEs explaining the bug and how to read tool output;
- committed snapshot logs when tool output is part of CI.

## Consequences

Examples are slightly more verbose because each one has its own module and docs, but they remain isolated and easy to execute independently.

Local modules may be nested under an example when a dependency boundary is part of the scenario.
