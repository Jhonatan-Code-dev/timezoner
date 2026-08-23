---
name: iso-29119-test-verifier
description: >-
  Lead ISO/IEC 29119 Software Testing Standards Verifier for Go. Enforces boundary value analysis (BVA),
  equivalence partitioning, mutation testing, fuzz testing, and exhaustive test documentation.
---

# ISO/IEC 29119 Software Testing Standards Verifier

This skill enforces compliance with the international standard **ISO/IEC 29119** (Software and systems engineering — Software testing), ensuring rigorous verification, boundary exploration, and regression prevention.

---

## 🧪 ISO/IEC 29119 Test Design Checklist

### 1. Dynamic Test Design Techniques
- **Boundary Value Analysis (BVA)**:
  - Epoch boundaries (`1970-01-01 00:00:00 UTC`, negative timestamps, Year 2038 / 32-bit vs 64-bit limits).
  - Leap years (Feb 28/29 in leap and non-leap years).
  - DST boundary transitions (02:00 -> 03:00 and 02:00 -> 01:00).
  - Maximum/minimum timezone offsets (UTC-12:00 to UTC+14:00).
- **Equivalence Partitioning (EP)**:
  - Valid IANA zones, valid aliases, invalid zones, empty strings, malformed strings.

### 2. Adversarial & Random Input Exploration (Fuzzing)
- Mandatory `Fuzz<Target>(f *testing.F)` for all parsers and converters.
- Verification that arbitrary corrupted byte streams never induce panics or unhandled exceptions.

### 3. Concurrency Verification
- High-concurrency stress tests with race detector enabled (`go test -race`).
- Testing simultaneous parallel access across hundreds of goroutines.

### 4. Traceability & Executable Specification
- Every public capability is tested with an end-to-end integration test and documented with a runnable `Example...()` test.
