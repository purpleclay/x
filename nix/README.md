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
