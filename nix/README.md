# Nix

Types and functions for interoperating with the Nix package manager, specifically its binary cache protocol.

```sh
go get github.com/purpleclay/x/nix
```

## Fuzz tests

Fuzz tests verify that parsers never panic on arbitrary input. Each target ships
a seed corpus of real inputs; running without `-fuzz` replays the corpus only
(fast, suitable for CI):

```sh
go test -run=FuzzParseHash ./...
```

To run the full fuzzer locally, pass `-fuzz` and a time limit:

```sh
go test -fuzz=FuzzParseHash        -fuzztime=60s
go test -fuzz=FuzzNARInfoUnmarshalText -fuzztime=60s
go test -fuzz=FuzzParseStorePath   -fuzztime=60s
go test -fuzz=FuzzParsePublicKey   -fuzztime=60s
go test -fuzz=FuzzParseSignature   -fuzztime=60s
```

The `base32` sub-package has two additional targets — one for panic safety and
one that checks the encode→decode round-trip property for all byte inputs:

```sh
cd base32
go test -fuzz=FuzzDecode                -fuzztime=60s
go test -fuzz=FuzzEncodeDecodeRoundTrip -fuzztime=60s
```

Any input that triggers a failure is automatically saved to
`testdata/fuzz/<FuzzName>/` and replayed on every subsequent `go test` run.

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
BenchmarkParseHash/Hex-12               	100000000	        30.39 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseHash/Nix32-12             	 51423673	        69.91 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseHash/Base64-12            	120269392	        29.91 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseHash/SRI-12               	127293162	        28.25 ns/op	       0 B/op	       0 allocs/op
BenchmarkHashOutput/RawHex-12           	 89672506	        39.38 ns/op	     128 B/op	       2 allocs/op
BenchmarkHashOutput/RawBase32-12        	 49447350	        71.75 ns/op	     128 B/op	       2 allocs/op
BenchmarkHashOutput/RawBase64-12        	100000000	        33.29 ns/op	      96 B/op	       2 allocs/op
BenchmarkHashOutput/SRI-12              	100000000	        35.65 ns/op	     128 B/op	       2 allocs/op
BenchmarkCompressHash-12                	208906142	        17.25 ns/op	       0 B/op	       0 allocs/op
BenchmarkNARInfoMarshalText-12          	  4236165	       854.2 ns/op	    2097 B/op	      22 allocs/op
BenchmarkNARInfoUnmarshalText-12        	  6391402	       563.0 ns/op	     840 B/op	       5 allocs/op
BenchmarkSignNARInfo-12                 	   258193	       13920 ns/op	    1224 B/op	      16 allocs/op
BenchmarkVerifyNARInfo-12               	   125496	       28841 ns/op	    1144 B/op	      15 allocs/op
PASS
ok  	github.com/purpleclay/x/nix	45.988s
goos: darwin
goarch: arm64
pkg: github.com/purpleclay/x/nix/base32
cpu: Apple M4 Pro
BenchmarkEncode-12    	65781057	        52.26 ns/op	       0 B/op	       0 allocs/op
BenchmarkDecode-12    	61282114	        58.72 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/purpleclay/x/nix/base32	7.300s
```
