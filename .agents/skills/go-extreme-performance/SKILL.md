---
name: go-extreme-performance
description: >-
  Extreme Performance and Low-Latency Go Optimization Specialist. Focuses on zero-allocations,
  CPU cache line locality, escape analysis, struct alignment, and high-throughput benchmarks.
---

# Extreme Performance & Low-Latency Go Specialist

This skill focuses on maximizing throughput and minimizing latency and memory footprint to enterprise ultra-low-latency standards.

---

## ⚡ Performance Optimization Framework

### 1. Escape Analysis & Zero-Allocation Patterns
- Check heap escapes using:
  ```bash
  go build -gcflags="-m -m" ./...
  ```
- Keep short-lived structures on the stack rather than escaping to the heap.
- Preallocate slice capacities with `make([]T, 0, expectedCap)` when the upper bound is known.

### 2. Struct Memory Layout & Padding
- Order struct fields by descending size (8-byte pointers/ints -> 4-byte ints -> 1-byte bools) to minimize compiler padding bytes and fit into CPU cache lines (64 bytes).

### 3. Location & Calculation Caching
- Cache frequently resolved objects with thread-safe lock-free or read-biased structures (`sync.Map`).
- Avoid repetitive string formatting allocations in hot loops; use buffer reuse or integer arithmetic when computing offsets.

### 4. Micro-benchmarking Standard
- Always measure changes using `go test -bench=. -benchmem -count=5`.
- Compare revisions with `benchstat` to prove statistical improvement.
