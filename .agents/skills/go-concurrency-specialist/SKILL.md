---
name: go-concurrency-specialist
description: >-
  Go Concurrency & Thread-Safety Specialist. Enforces correct use of sync.RWMutex, sync.Map,
  atomic operations, goroutine lifetime control, and detection of data races (CWE-362).
---

# Go Concurrency & Thread-Safety Specialist

This skill reviews Go code for concurrent correctness. A library consumed by HTTP servers, workers, and async pipelines must be 100% goroutine-safe on every exported symbol.

---

## Core Concurrency Rules

### 1. Global Maps Require sync.RWMutex (CWE-362)
```go
// FORBIDDEN — map is NOT goroutine-safe for concurrent reads + writes:
var aliases = map[string]string{"PET": "America/Lima"}
aliases["NEW"] = "America/Lima" // concurrent write: CRASH

// CORRECT — RWMutex allows concurrent reads, exclusive writes:
var mu sync.RWMutex
var aliases = map[string]string{"PET": "America/Lima"}

// Read (many goroutines can read simultaneously):
mu.RLock()
val := aliases[key]
mu.RUnlock()

// Write (exclusive, blocks all readers):
mu.Lock()
aliases[key] = value
mu.Unlock()
```

### 2. When to Use sync.Map vs sync.RWMutex
| Scenario | Recommended |
| :--- | :--- |
| Map is written once at startup, then only read | `sync.Map` (Load only) |
| Map is frequently written at runtime | `sync.RWMutex` + regular map |
| Cache that grows over time | `sync.Map` (Store + Load) |

### 3. Goroutine Leak Prevention
- Every spawned goroutine must have a bounded lifetime via `context.Context` or `done chan struct{}`.
- Libraries must never spawn background goroutines without the caller's explicit consent.

### 4. Detecting Races
```bash
go test -race ./...
# If CGO is not available (Windows):
CGO_ENABLED=1 go test -race ./...
```

### 5. Atomicity Traps
```go
// Non-atomic increment — RACE:
counter++

// Correct atomic:
atomic.AddInt64(&counter, 1)
```
