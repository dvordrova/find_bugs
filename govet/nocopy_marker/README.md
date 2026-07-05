# Explicit noCopy Marker

This example models a type that should not be copied after construction:

1. `StreamConsumer` represents a long-lived consumer handle.
2. The type embeds a private `noCopy` marker.
3. `noCopy` has `Lock` and `Unlock` methods on its pointer type.
4. `govet`'s `copylocks` analyzer recognizes that shape and reports accidental copies.

The bug is in [main.go](main.go): `Topic` has a value receiver, so calling it copies `StreamConsumer`.

## Run

```sh
make run
```

Expected result:

```text
consumer topic: payments
```

The program can appear to work because the marker has no runtime behavior. It exists for static analysis.

## Catch With govet Through golangci-lint

```sh
make lint
```

Expected report:

```text
main.go:19:9: copylocks: Topic passes lock by value: github.com/dvordrova/find_bugs/govet/nocopy_marker.StreamConsumer contains github.com/dvordrova/find_bugs/govet/nocopy_marker.noCopy (govet)
func (c StreamConsumer) Topic() string {
        ^
```

Read the report as an intentional non-copyable contract:

1. `Topic passes lock by value` means the method receiver copies the owner type.
2. `StreamConsumer contains ... noCopy` explains that this type opted into copy detection.
3. The caret points to the value receiver `c StreamConsumer`.

`noCopy` is not magic at runtime. It works because `govet copylocks` treats a type as lock-like when `*T` implements `Lock` and `Unlock`, while `T` itself does not.

`make tool-update` is a maintainer command for intentionally updating the pinned `golangci-lint` dependency.

## One Fix

Use pointer receivers and pass `*StreamConsumer` through APIs:

```go
func (c *StreamConsumer) Topic() string {
	return c.topic
}
```

This pattern is useful for handles that own goroutines, file descriptors, mutex-protected state, or lifecycle-sensitive resources.
