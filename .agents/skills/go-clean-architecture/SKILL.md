---
name: go-clean-architecture
description: >-
  Senior Clean Architecture, Hexagonal & DDD Specialist for Go. Enforces strict layer decoupling,
  dependency inversion (DIP), domain entities, value objects, use cases, and port/adapter boundaries.
---

# Senior Clean Architecture & Hexagonal Specialist in Go

This skill enforces strict Clean Architecture, Hexagonal Architecture (Ports & Adapters), and Domain-Driven Design (DDD) principles in Golang codebases to ensure maximum testability, maintainability, and independence from frameworks, drivers, and external infrastructure.

---

## 🏛️ Layer Organization & Dependency Rule

The core invariant of Clean Architecture is the **Dependency Rule**: *Source code dependencies must point only inward, toward higher-level policies.*

```
 ┌────────────────────────────────────────────────────────┐
 │ 4. Infrastructure / Adapters (DB, HTTP, CLI, Queues)   │
 │   ┌──────────────────────────────────────────────────┐ │
 │   │ 3. Interface Adapters (Controllers, Presenters) │ │
 │   │   ┌────────────────────────────────────────────┐ │ │
 │   │   │ 2. Application / Use Cases                 │ │ │
 │   │   │   ┌──────────────────────────────────────┐ │ │ │
 │   │   │   │ 1. Enterprise Domain / Entities      │ │ │ │
 │   │   │   └──────────────────────────────────────┘ │ │ │
 │   │   └────────────────────────────────────────────┘ │ │ │
 │   └──────────────────────────────────────────────────┘ │ │
 └────────────────────────────────────────────────────────┘
```

---

## 📐 Core Architecture Principles in Go

### 1. Domain Layer (Zero Dependencies)
- Contains pure business logic, Domain Entities, Value Objects, and Domain Errors.
- **Rule**: Imports ONLY standard Go packages. Never imports SQL drivers, HTTP frameworks, or external infrastructure.

### 2. Ports (Interfaces) Owned by the Consumer
- Interfaces (ports) are defined inside the application/use-case layer that needs them, not in the database/infrastructure package.
  ```go
  // Inside domain/application layer:
  type TimezoneRepository interface {
      SaveSchedule(ctx context.Context, s Schedule) error
      FindByID(ctx context.Context, id string) (Schedule, error)
  }
  ```

### 3. Adapters (Implementations)
- Database repositories (PostgreSQL, MySQL), HTTP Handlers, and CLI commands implement the ports without the domain knowing about their concrete implementations.

### 4. Dependency Inversion Principle (DIP)
- High-level modules do not depend on low-level modules. Both depend on abstractions.
- Use constructors (`NewService(repo TimezoneRepository)`) to inject dependencies explicitly at startup (Composition Root).
