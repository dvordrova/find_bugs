# ADR 0005: GitHub Actions Container Git Safety

## Status

Accepted

## Context

The CI job runs inside a container:

```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    container:
      image: golang:1.26.4
```

Snapshot CI ends with `git diff --exit-code`, so git must work inside the container.

During debugging, checkout logs showed that `.git` existed, but git failed with:

```text
fatal: detected dubious ownership in repository
```

This can happen with `actions/checkout`, containers, and self-managed runners such as actions-runner-controller. Checkout may configure safe.directory in a context that later container steps do not see.

## Decision

Keep an explicit safe-directory step after checkout:

```yaml
- name: Trust checkout directory
  run: git config --global --add safe.directory "$GITHUB_WORKSPACE"
```

Run repository checks from `$GITHUB_WORKSPACE`:

```yaml
- name: Run repository checks
  run: |
    cd "$GITHUB_WORKSPACE"
    make test
```

## Consequences

`git diff --exit-code` works inside the container despite ownership mismatch.

Debug-only steps such as printing `pwd`, `ls -al`, or `git status` were removed after the issue was understood.
