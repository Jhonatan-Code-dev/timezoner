---
name: go-api-design-expert
description: >-
  Senior Go Public API Design Expert. Enforces ergonomic API surfaces, backward compatibility,
  naming conventions, builder patterns, and prevention of mutable global functions that can be
  overwritten at runtime.
---

# Go Public API Design Expert

This skill acts as a ruthless API design critic for Go packages and libraries. Its goal is to produce APIs that are intuitive, safe, non-breaking, and impossible to misuse.

---

## Critical API Design Rules

### 1. Never Export Mutable Function Variables (Security / Misuse Prevention)
```go
// FORBIDDEN — Any consumer can overwrite this:
var ZonedFromLocal = func(s, z string) (ZonedTime, error) { ... }

// CORRECT — A named function cannot be replaced:
func ZonedFromLocal(s, z string) (ZonedTime, error) { ... }
```

### 2. Never Embed External Types in Domain Types (API Bloat)
```go
// FORBIDDEN — Exposes 50+ methods including dangerous ones like Local(), Add():
type DBTime struct { time.Time }

// CORRECT — Controlled accessors only:
type DBTime struct { t time.Time }
func (d DBTime) Time() time.Time { return d.t }
func (d DBTime) UTC() time.Time  { return d.t }
func (d DBTime) IsZero() bool    { return d.t.IsZero() }
```

### 3. Never Export Mutable Slices or Maps Directly
```go
// FORBIDDEN — Consumer can corrupt the entire package state:
var SupportedLayouts = []string{"2006-01-02", ...}

// CORRECT — Return a defensive copy:
func SupportedLayouts() []string {
    out := make([]string, len(internal))
    copy(out, internal)
    return out
}
```

### 4. Builder / Fluent API Must Accumulate Errors (Railway Pattern)
- Every method in the chain stores the first error in `tp.err`.
- Terminal methods (`MustTime()`, `MustDBTime()`) panic only when `err != nil`.
- Safe variants (`Time()`, `AsDBTime()`) return `(T, error)` for production code.
- Always expose `Err() error` to allow inspection without termination.

### 5. Semantic Versioning and Backward Compatibility
- Never remove or rename exported symbols between minor versions.
- New fields in public structs must use pointer types or have zero-value semantics.
- New parameters use functional options (`WithLayout(l string) Option`).
