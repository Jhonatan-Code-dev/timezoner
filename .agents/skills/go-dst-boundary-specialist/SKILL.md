---
name: go-dst-boundary-specialist
description: >-
  DST Boundary & Calendar Arithmetic Specialist for Go. Enforces correct handling of Daylight
  Saving Time transitions, ambiguous local times, skipped times, AddDate pitfalls, and wall-clock
  preservation across timezone changes.
---

# DST Boundary & Calendar Arithmetic Specialist

This skill is the deepest level of temporal correctness enforcement. It focuses on the subtle bugs that only appear at DST transition boundaries, which affect less than 0.3% of days per year but cause 100% incorrect results when they occur.

---

## DST Critical Rules

### 1. AddDate Does Not Preserve Local Time Across DST Boundaries
```go
// PROBLEM: AddDate moves the underlying instant, but the local time shifts by 1 hour on DST day
loc, _ := time.LoadLocation("America/New_York")
fri := time.Date(2026, 3, 6, 9, 0, 0, 0, loc) // Friday 09:00 (before spring-forward)
mon := fri.AddDate(0, 0, 3)                     // Monday 10:00 — WRONG! DST added 1 hour

// CORRECT: Reconstruct the local time components after AddDate
hour, min, sec := fri.Clock()
result := fri.AddDate(0, 0, 3)
year, month, day := result.Date()
fixed := time.Date(year, month, day, hour, min, sec, 0, loc) // Monday 09:00 — CORRECT
```

### 2. Ambiguous Times (Fall-Back / "Fold" Times)
When clocks go back (e.g., 02:00 → 01:00), the local time 01:30 exists twice.
```go
// Go 1.9+ exposes time.Time.IsDST() — but cannot tell which of the two 01:30 is which.
// Solution: Always store UTC and never create ambiguous wall-clock times directly.
```

### 3. Skipped Times (Spring-Forward)
When clocks go forward (e.g., 02:00 → 03:00), 02:30 does not exist.
```go
// time.Date(year, month, day, 2, 30, 0, 0, nyLoc) returns 03:30 silently.
// Always validate user-provided local times against this scenario.
```

### 4. DST Detection Algorithm (Standard)
```go
// Compare January and July offsets to determine the standard (non-DST) offset.
jan := time.Date(year, time.January, 1, 12, 0, 0, 0, loc)
jul := time.Date(year, time.July, 1, 12, 0, 0, 0, loc)
_, janOff := jan.Zone()
_, julOff := jul.Zone()
stdOffset := janOff
if julOff < janOff {
    stdOffset = julOff // Southern hemisphere: DST is in January
}
_, currentOffset := t.In(loc).Zone()
isDST := currentOffset > stdOffset
```

### 5. Business Day Arithmetic DST Test
```go
// Must test AddBusinessDays across the DST transition weekend.
loc, _ := time.LoadLocation("America/New_York")
preDST := time.Date(2026, 3, 6, 9, 0, 0, 0, loc) // Friday before spring-forward
result := AddBusinessDays(preDST, 1) // Should be Monday at 09:00, NOT 10:00
assert(result.Hour() == 9)
```
