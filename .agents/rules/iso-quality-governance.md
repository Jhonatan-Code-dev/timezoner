# Gobernanza de Calidad y Estándares ISO/IEC

Este archivo define las reglas de calidad que todo agente y desarrollador debe verificar y cumplir en este proyecto según las normas internacionales de ingeniería de software:

---

## 1. ISO/IEC 25010 (SQuaRE — Modelo de Calidad del Producto)
- **Portabilidad Universal**: El módulo debe funcionar en cualquier sistema operativo (Windows, Linux, Alpine, Scratch) sin requerir paquetes `tzdata` del sistema, mediante el uso de `_ "time/tzdata"`.
- **Confiabilidad e Inmunidad a Pánicos**: Ninguna función pública o método de conversión debe producir un pánico ante entradas inválidas, nulas o corruptas de los usuarios. Todo error debe retornarse estructurado (`error`).
- **Mantenibilidad y Modularidad**: Complejidad ciclomática baja ($\le 10$ por función). Código organizado por dominios desacoplados en `pkg/`.
- **Compatibilidad con Ecosistema SQL/JSON**: Todo tipo de persistencia (`DBTime`, `ZonedTime`) debe implementar fielmente `driver.Valuer`, `sql.Scanner`, `json.Marshaler` y `json.Unmarshaler`.

## 2. ISO/IEC 5055 (Medición de Calidad en Código Fuente)
- **Prevención de CWE-400 (Consumo No Controlado de Recursos)**: Los algoritmos de búsqueda (como `FindOverlap`) deben tener pasos temporales discretos y límites finitos de escaneo (24 horas).
- **Prevención de CWE-476 (Desreferenciación de Punteros Nulos)**: Comprobación defensiva `if value == nil` y `if tp.err != nil` en todas las entradas y métodos `Scan`.
- **Prevención de CWE-190 (Desbordamiento Aritmético)**: Toda aritmética temporal debe usar tipos seguros de 64 bits (`time.Duration`).

## 3. ISO/IEC 29119 (Estándares de Pruebas de Software)
- **Cobertura Mínima de Pruebas**: Cobertura obligatoria $\ge 80\%$ en todos los paquetes (`pkg/*` y raíz).
- **Análisis de Valores Límite (Boundary Value Analysis)**:
  - Validar límites de años bisiestos (29 de febrero).
  - Validar transiciones de horario de verano (DST spring-forward y fall-back).
  - Validar husos horarios fraccionarios (ej. UTC+05:30 India, UTC+05:45 Nepal).
  - Validar aritmética de días hábiles en viernes, domingos y fechas límites.
- **Fuzzing y Estrés Concurrente**: Todo parser de cadenas debe contar con pruebas de fuzzing nativas (`go test -fuzz`) y pruebas de estrés concurrente con múltiples goroutines.

## 4. ISO 8601 / RFC 3339 (Representación de Fecha y Hora)
- Todas las cadenas de intercambio de fecha y hora deben cumplir estrictamente con el formato universal ISO 8601 / RFC 3339 (`YYYY-MM-DDTHH:MM:SSZ` o `RFC3339Nano`).
