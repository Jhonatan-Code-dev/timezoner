---
name: go-benchmark-engineer
description: >-
  Go Benchmark Engineer. Designs, implements and interprets production-grade benchmarks using
  go test -bench. Identifies allocation hotspots, CPU cache misses, and validates performance
  claims with empirical evidence.
---

# Go Benchmark Engineer

This skill enforces that performance claims in Go libraries are backed by empirical, reproducible measurements using the Go testing framework's built-in benchmark infrastructure.

---

## Benchmark Design Rules

### 1. Every Performance-Critical Path Must Have a Benchmark
```go
func BenchmarkLoadLocation_Cached(b *testing.B) {
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = zone.LoadLocation("America/Lima")
    }
}
```

### 2. Always Use `-benchmem` to Detect Allocations
```
go test -bench=. -benchmem -run=^Benchmark ./...
```
- `allocs/op` must be 0 or 1 for hot-path operations.
- `B/op` must be minimal. Any unexpected growth signals a regression.

### 3. Prevent Compiler Optimization (Sink Variables)
```go
// Use package-level vars to prevent the compiler from eliminating benchmark code:
var sinkTime time.Time
var sinkErr  error
var sinkStr  string

func BenchmarkConvert(b *testing.B) {
    t := time.Now()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        sinkTime, sinkErr = zone.LoadLocation("America/Lima")
    }
}
```

### 4. Benchmark Cold vs Warm Cache Paths Separately
- Warm-cache benchmarks test the common case (sync.Map hit).
- Cold-cache benchmarks test first-time resolution and memory pressure.

### 5. Benchmark Interpretation Thresholds for timezoner
| Operation | Acceptable ns/op | Acceptable allocs/op |
| :--- | :---: | :---: |
| `LoadLocation` (cached) | < 100 ns | ≤ 1 |
| `NewDBTime` | < 20 ns | 0 |
| `IngestFromString` | < 500 ns | ≤ 2 |
| `ProjectForUser` | < 500 ns | ≤ 3 |
| `AddBusinessDays(5)` | < 500 ns | 0 |
| `FindOverlap(3 zones)` | < 1 ms | ≤ 10 |
