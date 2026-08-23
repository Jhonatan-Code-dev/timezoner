# 🌍 Timezoner

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Author](https://img.shields.io/badge/Author-Jhonatan-blue.svg)](https://github.com/)
[![License: Proprietary](https://img.shields.io/badge/License-Proprietary-red.svg)](LICENSE)
[![Zero Dependencies](https://img.shields.io/badge/Dependencies-0%20External-brightgreen)](go.mod)
[![Tests](https://img.shields.io/badge/Tests-100%25%20Passing-success)](timezoner_test.go)

**Timezoner** es un paquete y módulo de alto rendimiento en **Golang** creado y mantenido por **Jhonatan**, diseñado para la conversión precisa de husos horarios IANA, cálculo de diferencias temporales, detección de horario de verano (DST) y planificación de horarios coincidentes para equipos internacionales.

> ⚡ **Zero External Dependencies**: Desarrollado 100% sobre la biblioteca estándar de Go, con soporte integrado de caché concurrente en memoria (`sync.Map`) para máxima velocidad y eficiencia.

---

## 📋 Tabla de Contenidos

- [Características](#-características)
- [Instalación](#-instalación)
- [Guía Rápida](#-guía-rápida)
  - [1. Conversión de Zonas y Alias](#1-conversión-de-zonas-y-alias)
  - [2. Fluent API (Estilo Encadenado)](#2-fluent-api-estilo-encadenado)
  - [3. Cálculo de Diferencias y Detección de DST](#3-cálculo-de-diferencias-y-detección-de-dst)
  - [4. Planificador de Solapamiento para Equipos (FindOverlap)](#4-planificador-de-solapamiento-para-equipos-findoverlap)
  - [5. Comparador Múltiple de Zonas (Compare)](#5-comparador-múltiple-de-zonas-compare)
- [Estructura del Proyecto](#-estructura-del-proyecto)
- [Pruebas y Verificación](#-pruebas-y-verificación)
- [Autor y Licencia](#-autor-y-licencia)

---

## ✨ Características

- ⏱️ **Conversión y Parseo Seguro**: Conversión instantánea entre zonas IANA (`"America/Lima"`, `"UTC"`, `"Asia/Tokyo"`, etc.) y alias preconfigurados (`"PET"`, `"EST"`, `"CET"`, `"JST"`, `"COT"`, etc.).
- ⛓️ **API Fluida**: Sintaxis encadenable y expresiva para construir y transformar instancias temporales (`timezoner.Now().In("Europe/Madrid").Format(...)`).
- 📊 **Análisis de Zonas**: Obtención de offsets en segundos y formato legible (`+02:00`, `-05:00`), abreviaturas y cálculo de diferencia con UTC.
- ☀️ **Detección de Horario de Verano (DST)**: Detección automática de si un huso horario tiene activo el horario de verano en un instante dado.
- 🤝 **Overlap Meeting Planner**: Algoritmo para encontrar ventanas de horario laboral coincidentes entre múltiples zonas del mundo.
- 🚀 **Rendimiento Óptimo**: Sin dependencias de terceros y con resolución de `time.Location` optimizada mediante caché.

---

## 📦 Instalación

Para importar **Timezoner** en tu proyecto:

```bash
go get timezoner
```

---

## 🚀 Guía Rápida

### 1. Conversión de Zonas y Alias

```go
package main

import (
	"fmt"
	"time"
	"timezoner"
)

func main() {
	now := time.Now()

	// Convertir a huso horario de Tokio
	tokyoTime, err := timezoner.Convert(now, "Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hora en Tokio: %s\n", tokyoTime.Format("2006-01-02 15:04:05 MST"))

	// Hora actual utilizando alias común (PET = America/Lima)
	limaTime, _ := timezoner.NowIn("PET")
	fmt.Printf("Hora en Lima: %s\n", limaTime.Format("15:04:05"))
}
```

---

### 2. Fluent API (Estilo Encadenado)

```go
formatted, err := timezoner.Now().
	In("Europe/Paris").
	Format("2006-01-02 15:04:05 MST")

if err != nil {
	panic(err)
}

fmt.Printf("Hora en París: %s\n", formatted)
```

---

### 3. Cálculo de Diferencias y Detección de DST

```go
targetDate := time.Date(2026, 7, 15, 12, 0, 0, 0, time.UTC)

// Calcular diferencia horaria exacta
diff, _ := timezoner.Difference("Europe/Madrid", "America/Lima", targetDate)
fmt.Printf("Diferencia horaria: %v horas\n", diff.Hours()) // +7.0 horas

// Comprobar estado de Horario de Verano (DST)
isDST, _ := timezoner.IsDST("Europe/Madrid", targetDate)
fmt.Printf("¿Madrid en horario de verano?: %v\n", isDST) // true en Julio
```

---

### 4. Planificador de Solapamiento para Equipos (`FindOverlap`)

Encuentra ventanas de tiempo hábiles comunes para coordinar reuniones entre participantes de distintos países:

```go
package main

import (
	"fmt"
	"time"
	"timezoner"
)

func main() {
	slots, err := timezoner.FindOverlap(timezoner.OverlapRequest{
		Date:         time.Date(2026, 10, 15, 0, 0, 0, 0, time.UTC),
		Zones:        []string{"America/Lima", "America/New_York", "Europe/Madrid"},
		DefaultHours: timezoner.WorkingHours{StartHour: 9, EndHour: 18}, // 09:00 a 18:00
		SlotDuration: 1 * time.Hour,
	})

	if err != nil {
		panic(err)
	}

	for i, slot := range slots {
		fmt.Printf("Ventana #%d (Duración: %v):\n", i+1, slot.Duration)
		fmt.Printf("  • UTC:         %s - %s\n", slot.StartTimeUTC.Format("15:04"), slot.EndTimeUTC.Format("15:04"))
		fmt.Printf("  • Lima:        %s - %s\n", slot.ZoneTimes["America/Lima"].Format("15:04"), slot.ZoneTimes["America/Lima"].Add(slot.Duration).Format("15:04"))
		fmt.Printf("  • Nueva York:  %s - %s\n", slot.ZoneTimes["America/New_York"].Format("15:04"), slot.ZoneTimes["America/New_York"].Add(slot.Duration).Format("15:04"))
		fmt.Printf("  • Madrid:      %s - %s\n", slot.ZoneTimes["Europe/Madrid"].Format("15:04"), slot.ZoneTimes["Europe/Madrid"].Add(slot.Duration).Format("15:04"))
	}
}
```

---

### 5. Comparador Múltiple de Zonas (`Compare`)

```go
snapshots, err := timezoner.Compare(time.Now(), "America/Lima", "Europe/London", "Asia/Tokyo")
if err != nil {
	panic(err)
}

for _, s := range snapshots {
	fmt.Printf("%-22s | %s | Offset: %s | DST: %v\n", s.Zone, s.Formatted, s.OffsetFormatted, s.IsDST)
}
```

---

## 📂 Estructura del Proyecto

```
timezoner/
├── go.mod                      # Módulo Go puro (0 dependencias externas)
├── timezoner.go                # Funciones centrales y Fluent API
├── zones.go                    # Mapeo de alias, validación y caché de zonas IANA
├── timezoner_test.go           # Pruebas unitarias completas y benchmarks
├── examples_test.go            # Ejemplos ejecutables para godoc / pkg.go.dev
├── examples/                   # Casos de uso prácticos
│   ├── basic_usage/main.go     # Ejemplo de uso de funciones básicas
│   └── team_meeting_planner/main.go # Ejemplo del planificador de reuniones
├── LICENSE                     # Licencia Propietaria Exclusiva
└── README.md                   # Documentación técnica
```

---

## 🧪 Pruebas y Verificación

Para ejecutar la suite de pruebas unitarias y benchmarks:

```bash
# Ejecutar todas las pruebas con reporte de cobertura
go test -v -cover ./...

# Ejecutar benchmarks de rendimiento
go test -bench=. -benchmem ./...
```

---

## 👤 Autor y Licencia

- **Creador y Autor**: **Jhonatan**
- **Licencia**: **Licencia Propietaria / Todos los derechos reservados**.

> ⚠️ **Aviso de Propiedad Intelectual**: Este software es propiedad exclusiva de **Jhonatan**. Queda estrictamente prohibida la copia, distribución, modificación, sublicenciamiento, alteración o comercialización de este código fuente sin la autorización expresa y por escrito de su autor. Consulta el archivo [`LICENSE`](LICENSE) para más detalles.
