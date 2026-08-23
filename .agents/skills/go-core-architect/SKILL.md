---
name: go-core-architect
description: >-
  Senior Principal Go Architect skill specializing in designing world-class,
  idiomatic Go modules and libraries with pristine API boundaries, zero external dependencies,
  thread safety, and clean abstractions.
---

# Senior Principal Go Architect & Package Designer

This skill embodies the standards of Principal Software Engineers from tier-1 technology companies (Google, Cloudflare, Uber, Docker). It guides the design, implementation, and evolution of production-grade, highly reusable Go libraries and packages.

---

## 🏛️ Core Principles of Library Design in Go

### 1. Minimal & Purpose-Driven API Surface
- **Accept Interfaces, Return Concrete Types**: Keep consumer interfaces small (1-2 methods) and define them where they are used, not where they are implemented.
- **Export Only What Is Necessary**: Keep internal helper structs, constants, and utilities unexported (`private`).
- **Idiomatic Naming**: Avoid package name repetition (`timezoner.Timezone` ❌ -> `timezoner.Zone` or `timezoner.Detail` ✔️).

### 2. Zero Unnecessary External Dependencies
- Rely strictly on the Go Standard Library (`time`, `sync`, `fmt`, `strings`, `errors`, `math`).
- Keep `go.mod` clean and dependency-free to avoid dependency hell for downstream projects.

### 3. Immaculate Concurrency & Thread Safety
- Any exported struct or package-level function that maintains state MUST be safe for concurrent use by multiple goroutines.
- Use `sync.Map`, `sync.RWMutex`, or atomic primitives where appropriate without introducing lock contention bottlenecks.
- Never expose mutable package-level global variables.

### 4. Error Handling Philosophy
- Provide clear, contextual errors wrapping standard errors with `%w`.
- Define sentinel errors (`ErrInvalidZone`, `ErrInvalidTimeFormat`) or custom error types so callers can inspect with `errors.Is` and `errors.As`.
- Never panic in library functions; return descriptive `error`s.

---

## 📐 Fluent API & Builder Guidelines
- Builder chains should propagate errors internally and evaluate them on terminal methods (`.Time() (time.Time, error)`).
- Provide convenience `Must*` methods (`.MustTime()`) only for scripts and tests, clearly documented.
