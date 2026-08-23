---
name: go-code-hygiene-order
description: >-
  Obsessive Code Hygiene, Style and Project Structure Master for Go. Enforces immaculate Godoc formatting,
  idiomatic Go naming, strict linter compliance, and pristine repository order.
---

# Obsessive Code Hygiene & Repository Order Master

This skill enforces obsessive cleanliness, symmetry, and aesthetic code discipline across every file, function, comment, and directory in the Go codebase.

---

## 🧹 Code Hygiene & Cleanliness Standards

### 1. Godoc Perfection
- Every exported package, struct, interface, function, method, and constant MUST begin with a descriptive sentence starting with its own name.
  ```go
  // TimePoint encapsulates a specific instant in time and provides a fluent API for timezone operations.
  type TimePoint struct { ... }
  ```
- Package-level doc comment at the top of the root file explaining the entire module's architecture and usage.

### 2. Standard Repository Layout
- Root package contains only core, clean domain files (`timezoner.go`, `zones.go`, `timezoner_test.go`, `examples_test.go`).
- Examples placed in clean, self-contained sub-packages (`examples/basic_usage/`, `examples/team_meeting_planner/`).
- No dangling temporary files, unformatted code (`gofmt -s -w .`), or messy git artefacts.

### 3. Strict Linter Compliance
- 0 warnings from `gofmt`, `go vet`, `revive`, `staticcheck`, and `golangci-lint`.
- Clean error message strings (no capitalized words, no trailing punctuation per Go style guide: `"timezoner: invalid zone name"`).

### 4. Deterministic Code Organization
- File contents organized in a fixed order:
  1. Package declaration & docstring
  2. Imports (Standard library separated from internal packages)
  3. Constants & Sentinel Errors
  4. Core Types & Structs
  5. Constructors / Factory Functions
  6. Methods & Public Functions
  7. Private Helpers
