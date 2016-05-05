# Hacking PrefixMap

Profiling benchmarks
---

While making this project I've taken advantage of [pprof](https://golang.org/pkg/net/http/pprof/) tool provided by Golang.

Profiling the benchmarks is as easy as follows:

### Step 1.
```bash
go test -c && ./prefixmap.test -test.cpuprofile=cpu.prof -test.bench=<BenchmarkName>
```

### Step 2.
```bash
go tool pprof -text stringmap.test cpu.prof
```
