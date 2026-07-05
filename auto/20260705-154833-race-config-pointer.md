# Race Config Pointer Example

Commit: `6e2ec29 Add race config pointer example`

## Context

The planned example was `race/config_pointer`: config refresh replaces a shared pointer while request handlers read it.

## Options

1. Build a fresh `Config` inside the refresh goroutine on every iteration.
   Score: 6/10. Production-like, but the race detector also reports publication races on the new config object, which distracts from the pointer field.

2. Prebuild immutable configs, then race only on the shared `*Config` field.
   Score: 9/10. Keeps the report focused on the intended bug: unsynchronized pointer replacement.

3. Use `atomic.Pointer` in the example and intentionally misuse the object behind it.
   Score: 4/10. Interesting, but it teaches a more advanced failure mode than the planned one.

4. Use a package-level global config pointer.
   Score: 7/10. Common in real code, but less production-shaped than a small cache type with methods.

5. Reuse the same race snapshot normalization from `race/shared_map`.
   Score: 9/10. Keeps race examples consistent and stable across machines.

## Chosen

Options 2 and 5.

The test prebuilds immutable configs, then one goroutine calls `Refresh` while another calls `APIHost`.

## Why

This produces a single clean race report on `ConfigCache.current`. The fix story is also direct: use `atomic.Pointer[Config]` for immutable snapshots or `sync.RWMutex` for grouped state.

## Verification

- `make test` in `race/config_pointer`
- `make lint` in `race/config_pointer` failed with the expected race report.
- `make ci-test` in `race/config_pointer`
- `make test-update`
- `git diff --check`
