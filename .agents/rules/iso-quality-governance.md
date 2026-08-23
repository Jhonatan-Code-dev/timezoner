# Gobernanza de Calidad ISO/IEC (Nivel 100 — Tolerancia Cero)

Este archivo define los umbrales cuantitativos innegociables de calidad bajo normas internacionales de ingeniería de software:

---

## 1. ISO/IEC 25010 (SQuaRE — Modelo de Calidad de Producto)
- **Portabilidad**: Prohibido depender de `tzdata` del sistema operativo. La base de datos IANA debe estar obligatoriamente embebida mediante `_ "time/tzdata"`.
- **Confiabilidad**: Verificación de punteros nulos (`nil`) obligatoria en todos los métodos de serialización (`Scan`, `UnmarshalJSON`, `Value`).
- **Mantenibilidad**: Cohesión de dominio en `pkg/`. Complejidad ciclomática por función $\le 10$.

## 2. ISO/IEC 5055 (Defectos Estructurales en Código Fuente)
- **CWE-362 (Condición de Carrera)**: Cero accesos concurrentes no sincronizados.
- **CWE-400 (Consumo No Controlado de Recursos)**: Algoritmos con límites acotados de iteración.
- **CWE-476 (Puntero Nulo)**: Validación defensiva en toda frontera de paquete.
- **CWE-190 (Desbordamiento de Enteros)**: Aritmética basada en tipos seguros de 64 bits (`time.Duration`).
- **Acoplamiento Circular**: Cero ciclos en el grafo de paquetes (`go vet` / `go build`).

## 3. ISO/IEC 29119 (Testing y Verificación Dinámica)
- **Umbral de Cobertura**: Mínimo **80% de cobertura de sentencias** en cada paquete de dominio en `pkg/`.
- **Análisis de Valores Límite (Boundary Value Analysis)**:
  - Pruebas explícitas de 29 de febrero en años bisiestos (`DaysInMonth(2028, February) == 29`).
  - Pruebas explícitas de transición de horario de verano DST (`AddBusinessDays` cruzando el domingo de cambio de hora).
  - Pruebas explícitas de husos con minutos fraccionarios (`Asia/Kolkata` +05:30, `Asia/Kathmandu` +05:45).
- **Fuzzing y Concurrencia**:
  - `FuzzConvert` y `FuzzIngestFromString` nativos ejecutados sin pánicos.
  - Pruebas de estrés con $\ge 100$ goroutines simultáneas.

## 4. ISO 8601 / RFC 3339
- Toda serialización de fecha debe emitir formato canónico RFC 3339 o RFC 3339 Nano (`YYYY-MM-DDTHH:MM:SS.NNNNNNNNNZ07:00`), garantizando la preservación completa de la precisión temporal.
