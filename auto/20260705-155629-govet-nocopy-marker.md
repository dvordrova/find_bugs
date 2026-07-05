# Govet noCopy Marker Example

Commit: `74eb983 Add govet nocopy marker example`

## Context

The planned example was `govet/nocopy_marker`: use an explicit marker so `govet copylocks` catches accidental copies of a type that should not be copied.

## Options

1. Use the standard unexported `noCopy` struct with pointer `Lock` and `Unlock` methods.
   Score: 9/10. Matches the pattern `govet copylocks` recognizes and keeps the runtime behavior zero-cost.

2. Use an interface field named `noCopy`.
   Score: 4/10. It sounds close to the idea, but `copylocks` detects the lock-like method set on a concrete type more directly.

3. Reuse `sync.Mutex` as the marker.
   Score: 5/10. Works, but then the example becomes another accidental-lock example instead of an intentional non-copyable contract.

4. Trigger the copy through assignment instead of a value receiver.
   Score: 7/10. Valid, but the value receiver makes the API mistake easier to see.

5. Explain the pattern in README instead of adding comments around every line.
   Score: 9/10. Keeps code small and lets documentation carry the teaching.

## Chosen

Options 1, 4's receiver variant, and 5.

`StreamConsumer` embeds private `noCopy`; `Topic` has a value receiver, so `govet copylocks` reports that calling the method copies the non-copyable owner.

## Why

This complements `govet/copylocks`: one example shows an accidental `sync.Mutex` copy, the other shows a deliberate no-copy contract for lifecycle-sensitive handles.

## Verification

- `make tool-update` in `govet/nocopy_marker`
- `make test` in `govet/nocopy_marker`
- `make lint` in `govet/nocopy_marker` failed with the expected `copylocks` report.
- `make ci-test` in `govet/nocopy_marker`
- `make test-update`
- `git diff --check`
