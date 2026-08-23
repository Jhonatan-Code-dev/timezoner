---
name: go-temporal-datetime-master
description: >-
  Temporal & DateTime Engineering Specialist in Go. Deep expertise in time.Time internals,
  monotonic clocks, IANA tzdata, DST boundary transitions, ISO 8601/RFC 3339 compliance,
  temporal interval arithmetic, and business calendar calculations.
---

# Temporal & DateTime Engineering in Go

This skill provides comprehensive engineering guidelines for building mission-critical date, time, and timezone manipulation packages in pure Go.

---

## ⏱️ Core Knowledge & Principles

### 1. `time.Time` Internals & Clocks in Go
- **Wall Clock vs. Monotonic Clock**:
  - In Go 1.9+, `time.Now()` contains both a wall clock reading (calendar time) and a monotonic clock reading (for measurement).
  - Calling `t.Equal(u)` compares instantaneous real time and respects monotonic readings.
  - Calling `==` directly between two `time.Time` values compares both the wall clock and the monotonic clock. Always use `t.Equal(u)` to compare timestamps across process boundaries or serialization.
  - Calling `t.Round(0)` or `t.UTC()` strips the monotonic reading for deterministic equality and serialization.

### 2. Timezone Offsets & IANA Database
- **Non-Integral Hour Offsets**:
  - Never assume offsets are always multiples of 60 minutes.
  - Support half-hour offsets (e.g. `Asia/Kolkata` +05:30, `Asia/Tehran` +03:30, `Australia/Adelaide` +09:30).
  - Support quarter-hour offsets (e.g. `Asia/Kathmandu` +05:45, `Pacific/Chatham` +12:45).
- **Daylight Saving Time (DST) Transitions**:
  - **Spring Forward**: 02:00 -> 03:00 (1 hour does not exist in local time).
  - **Fall Back**: 02:00 -> 01:00 (1 hour occurs twice).
  - Conversions between local wall clock and UTC must handle ambiguous and non-existent local times gracefully.

### 3. Date Interval Arithmetic & Overlaps
- Standard Interval Intersection Formula:
  Two intervals `[A_start, A_end)` and `[B_start, B_end)` overlap if and only if:
  $$\text{Overlap} \iff A_{\text{start}} < B_{\text{end}} \land B_{\text{start}} < A_{\text{end}}$$
- Interval Overlap Duration:
  $$\text{Start} = \max(A_{\text{start}}, B_{\text{start}})$$
  $$\text{End} = \min(A_{\text{end}}, B_{\text{end}})$$
  $$\text{Duration} = \text{End} - \text{Start}$$

### 4. Leap Years and Date Math
- Always use `time.Date(year, month, day, ...)` rather than manual second additions (`+ 24*time.Hour`) when performing calendar day addition, because daylight saving transitions make some calendar days 23 or 25 hours long.
- A year is a leap year if:
  $$(\text{year} \pmod 4 == 0 \land \text{year} \pmod{100} \neq 0) \lor (\text{year} \pmod{400} == 0)$$

### 5. Serialization Standards
- **RFC 3339 / ISO 8601**: Standard interchange layout (`2006-01-02T15:04:05Z07:00` or `time.RFC3339Nano`).
- **Database Storage**: Always store timestamps in UTC (`t.UTC()`) with integer epoch offsets or RFC 3339 strings.
