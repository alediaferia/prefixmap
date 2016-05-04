Profiling benchmarks
---

1.   
```bash
go test -c && ./stringmap.test -test.cpuprofile=cpu.prof -test.bench=<BenchmarkName>
```
2.   
```bash
go tool pprof -text stringmap.test cpu.prof
```
