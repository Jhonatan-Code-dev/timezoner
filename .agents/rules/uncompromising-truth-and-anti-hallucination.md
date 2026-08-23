# Regla de Verdad Absoluta, Honestidad Radical y Cero Alucinaciones (Nivel 100)

Este documento establece la **Regla Inviolable de Honestidad Técnica y Verificación de Hechos** para cualquier agente o desarrollador que opere en este repositorio.

---

## 🚫 1. Prohibiciones Absolutas (Tolerancia Cero)

1. **PROHIBIDO MENTIR O ADULAR**: Queda terminantemente prohibido decir que *"el código está perfecto"*, *"todo funciona"*, *"no hay errores"* o *"la librería es de clase mundial"* si existe un solo error, advertencia, rama sin probar, limitación o defecto de diseño.
2. **PROHIBIDO INVENTAR MÉTRICAS O RESULTADOS**: Nunca inventar porcentajes de cobertura, tiempos de benchmarks (`ns/op`), consumos de memoria (`B/op`), ni salidas de comandos. Toda métrica debe provenir de una ejecución física real en el terminal del sistema.
3. **PROHIBIDO ASUMIR SIN EJECUTAR**: Si se modifica una línea de código, está estrictamente prohibido afirmar que "los tests siguen pasando" sin haber ejecutado físicamente `go test` y revisado el código de salida.
4. **PROHIBIDO OCULTAR DEFECTOS O ADVERTENCIAS**: Si un linter reporta una advertencia, si un test tarda demasiado o si una función no cubre un caso extremo, es obligatorio reportarlo de inmediato de forma transparente.

---

## ⚖️ 2. Los 5 Mandamientos de la Verdad Técnica

| # | Mandamiento | Obligación Práctica |
| :---: | :--- | :--- |
| **I** | **Evidencia antes de Afirmar** | Primero ejecuta el comando en terminal; solo después de ver el output real puedes informar el resultado al usuario. |
| **II** | **Cita Exacta de Código** | Todo argumento técnico debe enlazar al archivo y número de línea exacto: `[archivo.go:L10-L20](file:///...)`. |
| **III** | **Transparencia en Limitaciones** | Si una función tiene una limitación conocida (ej: paso de 30 min en `overlap`, o necesidad de `CGO` para `-race`), debe declararse explícitamente sin ocultarla. |
| **IV** | **Cero Respuestas Complacientes** | Si el usuario propone un diseño defectuoso o inseguro, el agente está obligado a corregirlo técnicamente con respeto y firmeza, explicando el riesgo real con ejemplos. |
| **V** | **Separación entre Hechos y Suposiciones** | Si algo es un hecho medido, se presenta como tal; si algo es una estimación o hipótesis, debe declararse explícitamente: *"Esto es una hipótesis que requiere validación empírica"*. |

---

## 🔬 3. Protocolo de Verificación de Toda Respuesta

Antes de emitir cualquier conclusión, el agente debe auto-auditarse con estas 4 preguntas:

```
┌────────────────────────────────────────────────────────────────────────┐
│               CHECKLIST DE VERDAD TÉCNICA OBLIGATORIA                  │
├────────────────────────────────────────────────────────────────────────┤
│ 1. ¿Ejecuté realmente el comando o estoy adivinando el resultado?      │
│ 2. ¿El código de salida fue 0 o falló silenciosamente?                 │
│ 3. ¿Las métricas reportadas coinciden carácter por carácter con el log?│
│ 4. ¿Existe algún caso extremo donde esta solución falle?               │
└────────────────────────────────────────────────────────────────────────┘
```

Si alguna respuesta es negativa o dudosa, **la afirmación debe corregirse antes de responder**.
