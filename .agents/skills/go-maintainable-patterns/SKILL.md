---
name: go-maintainable-patterns
description: >-
  Go Clean Code, Maintainable Architecture & Design Patterns Specialist. Enforces functional options,
  strict backward compatibility (SemVer), interface segregation, immutability, and production-grade library idioms.
---

# Go Maintainable Architecture & Production-Grade Design Patterns

This skill guides the creation of clean, maintainable, and backward-compatible Go packages and modules. It provides battle-tested design patterns and engineering practices that ensure libraries remain robust, easy to test, and effortless to maintain over years of evolution.

---

## 🏛️ Essential Go Library Patterns

### 1. Functional Options Pattern for Configurable APIs
When constructing configurable instances, avoid complex constructors with long parameter lists or breaking changes on new options. Use the Functional Options pattern:

```go
type Options struct {
    Timeout     time.Duration
    MaxRetries  int
    CacheSize   int
    Logger      Logger
}

type Option func(*Options)

func WithTimeout(d time.Duration) Option {
    return func(o *Options) {
        if d > 0 {
            o.Timeout = d
        }
    }
}

func NewService(opts ...Option) *Service {
    cfg := defaultOptions()
    for _, opt := range opts {
        opt(&cfg)
    }
    return &Service{cfg: cfg}
}
```

### 2. Monotonic vs. Wall Clock Discipline
- **Comparisons & Durations**: Use monotonic time (`t1.Sub(t0)`).
- **Serialization & Persistence**: Strip monotonic readings using `t.UTC()` or `t.Round(0)` before serializing or storing in databases.

### 3. Strict Semantic Versioning & Non-Breaking API Evolution
- **Never Change Exported Function Signatures**: Add new functions (e.g. `ConvertWithContext`) rather than modifying existing ones.
- **Never Remove Exported Fields from Public Structs**: Use unexported fields or extensible options to avoid breaking caller initialization.
- **Deprecation Lifecycle**: Mark obsolete APIs with standard Godoc comments:
  ```go
  // Deprecated: Use ConvertWithLocation instead.
  func ConvertLegacy(...)
  ```

### 4. Immutability & Concurrency Safety
- Package functions and methods must not mutate inputs passed by callers.
- Return new instances or copies of slices/maps to prevent data races in concurrent caller code.

---

## 🔍 Code Review & Maintainability Checklist

- [ ] **No Global Mutables**: Package level variables are read-only (`var ErrNotFound = ...`) or strictly synchronized.
- [ ] **Interface Pollution Check**: Interfaces are defined by consumers where needed, not prematurely declared in the library root.
- [ ] **Context Propagation**: Blocking, I/O-bound, or long-running operations accept `context.Context` as the first parameter.
- [ ] **Deterministic Behavior**: Package functions produce identical outputs for identical inputs without hidden side-effects.
- [ ] **Defensive Copying**: Slices and maps returned from package internals are cloned to protect internal state.
