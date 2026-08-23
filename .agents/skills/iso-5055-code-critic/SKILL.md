---
name: iso-5055-code-critic
description: >-
  Lead ISO/IEC 5055 Code Quality & Architecture Flaw Critic. Scans source code for structural weaknesses,
  CWE security patterns, memory leakage vectors, data race hazards, and architectural coupling violations.
---

# ISO/IEC 5055 Automated Source Code Quality Critic

This skill enforces **ISO/IEC 5055** (Automated Source Code Quality Measures), which measures software quality directly at the source code level by detecting critical structural flaws, security weaknesses (mapped to CWEs), and performance bottlenecks.

---

## 🛡️ ISO/IEC 5055 Structural Dimension Checklist

### 1. Architectural Integrity & Coupling (Maintainability)
- **Zero Circular Package Dependencies**: Package graph is an acyclic directed graph (DAG).
- **Layer Violations**: Domain packages (`pkg/`) must never import delivery or adapter layers.
- **Dead Code Elimination**: 0 unused variables, unreferenced private functions, or abandoned constants.

### 2. Concurrency & Race Safety (Reliability)
- **Lock Contention**: Minimize critical sections in `sync.RWMutex` / `sync.Mutex`.
- **Shared Mutable State**: Forbid package-level mutable variables.
- **Goroutine Leak Prevention**: Any spawned goroutine must have a bounded lifetime.

### 3. Memory & Resource Safety (Security / CWEs)
- **CWE-400 (Uncontrolled Resource Consumption)**: Caches have bounds or only store validated canonical keys to prevent memory explosion.
- **CWE-476 (NULL Pointer Dereference)**: Defensive validation on all inputs and receivers.
- **CWE-190 (Integer Overflow)**: Safe time/duration conversions without numerical truncation.

### 4. Computational Efficiency (Performance)
- Preallocate slices (`make([]T, 0, cap)`) when the size is known.
- Avoid repetitive string concatenations in hot loops; use `strings.Builder` or integer math.
