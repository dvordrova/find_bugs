# Snapshot Download Noise Decision

Commit: `56bb13f Preload modules before snapshot checks`

## Options

1. Commit the `go: downloading ...` lines into `*.logs`.
   - Score: 1/5
   - It would make CI green once, but snapshots would depend on module cache state and become noisy.

2. Filter `go: downloading ...` lines after every captured command.
   - Score: 4/5
   - Robust, but it adds more shell plumbing around every snapshot recipe.

3. Run `go mod download` before commands that write snapshot logs.
   - Score: 5/5
   - It keeps snapshot logs focused on detector output while making cold-cache CI match warm local runs.

4. Add a GitHub Actions cache.
   - Score: 2/5
   - It can speed up CI, but it does not fix correctness because a cold cache should still produce stable snapshots.

## Selected

Option 3.

The failed pipeline showed `go: downloading ...` lines inside `lint.logs` on a clean GitHub runner. Local runs passed because the module cache was warm. Adding `go mod download` before snapshot-producing commands preloads dependencies without writing that download noise into committed logs.

Verification included a local cold-cache run with temporary `GOMODCACHE` and `GOCACHE`; the only diff was the intended Makefile change, not any `*.logs` file.
