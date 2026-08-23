---
name: iso-25010-quality-auditor
description: >-
  Lead ISO/IEC 25010 Software Quality Standards Auditor for Go. Rigorously inspects systems across
  the 8 quality characteristics: Maintainability, Reliability, Performance Efficiency, Security,
  Portability, Compatibility, Usability, and Functional Suitability.
---

# ISO/IEC 25010 Software Product Quality Auditor (SQuaRE)

This skill acts as a ruthless quality auditor evaluating Go codebases against the international standard **ISO/IEC 25010** (Systems and software engineering — Systems and software Quality Requirements and Evaluation).

---

## 🏛️ The 8 ISO/IEC 25010 Quality Characteristics

```
 ┌────────────────────────────────────────────────────────────────────────┐
 │                      ISO/IEC 25010 QUALITY MODEL                       │
 ├───────────────────┬───────────────────┬────────────────────────────────┤
 │ 1. Maintainability│ 2. Reliability    │ 3. Performance Efficiency      │
 │ • Modularity      │ • Fault tolerance │ • Time behavior (latency)      │
 │ • Reusability     │ • Recoverability  │ • Resource utilization (allocs)│
 │ • Testability     │ • Maturity        │ • Capacity                     │
 ├───────────────────┼───────────────────┼────────────────────────────────┤
 │ 4. Security       │ 5. Portability    │ 6. Compatibility               │
 │ • Confidentiality │ • Adaptability    │ • Co-existence                 │
 │ • Integrity       │ • Installability  │ • Interoperability             │
 │ • Accountability  │ • Replaceability  │                                │
 ├───────────────────┴───────────────────┴────────────────────────────────┤
 │ 7. Functional Suitability             │ 8. Usability                   │
 │ • Completeness, Correctness, Precision│ • Learnability, Operability    │
 └────────────────────────────────────────────────────────────────────────┘
```

---

## 🔍 ISO/IEC 25010 Audit Checklist for Go Packages

### 1. Maintainability (Sub-characteristics: Modularity, Testability, Analyzability)
- **Modularity**: Code is partitioned into bounded contexts (`pkg/zone`, `pkg/types`, `pkg/calendar`).
- **Low Cyclomatic Complexity**: Functions have single responsibilities; complexity score $\le 10$.
- **Zero Interface Pollution**: Interfaces declared where consumed.

### 2. Reliability (Sub-characteristics: Fault Tolerance, Recoverability)
- **Panic Immunity**: Public library APIs NEVER panic on invalid user input; always return structured `error`s.
- **Nil Safety**: Pointer receivers and inputs check for `nil` before dereferencing.

### 3. Performance Efficiency (Sub-characteristics: Time Behavior, Resource Utilization)
- **Zero-Allocation Critical Paths**: Hot paths avoid escaping heap allocations.
- **Concurrent In-Memory Caching**: Repetitive lookups are cached thread-safely (`sync.Map`).

### 4. Portability (Sub-characteristics: Adaptability, Installability)
- **Embedded Dependencies (`_ "time/tzdata"`)**: Works universally on Windows, Linux, Alpine, and Scratch Docker containers without OS dependencies.
