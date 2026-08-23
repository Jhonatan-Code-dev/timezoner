# Estándares de Arquitectura e Ingeniería Go (Top-Tier FAANG)

Este documento define las reglas de observancia obligatoria para cualquier desarrollo, modificación, refactorización o auditoría en el repositorio **timezoner**.

---

## 1. Arquitectura Modular Limpia (Modular Monolith)
- **Dominio Desacoplado en `pkg/`**: Cada funcionalidad especializada (`pkg/zone`, `pkg/types`, `pkg/calendar`, `pkg/humanize`, `pkg/ingest`, `pkg/project`, `pkg/overlap`) es un paquete autónomo con alta cohesión y bajo acoplamiento.
- **Fachada Pública Unificada**: El paquete raíz `timezoner` expone la Fluent API (`timezoner.At()`, `timezoner.Now()`) y constructores principales como punto de entrada único para los consumidores de la librería.
- **Regla de Dependencia (DIP)**: Los módulos de dominio no dependen de infraestructura externa. El árbol de dependencias internas debe ser un grafo acíclico dirigido (DAG).

## 2. Inmutabilidad y Protección de Estado
- **Cero Variables Globales Mutables Exportadas**: Ningún slice, mapa ni variable de paquete exportada puede ser modificable externamente.
- **Retorno de Copias Defensivas**: Toda función que exponga listas o catálogos internos (como `CommonZones()` o `SupportedLayouts()`) debe retornar una copia independiente (`copy()`).
- **Encapsulación de Tipos de Persistencia**: Tipos como `DBTime` y `ZonedTime` no deben embeber `time.Time` directamente; deben usar campos privados con accesores controlados (`.Time()`, `.UTC()`, `.IsZero()`) para prevenir mutaciones que corrompan la zona horaria en la base de datos.
- **Funciones Nombradas**: No declarar funciones exportadas como variables lambda `var Func = func()`; siempre usar declaraciones de función estándar `func Func()`.

## 3. Concurrencia y Seguridad (CWE-362)
- **Thread-Safety Total**: Toda función pública, estructura compartida y mecanismo de caché debe ser 100% seguro para ejecución concurrente con goroutines masivas.
- **Protección de Mapas**: Los mapas globales deben protegerse obligatoriamente con `sync.RWMutex` (lecturas con `RLock()`, escrituras con `Lock()`) o `sync.Map`.
- **Cero Fugas de Goroutines**: Cualquier proceso asíncrono debe tener un ciclo de vida delimitado y cancelable mediante `context.Context`.

## 4. Rendimiento y Cero Alocaciones Innecesarias
- **Optimización de Rutas Críticas**: Métodos de alta frecuencia (`NewDBTime`, `LoadLocation` en caché, conversiones) deben ejecutar en sub-microsegundos con $\le 1$ alocación en el heap.
- **Prealocación de Memoria**: Siempre que se conozca la longitud o capacidad de una colección, inicializar con `make([]T, 0, cap)`.
- **Caché en Memoria**: Mantener caché concurrente para resolución de `*time.Location`.
- **Evidencia Empírica**: Toda optimización de rendimiento debe ser validada con benchmarks reproducibles (`go test -bench=. -benchmem`).
