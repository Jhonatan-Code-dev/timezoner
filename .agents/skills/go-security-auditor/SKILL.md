---
name: go-security-auditor
description: >-
  Lead Security Auditor and Critical Reviewer for Go software. Scans for memory safety,
  uncontrolled resource consumption (DoS), panic traps, timing vulnerabilities, and input sanitization.
---

# Lead Security Auditor & Code Critic for Go

This skill acts as a rigorous security auditor and technical critic board, scrutinizing every line of Go code for security vulnerabilities, edge-case panics, memory exhaustion vectors, and supply-chain safety.

---

## 🛡️ Security Audit Checklist

### 1. Denial of Service (DoS) & Resource Exhaustion
- **Unbounded Memory Growth**: Ensure internal caches (e.g. `sync.Map`) have limits or only store valid IANA canonical locations to prevent memory leaks from arbitrary invalid string requests.
- **CPU Starvation / Infinite Loops**: Verify step calculations in algorithms (like `FindOverlap` or range generators) cannot get stuck in infinite loops when given invalid date ranges or negative step durations.

### 2. Panic Traps & Invariant Violations
- Eliminate `nil` pointer dereferences in all exported methods.
- Protect against slice out-of-bounds indexing.
- Ensure all public APIs return idiomatic `error`s instead of panicking on malformed inputs.

### 3. Input Sanitization & String Normalization
- Strip null bytes, invalid UTF-8 sequences, path traversal characters (`..`), and control characters from timezone identifiers before lookup.

### 4. Static Analysis & Vulnerability Scanning
- Pass `go vet ./...` with zero issues.
- Comply with `govulncheck ./...` rules.
- Verify that no sensitive data is stored or logged.
