---
name: go-side-effect-inspector
description: >-
  Go Side Effect Inspector. Analiza TODO lo que hace una función o módulo más allá de su valor
  de retorno: mutaciones de estado global, alocaciones de heap, llamadas externas, modificaciones
  por referencia, efectos sobre parámetros, y comportamiento ante entradas inesperadas. Detecta
  lo que el código hace pero que NO se ve en la firma de la función.
---

# Go Side Effect Inspector — Análisis de Efectos Secundarios

Este skill es el más profundo a nivel de comportamiento de funciones. Un **efecto secundario**
(side effect) es cualquier cosa que una función hace que **no está expresada en su tipo de retorno**.
Son los comportamientos invisibles que causan los bugs más difíciles de diagnosticar.

---

## ¿Qué son los Efectos Secundarios?

```go
// Esta función parece inocente:
func LoadLocation(zoneName string) (*time.Location, error) { ... }

// Pero tiene efectos secundarios ocultos:
// 1. ESCRIBE en un sync.Map global (locationCache)
// 2. LEE de un map global (zoneAliases) — potencial race condition
// 3. ALOCA memoria en el heap para *time.Location
// 4. PUEDE hacer I/O para leer tzdata del sistema operativo
// 5. El resultado es COMPARTIDO entre goroutines (mismo puntero cacheado)
```

---

## Categorías de Efectos Secundarios a Inspeccionar

### 1. Mutaciones de Estado Global
```
PREGUNTA: ¿La función modifica alguna variable de paquete o global?
SEÑALES:  Escrituras en sync.Map, maps globales, slices globales, variables var
RIESGO:   Comportamiento no determinista entre llamadas, condiciones de carrera
```
**Cómo detectarlo:**
- Buscar `var ` a nivel de paquete que sea modificable
- Buscar `.Store(`, `.Lock()`, `.Set(`, `append(global,`
- Buscar `[key] =` donde el mapa sea de paquete

---

### 2. Alocaciones de Heap (Escape Analysis)
```
PREGUNTA: ¿La función crea objetos que escapan al heap?
SEÑALES:  Retornar punteros, interfaces, slices, o cualquier tipo que Go "escape"
RIESGO:   Presión de GC, latencia impredecible, degradación bajo carga
```
**Cómo detectarlo:**
```bash
go build -gcflags="-m -m" ./... 2>&1 | grep "escapes to heap"
go test -bench=. -benchmem ./...   # allocs/op > 0 indica heap allocations
```

---

### 3. Modificaciones de Parámetros por Referencia
```
PREGUNTA: ¿La función modifica el contenido de sus parámetros (slice, map, pointer)?
SEÑALES:  Funciones que reciben *T, []T, map[K]V y escriben en ellos
RIESGO:   El llamador ve sus datos modificados sin esperarlo (violación de contratos)
```
**Ejemplo peligroso:**
```go
// El caller no sabe que su slice fue modificado:
func processZones(zones []string) {
    zones[0] = "UTC" // EFECTO SECUNDARIO: modifica el slice original del llamador
}
```

---

### 4. Comportamiento Dependiente de Tiempo Real (`time.Now()`)
```
PREGUNTA: ¿La función llama time.Now() internamente sin que el llamador lo sepa?
SEÑALES:  NowDBTime(), ingest.Now(), humanize.Humanize() sin relativeTo
RIESGO:   Tests no deterministas, comportamiento diferente en UTC vs hora local del servidor
```
**Solución:** Inyección de tiempo vía parámetro opcional:
```go
func Humanize(t time.Time, relativeTo ...time.Time) string {
    now := time.Now()
    if len(relativeTo) > 0 { now = relativeTo[0] } // inyectable para tests
    ...
}
```

---

### 5. Efectos sobre el Sistema de Archivos o Red
```
PREGUNTA: ¿La función lee archivos del sistema, hace llamadas HTTP, o accede a DNS?
SEÑALES:  time.LoadLocation() puede leer /usr/share/zoneinfo en sistemas sin tzdata embebido
RIESGO:   Fallo en contenedores scratch o en ambientes sin permisos de lectura
MITIGACIÓN: _ "time/tzdata" embebe los datos en el binario eliminando la dependencia de I/O
```

---

## Plantilla de Inspección por Función

Para cada función pública analizar:

```
Función: [nombre]
Firma:   [parámetros → retornos]

EFECTOS SECUNDARIOS DETECTADOS:
  □ Estado global: [describe qué variables de paquete modifica o lee]
  □ Alocaciones heap: [allocs/op del benchmark]
  □ Mutación de parámetros: [¿modifica algún parámetro de entrada?]
  □ Dependencia de tiempo real: [¿llama time.Now() internamente?]
  □ I/O externo: [¿lee archivos, red, o recursos del OS?]
  □ Pánico potencial: [¿puede entrar en pánico con ciertos inputs?]
  □ Goroutines: [¿inicia alguna goroutine sin el conocimiento del llamador?]

NIVEL DE RIESGO: [CRÍTICO / MAYOR / MENOR / ACEPTABLE]
ACCIÓN:          [Corrección propuesta]
```

---

## Aplicación en timezoner

Ejecutar esta inspección en cada función de:
- `pkg/zone/zone.go` — `LoadLocation`, `RegisterAlias`, `GetInfo`
- `pkg/types/dbtime.go` — `Scan`, `MarshalJSON`, `Value`
- `pkg/types/zonedtime.go` — `unmarshalString`, `Local`, `Scan`
- `pkg/ingest/ingest.go` — `FromString`, `FromLocal`, `Now`
- `pkg/calendar/calendar.go` — `AddBusinessDays` (efectos de DST)
- `pkg/overlap/overlap.go` — `Find` (goroutines, alocaciones de maps)
- `timezoner.go` — toda la fachada pública (efectos heredados)
