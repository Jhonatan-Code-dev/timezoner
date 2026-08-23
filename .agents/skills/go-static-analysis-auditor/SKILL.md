---
name: go-static-analysis-auditor
description: >-
  Go Static Analysis Auditor. Ejecuta y interpreta el análisis estático completo del código fuente
  sin ejecutarlo: go vet, staticcheck, golangci-lint, govulncheck. Detecta bugs, vulnerabilidades,
  código muerto, complejidad ciclomática y anti-patrones idiomáticos de Go antes de que exploten en producción.
---

# Go Static Analysis Auditor — Experto en Análisis Estático

Este skill actúa como auditor de análisis estático de nivel senior. Examina el código fuente sin
ejecutarlo para detectar defectos estructurales, vulnerabilidades de seguridad, código muerto y
anti-patrones que el compilador de Go no reporta pero que causan bugs en producción.

---

## ¿Qué es el Análisis Estático?

El análisis estático examina el código **sin ejecutarlo**. Detecta lo que nosotros no vemos a
simple vista: condiciones de carrera potenciales, punteros nulos garantizados, flujos de error
ignorados, variables no usadas y violaciones de las convenciones idiomáticas de Go.

---

## Herramientas Obligatorias (en orden de severidad)

### 1. `go vet` — Análisis Nativo del Compilador
```bash
go vet ./...
```
Detecta:
- Llamadas a `fmt.Printf` con formato incorrecto (`printf("value: %d", "string")`)
- Locks copiados por valor (`sync.Mutex` no debe pasarse por valor)
- Comparaciones de structs que contienen campos no comparables
- Uso incorrecto de `defer` dentro de bucles

### 2. `staticcheck` — Análisis Avanzado (SA, S, QF, ST)
```bash
go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck ./...
```
Detecta:
- **SA1019**: Uso de APIs deprecadas
- **SA4009**: Argumento de función que siempre es sobreescrito antes de usarse
- **SA5011**: Posible desreferenciación de puntero nil
- **S1028**: `errors.New(fmt.Sprintf(...))` → debe ser `fmt.Errorf`
- **ST1003**: Nombres de exportados que no siguen convenciones de Go

### 3. `golangci-lint` — Meta-Analizador (50+ linters)
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint run ./...
```
Incluye automáticamente:
- `errcheck`: verifica que todos los errores sean manejados
- `gosimple`: simplificaciones de código
- `ineffassign`: asignaciones que nunca se usan
- `typecheck`: errores de tipo
- `unused`: código no utilizado (funciones, variables, tipos)
- `gocritic`: detecta anti-patrones y code smells
- `revive`: sucesor de `golint` con reglas configurables

### 4. `govulncheck` — Vulnerabilidades en Dependencias
```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```
Analiza si alguna dependencia (directa o transitiva) tiene CVEs publicados que afecten al código
que realmente se ejecuta (no solo las que están en go.mod).

---

## Checklist de Análisis Estático para Go

### Flujo de Errores
- [ ] Todos los valores de retorno de error son manejados (`errcheck`)
- [ ] No hay `errors.New(fmt.Sprintf(...))` — usar `fmt.Errorf`
- [ ] No hay `_ = functionThatReturnsError()`  en código de producción
- [ ] Los errores centinela usan `errors.Is` / `errors.As`, nunca comparación directa `==`

### Gestión de Memoria
- [ ] No hay slices globales exportados mutables
- [ ] No hay maps globales sin sincronización (`sync.Map` o `sync.RWMutex`)
- [ ] Las funciones que crecen en el heap son identificadas y optimizadas
- [ ] No hay conversiones `interface{}` / `any` innecesarias que escapan al heap

### Concurrencia
- [ ] `sync.Mutex` y `sync.RWMutex` nunca se pasan por valor
- [ ] No hay goroutines sin lifetime acotado (goroutine leaks)
- [ ] No hay acceso concurrente a variables sin sincronización

### API y Nombres
- [ ] Todos los exportados tienen comentario Godoc
- [ ] Nombres siguen convenciones (`MixedCaps`, no `snake_case`)
- [ ] No hay parámetros de función que nunca se leen antes de sobrescribirse
- [ ] No hay código muerto (funciones privadas no referenciadas)

---

## Cómo Aplicar Este Skill

Cuando se active este skill, ejecutar en orden:

```bash
# Paso 1: Análisis nativo
go vet ./...

# Paso 2: Análisis avanzado
staticcheck ./...

# Paso 3: Meta-analizador completo
golangci-lint run ./...

# Paso 4: Vulnerabilidades
govulncheck ./...
```

Reportar cada hallazgo con:
- Archivo y número de línea
- Categoría del defecto (Error, Warning, Info)
- Descripción del impacto en producción
- Corrección propuesta con código
