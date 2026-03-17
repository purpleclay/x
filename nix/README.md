# Nix

Types and functions for for interoperating with the Nix package manager, specifically its binary cache protocol.

```sh
go get github.com/purpleclay/x/nix
```

## Benchmarks

Benchmarks are run on an Apple M4 Pro (darwin/arm64) using Go's built-in testing framework.

```
go test -bench=. -benchmem ./...
```

```sh
goos: darwin
goarch: arm64
pkg: github.com/purpleclay/x/nix/base32
cpu: Apple M4 Pro
BenchmarkEncode-12      25347848                46.97 ns/op            0 B/op          0 allocs/op
BenchmarkDecode-12      21737488                56.20 ns/op            0 B/op          0 allocs/op
PASS
ok      github.com/purpleclay/x/nix/base32      2.855s
```
