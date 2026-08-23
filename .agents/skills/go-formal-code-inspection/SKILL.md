---
name: go-formal-code-inspection
description: >-
  Go Formal Code Inspection Expert (IEEE 1028 / ISO/IEC 20246). Conduce revisiones formales
  estructuradas del código fuente línea por línea usando checklists estandarizados. Analiza
  trazabilidad, flujo de datos, contratos de función, invariantes, precondiciones, postcondiciones
  y todos los caminos de ejecución posibles (happy path + error paths + edge cases).
---

# Formal Code Inspection Expert (IEEE 1028 / ISO/IEC 20246)

La **Inspección Formal de Código** es el proceso más riguroso de revisión de software. A diferencia
de un code review informal, sigue un protocolo estructurado con roles definidos, checklists
estándar y registro de defectos. Fue definida por el estándar **IEEE 1028** y adoptada en
**ISO/IEC 20246**.

---

## Los 5 Roles en una Inspección Formal

| Rol | Responsabilidad |
| :--- | :--- |
| **Moderador** | Dirige la sesión, mantiene el tiempo, asegura cobertura completa |
| **Autor** | Presenta el código, responde preguntas técnicas |
| **Lector** | Lee el código en voz alta, paráfrasis de la lógica |
| **Inspector** | Detecta defectos contra el checklist |
| **Escriba** | Registra todos los defectos encontrados |

---

## Tipos de Inspección

### 1. Inspección de Flujo de Control
Analizar todos los caminos posibles de ejecución de una función:
- **Happy path**: entrada válida → resultado esperado
- **Error paths**: cada condición de error retorna el error correcto
- **Edge cases**: valores límite (vacío, nil, cero, máximo)
- **Dead code**: ramas que nunca pueden ser alcanzadas

**Plantilla:**
```
Función: pkg/ingest.FromString(dateStr, defaultZone string) (time.Time, error)

CAMINOS DE EJECUCIÓN:
  Path 1: dateStr == "" → return ErrEmptyDateString ✓
  Path 2: defaultZone inválido → return ErrInvalidZone ✓
  Path 3: ningún layout parsea → return ErrInvalidInput ✓
  Path 4: RFC3339Nano parsea → return UTC normalizado ✓
  Path 5: "02/01/2006" parsea → return UTC en defaultZone ✓

CAMINOS NO CUBIERTOS POR TESTS:
  Path 6: dateStr con solo espacios (TrimSpace lo neutraliza → Path 1) — verificar
  Path 7: defaultZone con espacios al inicio/fin (TrimSpace lo maneja) — verificar
```

---

### 2. Inspección de Flujo de Datos (Data Flow Analysis)
Rastrear cómo viaja un dato desde que entra hasta que sale:

```
Variable: zoneName en LoadLocation()

DEFINICIÓN:  parámetro de entrada string
→ TrimSpace(zoneName)          → s string (nueva variable)
→ ToUpper(s)                   → upper string
→ zoneAliases[upper]           → canonical string (si existe alias)
→ locationCache.Load(name)     → *time.Location (si hay caché)
→ time.LoadLocation(name)      → *time.Location (I/O potencial)
→ locationCache.Store(name, loc) → EFECTO SECUNDARIO (modifica caché global)
→ return loc, nil              → *time.Location compartido entre goroutines

INVARIANTE CRÍTICA: 'name' después del alias lookup es siempre un ID IANA canónico
PREGUNTA: ¿Qué pasa si zoneAliases["PET"] == "" (string vacío)?
          → time.LoadLocation("") retorna time.UTC (comportamiento no documentado)
```

---

### 3. Inspección de Contratos de Función (Design by Contract)

Para cada función pública verificar:

```
PRECONDICIONES  (lo que debe ser verdad ANTES de llamar):
POSTCONDICIONES (lo que debe ser verdad DESPUÉS de que retorna):
INVARIANTES     (lo que siempre debe ser verdad durante la ejecución):

Ejemplo — NewDBTime(t time.Time) DBTime:
  PRE:  t puede ser cualquier time.Time (incluso zero value)
  POST: resultado.UTC().Location() == time.UTC (siempre en UTC)
  POST: resultado.Nanosecond() == 0 (Round(0) elimina reloj monotónico)
  INV:  nunca modifica el parámetro 't' (value receiver)
```

---

### 4. Inspección de Manejo de Errores

Checklist por función:
- [ ] ¿Todos los retornos de error son revisados internamente?
- [ ] ¿Los errores wrappean la causa raíz con `%w`?
- [ ] ¿Los errores centinela son usados con `errors.Is` en los tests?
- [ ] ¿Los mensajes de error incluyen el contexto suficiente para diagnóstico?
- [ ] ¿Los errores de terceros son traducidos a errores del dominio?

**Ejemplo correcto:**
```go
// Mal: pierde información de la causa
return nil, fmt.Errorf("zona inválida")

// Bien: wrap con %w permite errors.Is() y preserva cadena de error
return nil, fmt.Errorf("%w: zona destino '%s'", ErrInvalidZone, ianaZone)
```

---

### 5. Inspección de Invariantes de Tipos

Para los tipos de dominio (`DBTime`, `ZonedTime`):

```
DBTime — Invariantes:
  □ dbTime.UTC().Location() == time.UTC (siempre)
  □ dbTime.UTC().Nanosecond() nunca tiene residuo monotónico (Round(0))
  □ El método Local() NO EXISTE (por diseño: eliminar la posibilidad de error)
  □ dbTime.String() siempre retorna RFC3339Nano o "" (nunca formato ambiguo)

ZonedTime — Invariantes:
  □ zoned.Zone es siempre un ID IANA canónico (ej: "America/Lima", no "PET")
  □ zoned.UTC.IsZero() == true implica zoned.Zone == "" (representación zero coherente)
  □ zoned.Local() es idempotente: llamarlo n veces produce el mismo resultado
```

---

## Plantilla de Reporte de Inspección Formal

```markdown
## Inspección Formal — [Módulo/Función]
Fecha: YYYY-MM-DD
Inspector: [nombre]
Commit: [hash]

### Defectos Encontrados
| ID | Tipo | Severidad | Línea | Descripción | Estado |
|----|------|-----------|-------|-------------|--------|
| D-01 | Error | CRÍTICO | 45 | Race condition en mapa global | ABIERTO |
| D-02 | Warning | MAYOR | 88 | RFC3339 pierde nanosegundos | CORREGIDO |

### Caminos No Cubiertos por Tests
- [ ] Path 3: zona con alias registrado dinámicamente + concurrencia
- [ ] Path 7: dateStr con caracteres Unicode en TrimSpace

### Preguntas de Diseño Abiertas
- [ ] ¿Debe ZonedTime.Local() retornar un error si la zona ya no existe en la IANA DB?
```
