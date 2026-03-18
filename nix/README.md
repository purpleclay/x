# Nix

Types and functions for interoperating with the Nix package manager, specifically its binary cache protocol.

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
pkg: github.com/purpleclay/x/nix
cpu: Apple M4 Pro
BenchmarkParseHash/Nix32-12          	16526563	        68.28 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseHash/Base64-12         	42235736	        29.00 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseHash/SRI-12            	45654208	        26.37 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseHash/Hex-12            	41468682	        29.60 ns/op	       0 B/op	       0 allocs/op
BenchmarkHashOutput/RawHex-12        	32113000	        36.86 ns/op	     128 B/op	       2 allocs/op
BenchmarkHashOutput/RawBase32-12     	18220282	        65.43 ns/op	     128 B/op	       2 allocs/op
BenchmarkHashOutput/RawBase64-12     	39197059	        30.84 ns/op	      96 B/op	       2 allocs/op
BenchmarkHashOutput/SRI-12           	35127734	        33.53 ns/op	     128 B/op	       2 allocs/op
BenchmarkCompressHash-12             	72540246	        17.32 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/purpleclay/x/nix	11.088s
goos: darwin
goarch: arm64
pkg: github.com/purpleclay/x/nix/base32
cpu: Apple M4 Pro
BenchmarkEncode-12    	22963064	        46.71 ns/op	       0 B/op	       0 allocs/op
BenchmarkDecode-12    	20772042	        56.99 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/purpleclay/x/nix/base32	2.416s
```
