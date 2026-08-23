# Mandato de Evidencia Empírica Obligatoria (Nivel 100 — Tolerancia Cero)

Este documento establece la **Regla Suprema de Verificación** para todos los agentes, desarrolladores y procesos automatizados en el repositorio **timezoner**. 

---

## 🚫 Principio Fundamental: Prohibición Absoluta de Afirmaciones sin Evidencia

> **"Una afirmación sin evidencia empírica verificable se considera falsa por defecto."**

Queda **estrictamente prohibido**:
1. Decir *"todo está bien"*, *"es seguro"*, *"el rendimiento es óptimo"*, *"cumple con los estándares"* sin adjuntar la prueba física ejecutable.
2. Reportar porcentajes de cobertura inventados, aproximados o inferidos.
3. Declarar que un bug o condición de carrera está "resuelto" sin mostrar el test que fallaba antes y pasa ahora, junto con el enlace al código exacto (`file:///...#Lxx-Lyy`).
4. Alabar la velocidad o eficiencia sin adjuntar la salida textual de `go test -bench=. -benchmem`.

---

## 📋 Estructura Obligatoria de Todo Informe o Respuesta Técnica

Cada vez que se afirme un estado, corrección o nivel de calidad, la respuesta DEBE incluir obligatoriamente las **4 Pruebas de Evidencia**:

```
┌────────────────────────────────────────────────────────────────────────┐
│                   PROTOCOLO DE EVIDENCIA EMPÍRICA                      │
├────────────────────────────────────────────────────────────────────────┤
│ 1. COMANDO EJECUTADO: Comando exacto del sistema (ej: go test -v)      │
│ 2. SALIDA CRUDA DEL SISTEMA: Stdout/Stderr textual con métricas reales │
│ 3. ENLACE A LÍNEAS EXACTAS: [archivo.go:L10-L25](file:///...)          │
│ 4. DEMOSTRACIÓN TÉCNICA: Por qué el código garantiza el comportamiento │
└────────────────────────────────────────────────────────────────────────┘
```

---

## 🔬 Matriz de Evidencia por Tipo de Afirmación

| Si afirmas que... | La EVIDENCIA OBLIGATORIA que debes presentar es... |
| :--- | :--- |
| **"Los tests pasan"** | El log literal de `go test ./...` con estado `ok`, tiempo en segundos y código de salida `0`. |
| **"La cobertura es del X%"** | El output directo de `go test -cover ./...` por cada submódulo individual en `pkg/`. |
| **"No hay condiciones de carrera"** | El log de pruebas concurrentes con múltiples goroutines y el enlace a la sincronización (`sync.RWMutex` / `sync.Map`). |
| **"Es de alto rendimiento / 0 allocs"** | La tabla de `go test -bench=. -benchmem` mostrando `ns/op`, `B/op` y `allocs/op` exactos. |
| **"El código es inmutable"** | El enlace al código fuente mostrando el tipo privado y la función que retorna una copia defensiva (`copy()`). |
| **"Maneja DST / bisiestos"** | El test específico en `calendar_test.go` que evalúa la fecha exacta (ej: 29 Feb 2028 o 8 Mar 2026 en New York). |
| **"No hay errores de compilación"** | El log de `go vet ./...` y `go build ./...` con código de salida `0`. |

---

## ⚖️ Consecuencia de Incumplimiento
Cualquier auditoría, refactorización o cambio que no cumpla con este protocolo será rechazado de inmediato y catalogado como **NO VERIFICADO / NO CERTIFICABLE**.
