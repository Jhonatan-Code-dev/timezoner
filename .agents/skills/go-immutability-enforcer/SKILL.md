---
name: go-immutability-enforcer
description: >-
  Go Immutability & State Protection Enforcer. Prevents mutable global state, exported mutable
  slices/maps, dangerous type embedding, and ensures defensive copies on all public collection returns.
---

# Go Immutability & State Protection Enforcer

This skill enforces strict immutability rules to prevent external consumers from corrupting the internal state of a Go library through accidental or intentional mutation.

---

## Immutability Rules

### 1. Never Export Mutable Slices or Maps
```go
// FORBIDDEN — Anyone can overwrite elements:
var SupportedLayouts = []string{"2006-01-02", "2006-01-02 15:04"}

// CORRECT — Return a defensive copy every time:
var supportedLayouts = []string{"2006-01-02", "2006-01-02 15:04"} // lowercase = unexported

func SupportedLayouts() []string {
    out := make([]string, len(supportedLayouts))
    copy(out, supportedLayouts)
    return out
}
```

### 2. Never Use Type Embedding to Expose Parent Methods Unintentionally
```go
// FORBIDDEN — Exposes dangerous methods like Local(), Add(), AddDate():
type DBTime struct { time.Time }

// Consumer can do: myDBTime.Local() -> wrong timezone in DB!
// Consumer can do: myDBTime.Add(time.Hour) -> no UTC guarantee

// CORRECT — Private field with controlled accessors:
type DBTime struct { t time.Time }
func (d DBTime) Time() time.Time  { return d.t }
func (d DBTime) IsZero() bool     { return d.t.IsZero() }
```

### 3. Never Export Package-Level Mutable State
```go
// FORBIDDEN:
var DefaultTimeout = 30 * time.Second // Consumer can do: timezoner.DefaultTimeout = -1

// CORRECT: Use functional options or constructor parameters.
```

### 4. Returning Collections: Always Copy
```go
func CommonZones() []string {
    result := make([]string, len(commonZonesList))
    copy(result, commonZonesList) // Consumer cannot mutate the original
    return result
}
```

### 5. Value Types vs Pointer Types
- `DBTime` and `ZonedTime` should be value types (no pointer receiver on constructors).
- Value semantics prevent aliasing bugs where two variables unexpectedly share state.
