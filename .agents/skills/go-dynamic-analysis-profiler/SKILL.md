---
name: go-dynamic-analysis-profiler
description: >-
  Go Dynamic Analysis & Profiler Expert. Analiza el comportamiento real del código mientras se
  ejecuta: uso de CPU, memoria en heap, goroutines activas, bloqueos de mutex, latencia de llamadas,
  condiciones de carrera en tiempo de ejecución y fugas de goroutines. Usa pprof, trace, race
  detector y análisis de escape para ver lo invisible.
---

# Go Dynamic Analysis & Profiler — El Código en Movimiento

A diferencia del análisis estático que lee el código sin ejecutarlo, el **análisis dinámico**
observa el comportamiento **real** del programa mientras corre. Revela problemas que solo aparecen
bajo carga, con datos reales o en condiciones de concurrencia.

---

## Las 4 Dimensiones del Análisis Dinámico en Go

### 1. Race Detector — Condiciones de Carrera en Tiempo Real

El detector de carreras de Go instrumenta el binario para detectar accesos concurrentes no sincronizados.

```bash
# Requiere CGO (Linux/Mac). En Windows con CGO habilitado:
CGO_ENABLED=1 go test -race ./...

# Interpretar el output:
# DATA RACE
# Write at 0x00c000... by goroutine 7:    ← quién escribe
# Previous read at 0x00c000... by goroutine 5:  ← quién lee concurrentemente
```

**Qué detecta en timezoner:**
- Escrituras concurrentes en `zoneAliases` (antes de la corrección con `RWMutex`)
- Accesos no sincronizados a cualquier variable de paquete

---

### 2. CPU Profiling — ¿Dónde Gasta Tiempo el Código?

```go
import _ "net/http/pprof"  // En tests/benchmarks:

go test -bench=BenchmarkFindOverlap -cpuprofile cpu.prof .
go tool pprof cpu.prof
# Luego: top10, web (abre grafo de llamas en el navegador)
```

**Interpretar:**
- `flat%`: tiempo que la función pasa en sí misma
- `cum%`: tiempo incluyendo todas las funciones que llama (acumulado)
- Si `zone.LoadLocation` aparece en el top = el caché sync.Map no está funcionando bien

---

### 3. Memory Profiling — ¿Qué Aloca en el Heap?

```bash
go test -bench=BenchmarkConvert -memprofile mem.prof .
go tool pprof mem.prof
# Comandos: top, alloc_space, alloc_objects
```

**Qué buscar en timezoner:**
- `DBTime.Scan` no debe alocar si el valor es `time.Time` directo
- `NewDBTime` debe ser 0 allocs (solo valor en stack)
- `parseString` puede alocar por el slice de layouts (mitigar con variable de paquete)

---

### 4. Escape Analysis — ¿Qué Escapa al Heap?

```bash
go build -gcflags="-m" ./... 2>&1
# Output: "X escapes to heap" indica alocación dinámica
# Output: "X does not escape" indica que permanece en el stack (óptimo)
```

**Reglas de escaping en Go:**
| Causa | Escapa al Heap |
| :--- | :---: |
| Retornar un puntero a variable local | Sí |
| Asignar a interfaz (`any`) | Sí |
| Pasar a función que acepta `interface{}` | Sí |
| Variable cuyo tamaño se desconoce en compilación | Sí |
| Valor pequeño retornado por copia | No |

---

### 5. Goroutine Trace — Visualizar Todas las Goroutines

```bash
go test -bench=BenchmarkFindOverlap -trace trace.out .
go tool trace trace.out
# Abre interfaz web mostrando:
# - Timeline de goroutines
# - Bloqueos de mutex
# - GC pauses
# - Network waits
```

---

### 6. Análisis de Bloqueo (Mutex Contention)

```bash
go test -bench=. -mutexprofile mutex.prof .
go tool pprof mutex.prof
```

**Qué detecta en timezoner:**
- Si muchas goroutines compiten por `aliasesMu.Lock()` → indica que `RegisterAlias` es un cuello de botella
- Si la contención de `locationCache` (sync.Map) es alta → evaluar estrategia de caché alternativa

---

## Checklist de Análisis Dinámico

```
Para cada módulo crítico de timezoner:

pkg/zone:
  □ BenchmarkLoadLocation muestra < 100 ns/op y ≤ 1 alloc/op
  □ BenchmarkLoadLocation_Concurrent con 1000 goroutines sin DATA RACE
  □ go build -gcflags="-m" confirma que *time.Location no escapa al heap repetidamente

pkg/types:
  □ BenchmarkNewDBTime muestra 0 allocs/op
  □ DBTime.Scan(time.Time) muestra 0 allocs/op
  □ ZonedTime.Value() muestra ≤ 2 allocs/op (por el marshal JSON)

pkg/calendar:
  □ BenchmarkAddBusinessDays(10) < 500 ns/op, 0 allocs/op
  □ Testar con DST boundary de New York (8 marzo 2026)

pkg/overlap:
  □ BenchmarkFindOverlap(3 zonas) < 500 µs, verificar alocaciones del map interno
  □ Trace confirma que no hay goroutines spawneadas dentro del algoritmo
```

---

## Comandos de Análisis Dinámico Completo

```bash
# Ejecutar todos los benchmarks con métricas de memoria
go test -bench=. -benchmem -count=3 ./...

# Detección de carreras (con CGO)
CGO_ENABLED=1 go test -race ./...

# Escape analysis completo
go build -gcflags="-m -m" ./... 2>&1 | grep -v "^#"

# Profiling CPU de la fachada
go test -bench=BenchmarkConvert -cpuprofile cpu.prof -count=5 .
go tool pprof -top cpu.prof
```
