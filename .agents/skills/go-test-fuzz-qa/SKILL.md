---
name: go-test-fuzz-qa
description: >-
  Elite QA, Fuzzing & Mutation Testing Specialist for Go. Enforces 100% boundary testing,
  property-based tests, race condition hunting (-race), fuzz testing, and zero-allocation assertions.
---

# Elite QA & Fuzzing Specialist for Go

This skill applies extreme testing methodologies to ensure that the Go package achieves a 10/10 level of reliability, robustness, and stability under adversarial inputs and high concurrency.

---

## 🧪 Testing Checklist & Methodologies

### 1. Fuzz Testing (`go test -fuzz`)
- Implement native Go fuzz targets (`FuzzConvert`, `FuzzParse`, `FuzzFindOverlap`) with random byte slices, malformed zone names, leap years, epoch boundaries, and corrupted dates.
- Assert that functions never crash or panic on any input string.

### 2. Concurrency & Race Detector (`-race`)
- Execute tests under high concurrency (`go test -race -count=10 ./...`).
- Verify simultaneous read/write operations on internal caches or lookup maps across thousands of goroutines.

### 3. Table-Driven Tests & Edge Cases
- **Temporal Edge Cases**:
  - Solstices & Equinoxes.
  - Daylight Saving Time transition boundaries (the "lost hour" and "repeated hour" at 02:00 -> 03:00).
  - Leap seconds and Leap years (e.g. Feb 29, 2024, 2028).
  - Timezone offsets with half-hour and 45-minute increments (e.g. `Asia/Kolkata` +05:30, `Asia/Kathmandu` +05:45, `Australia/Eucla` +08:45).
  - Extreme coordinates / UTC-12:00 to UTC+14:00 (Line Islands / Kiritimati).

### 4. Zero-Allocation Assertions in Tests
- Track heap allocations using `testing.AllocsPerRun`.
- Ensure critical conversion paths produce 0 heap allocations when locations are cached.

### 5. Executable Godoc Examples
- Every public API method must have a corresponding `Example<Function>()` in `examples_test.go` with verifiable `// Output:` comments checked by `go test`.
