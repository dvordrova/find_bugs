# Scannererr Vettool Example Decision

Commit: `e85a932 Add scannererr vettool example`

## Options

1. Only keep the scanner case in `BUGS.md`.
   - Score: 2/5
   - It documents the failure mode, but readers cannot run a detector locally.

2. Treat it as an IDE-only `gopls` note.
   - Score: 3/5
   - Useful for developers, but it does not satisfy the repository goal of a CI-friendly command.

3. Add a small `go vet -vettool` wrapper around `scannererr` from `golang.org/x/tools`.
   - Score: 5/5
   - It is executable, CI-friendly, and uses the upstream analyzer instead of copying analyzer code into the repo.

4. Write a custom scanner analyzer in this repository.
   - Score: 1/5
   - More maintenance burden, and worse than using the upstream analyzer that is already published.

## Selected

Option 3.

The example builds `scannererr-vettool` from a tiny `singlechecker` wrapper and runs it through `go vet -vettool=...`. This matches the current upstream status: `scannererr` exists in `golang.org/x/tools`, while standard `go vet` does not include it yet in Go 1.26. The README links to Go issue #17747 and notes the Go 1.28 milestone, so future cleanup is obvious.
