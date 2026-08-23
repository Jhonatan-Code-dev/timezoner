# Estándares de Arquitectura e Ingeniería Go (Nivel 100 — Drástico)

Este documento define las reglas de ingeniería de observancia **innegociable** para cualquier desarrollo, modificación, refactorización o auditoría en el repositorio **timezoner**.

---

## 1. Arquitectura Modular Limpia (Modular Monolith)
- **Dominio Desacoplado en `pkg/`**: Cada módulo (`pkg/zone`, `pkg/types`, `pkg/calendar`, `pkg/humanize`, `pkg/ingest`, `pkg/project`, `pkg/overlap`) debe ser 100% independiente y autosuficiente.
- **Fachada Pública Unificada**: El paquete raíz `timezoner` es el único punto de contacto oficial para los consumidores.
- **Árbol de Dependencias Aislado**: Los paquetes en `pkg/` jamás importan paquetes de aplicación ni controladores.

## 2. Inmutabilidad y Protección Absoluta de Estado
- **Prohibición Total de Variables Globales Mutables Exportadas**: Cualquier `var` exportada mutable (slices, maps, variables de configuración) es considerada un defecto crítico inmediato.
- **Copias Defensivas Obligatorias**: Toda función que exponga slices o listas internas (`CommonZones()`, `SupportedLayouts()`) debe ejecutar `copy()` a un nuevo slice antes de retornar.
- **Tipos de Persistencia Encapsulados**: `DBTime` y `ZonedTime` tienen campos privados y accesores controlados (`.Time()`, `.UTC()`, `.IsZero()`). Jamás se permite embeber `time.Time` directamente.
- **Funciones Nombradas Inmutables**: Prohibido usar variables de tipo función `var Func = func()`. Siempre usar declaraciones de función estándar `func Func()`.

## 3. Concurrencia y Detección de Carreras (CWE-362)
- **Thread-Safety Obligatorio**: Todo mapa global o caché debe estar respaldado por `sync.RWMutex` (lecturas con `RLock()`, escrituras con `Lock()`) o `sync.Map`.
- **Inmunidad a Deadlocks**: Bloqueos mínimos; jamás invocar funciones desconocidas o de terceros dentro de una sección crítica bloqueada.

## 4. Rendimiento Extremo con Evidencia Empírica
- **Límites de Rendimiento en Rutas Calientes**:
  - `LoadLocation` (con caché): $< 100\text{ ns/op}$, $\le 1\text{ alloc/op}$.
  - `NewDBTime`: $< 20\text{ ns/op}$, $0\text{ allocs/op}$.
  - `AddBusinessDays`: $< 500\text{ ns/op}$, $0\text{ allocs/op}$.
- **Prealocación de Slices**: Uso obligatorio de `make([]T, 0, cap)` cuando la dimensión es determinable.
- **Prueba Obligatoria**: Toda afirmación de velocidad debe incluir la salida textual de `go test -bench=. -benchmem`.

## 5. Inmunidad a Pánicos y Manejo de Errores
- **Cero Pánicos en APIs Públicas**: Las funciones públicas ordinarias nunca deben hacer `panic()`. Solo las variantes con prefijo explícito `Must*` (`MustTime()`, `MustDBTime()`) tienen permitido el pánico ante error.
- **Errores Centinela Estructurados**: Todo error debe usar errores tipados / centinela envueltos con `%w` (`ErrInvalidZone`, `ErrEmptyDateString`, `ErrNoZonesProvided`).
