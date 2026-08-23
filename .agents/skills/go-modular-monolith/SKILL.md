---
name: go-modular-monolith
description: >-
  Modular Monolith & Scalable Systems Architect for Go. Enforces domain-driven module boundaries (pkg/),
  internal encapsulation (internal/), clear contracts, independent scalability, and microservices-ready design.
---

# Modular Monolith & Scalable Systems Architect in Go

This skill guides the design and evolution of **Modular Monoliths** in Golang. It combines the deployment simplicity and ultra-fast performance of a single binary with the clean module boundaries, bounded contexts, and independence of microservices.

---

## 🏛️ Modular Monolith Structure in Go

```
project-root/
├── cmd/                      # Application entry points (HTTP API, Workers, CLI)
│   └── server/main.go
├── pkg/                      # Independent public domain modules (Exported)
│   ├── zone/                 # Bounded context: Timezones & IANA catalog
│   ├── calendar/             # Bounded context: Business days & bounds
│   ├── ingest/               # Bounded context: Normalization pipeline
│   ├── project/              # Bounded context: User projections
│   └── types/                # Bounded context: SQL / JSON types (DBTime, ZonedTime)
├── internal/                 # Private package internals (Cannot be imported outside)
│   └── platform/             # Infrastructure shared primitives (Telemetry, DB pool)
├── examples/                 # Self-contained executable demonstrations
└── go.mod
```

---

## 📐 Core Tenets of Modular Monoliths

### 1. High Cohesion & Low Coupling (Bounded Contexts)
- Each module in `pkg/<module>/` represents a distinct bounded context.
- A module encapsulates its own data structures, business logic, errors, and isolated unit tests.

### 2. Explicit Cross-Module Communication
- Modules must not reach directly into the internal state or private structs of other modules.
- Communication happens strictly through public exported functions, contracts, or events.

### 3. Microservices-Ready Decoupling
- Because each module in `pkg/` is isolated and self-contained with its own tests, extracting any module into a separate microservice or independent Go repository requires zero architectural rewriting.

### 4. Single Deployment, Zero Distributed Latency
- High throughput and in-memory function call performance (nanoseconds) without network serialization overhead or network partition failures.
