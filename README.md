# Timezoner

Timezoner es un paquete de alto rendimiento en Go puro (cero dependencias externas) diseñado para la conversión precisa de husos horarios IANA, cálculo de diferencias temporales, detección de horario de verano (DST), aritmética de calendario laboral, formateo de tiempo relativo humano y planificación de reuniones para equipos distribuidos globalmente.

Proporciona soporte de base de datos IANA embebida (`time/tzdata`) para portabilidad absoluta, integración nativa con `database/sql` y dos patrones empresariales de persistencia en base de datos (`DBTime` y `ZonedTime`).

---

## Los 2 Patrones Empresariales de Persistencia

Timezoner implementa de forma nativa los dos estándares de la industria para almacenamiento de fechas:

| Patrón | Tipo | Caso de Uso Principal | Estrategia de Almacenamiento |
| :--- | :--- | :--- | :--- |
| **Patrón 1: Transaccional / Auditoría** | `timezoner.DBTime` | Pagos, pedidos, logs, mensajes de chat, eventos históricos. | Columna única almacenada en UTC absoluto (`TIMESTAMPTZ`). |
| **Patrón 2: Citas Futuras y Calendarios** | `timezoner.ZonedTime` | Citas médicas, vuelos, webinars, alarmas, recordatorios. | Dos columnas (o JSONB) guardando el instante UTC + la zona IANA de origen (`"America/Lima"`). Preserva la hora local original ante cambios gubernamentales de DST. |

---

## Capacidades Principales

- **Portabilidad Universal (`time/tzdata`)**: Incluye la base de datos IANA oficial dentro del binario. Funciona en Windows, Alpine Linux, AWS Lambda y contenedores Scratch de Docker sin requerir paquetes `tzdata` del sistema operativo.
- **Tipos Nativos para SQL y JSON (`DBTime`, `ZonedTime`)**: Implementan `driver.Valuer`, `sql.Scanner`, `json.Marshaler` y `json.Unmarshaler` garantizando almacenamiento limpio en UTC sin distorsiones de reloj monotónico.
- **Aritmética de Calendario y Días Hábiles**: Añade o resta días laborables (`AddBusinessDays`) ignorando fines de semana, y calcula límites de fecha (`StartOfDay`, `EndOfMonth`).
- **Tiempo Relativo Humano (`Humanize`)**: Convierte diferencias de tiempo en lenguaje natural en español e inglés (`"hace 5 minutos"`, `"en 2 horas"`).
- **Ingesta y Proyección para Base de Datos**: Pipelines completos para normalizar entradas locales a UTC y proyectar registros UTC a la zona de cualquier usuario.
- **Planificador de Solapamiento para Equipos**: Calcula intervalos de trabajo coincidentes entre participantes de múltiples países para cualquier fecha.
- **Caché Concurrente en Memoria**: Optimización de resolución de `*time.Location` mediante `sync.Map` para máxima velocidad.
- **Cero Dependencias Externas**: Desarrollado 100% sobre la biblioteca estándar de Go (`time`, `sync`, `errors`, `database/sql/driver`, `fmt`).

---

## Instalación

```bash
go get timezoner
```

Requiere Go 1.22 o superior.

---

## Ejemplos de Uso

### 1. Patrón 1: Pagos y Transacciones (UTC Puro con `DBTime`)

```go
package main

import (
	"fmt"
	"timezoner"
)

type Payment struct {
	ID        string           `json:"id"`
	Amount    float64          `json:"amount"`
	PaidAtUTC timezoner.DBTime `json:"paid_at_utc"`
}

func main() {
	p := Payment{
		ID:        "TX-1001",
		Amount:    250.00,
		PaidAtUTC: timezoner.NowDBTime(),
	}

	// Proyectar para un usuario en Madrid
	madridView, _ := timezoner.ProjectForUser(p.PaidAtUTC.Time, "Europe/Madrid")
	fmt.Printf("El usuario en Madrid ve: %s (%s)\n", madridView.Formatted, madridView.OffsetFormatted)
}
```

---

### 2. Patrón 2: Citas Futuras y Calendarios (`ZonedTime`)

```go
package main

import (
	"fmt"
	"timezoner"
)

type Appointment struct {
	ID          int                 `json:"id"`
	Doctor      string              `json:"doctor"`
	ScheduledAt timezoner.ZonedTime `json:"scheduled_at"`
}

func main() {
	// Un paciente en Lima agenda cita para las 10:00 AM
	zoned, err := timezoner.ZonedFromLocal("2026-10-01 10:00", "America/Lima")
	if err != nil {
		panic(err)
	}

	app := Appointment{
		ID:          1,
		Doctor:      "Dra. García",
		ScheduledAt: zoned,
	}

	// 1. Hora local garantizada en Lima
	localTime, _ := app.ScheduledAt.Local()
	fmt.Println("Lima:", localTime.Format("2006-01-02 15:04")) // 2026-10-01 10:00

	// 2. Proyectada para un especialista en Tokio
	tokyoView, _ := app.ScheduledAt.ForViewer("Asia/Tokyo")
	fmt.Printf("Tokio: %s (%s)\n", tokyoView.Formatted, tokyoView.OffsetFormatted)
}
```

---

### 3. Días Hábiles y Límites de Calendario

```go
package main

import (
	"fmt"
	"time"
	"timezoner"
)

func main() {
	// Partiendo de un Viernes
	viernes := time.Date(2026, 9, 4, 10, 0, 0, 0, time.UTC)

	// Sumar 5 días hábiles (salta sábado y domingo -> siguiente viernes) y mover al final del día
	fechaVencimiento := timezoner.At(viernes).
		AddBusinessDays(5).
		EndOfDay().
		MustTime()

	fmt.Println("Fecha de vencimiento:", fechaVencimiento.Format("2006-01-02 15:04:05 MST"))
	// Salida: Fecha de vencimiento: 2026-09-11 23:59:59 UTC
}
```

---

### 4. Tiempo Relativo Humano (Humanize)

```go
package main

import (
	"fmt"
	"time"
	"timezoner"
)

func main() {
	t := time.Now().Add(-2 * time.Hour)

	fmt.Println("Español:", timezoner.Humanize(t))   // "hace 2 horas"
	fmt.Println("Inglés:", timezoner.HumanizeEn(t))  // "2 hours ago"
}
```

---

### 5. Planificador de Solapamiento para Equipos Globales

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
		DefaultHours: timezoner.WorkingHours{StartHour: 9, EndHour: 18}, // 09:00 - 18:00
		SlotDuration: 1 * time.Hour,
	})
	if err != nil {
		panic(err)
	}

	for i, slot := range slots {
		fmt.Printf("Ventana #%d (Duración: %v):\n", i+1, slot.Duration)
		fmt.Printf("  UTC:         %s - %s\n", slot.StartTimeUTC.Format("15:04"), slot.EndTimeUTC.Format("15:04"))
		fmt.Printf("  Lima:        %s\n", slot.ZoneTimes["America/Lima"].Format("15:04"))
		fmt.Printf("  Nueva York:  %s\n", slot.ZoneTimes["America/New_York"].Format("15:04"))
		fmt.Printf("  Madrid:      %s\n", slot.ZoneTimes["Europe/Madrid"].Format("15:04"))
	}
}
```

---

## Arquitectura de Paquetes y Tipos

| Tipo / Paquete | Propósito |
| :--- | :--- |
| **`timezoner.DBTime`** | Tipo SQL/JSON de columna única en UTC para transacciones, pagos y auditoría. |
| **`timezoner.ZonedTime`** | Tipo SQL/JSON de doble columna (UTC + zona IANA) para citas, calendarios y vuelos. |
| **`timezoner.Ingest*`** | Sanitiza y normaliza fechas locales de usuarios a UTC antes de persistir en BD. |
| **`timezoner.Project*`** | Proyecta registros UTC de la BD a la zona específica de cada usuario con ISO 8601 y estado DST. |

---

## Modelo de Errores

Timezoner exporta errores centinela tipados para verificación explícita con `errors.Is`:

- `timezoner.ErrEmptyZoneName`: Retornado cuando se ingresa una cadena vacía como identificador de zona.
- `timezoner.ErrInvalidZone`: Retornado cuando el nombre de la zona no existe en la base de datos IANA.
- `timezoner.ErrInvalidTimeFormat`: Retornado cuando una cadena de fecha no coincide con el layout provisto.
- `timezoner.ErrNoZonesProvided`: Retornado cuando se ejecuta una operación multizona sin proporcionar zonas.
- `timezoner.ErrEmptyDateString`: Retornado cuando se intenta procesar un texto de fecha vacío.

---

## Pruebas y Benchmarks

Ejecutar la suite completa de pruebas unitarias con cobertura:

```bash
go test -v -cover ./...
```

Ejecutar benchmarks de rendimiento y memoria:

```bash
go test -bench=. -benchmem ./...
```

---

## Estructura del Repositorio

```
.
├── go.mod                      # Definición del módulo Go
├── timezoner.go                # API central y Fluent API
├── zones.go                    # Catálogo IANA, caché en memoria, tzdata embebido y alias
├── calendar.go                 # Días laborables, fines de semana y límites de fecha
├── humanize.go                 # Tiempo relativo en lenguaje natural (ES/EN)
├── dbtime.go                   # Driver SQL Valuer / Scanner en UTC puro (Patrón 1)
├── zonedtime.go                # Tipo SQL/JSON con UTC + Zona IANA para calendarios (Patrón 2)
├── ingest.go                   # Ingesta y normalización a UTC para BD
├── project.go                  # Proyección y formateo para usuarios finales
├── timezoner_test.go           # Tests unitarios centrales, fuzzing y pruebas de concurrencia
├── calendar_test.go            # Tests unitarios de calendario y días hábiles
├── humanize_test.go            # Tests de tiempo relativo
├── dbtime_test.go              # Tests de persistencia SQL y JSON
├── zonedtime_test.go           # Tests de ZonedTime para citas y eventos futuros
├── examples_test.go            # Ejemplos verificables para godoc / pkg.go.dev
├── examples/
│   ├── basic_usage/            # Ejemplos básicos de conversión
│   ├── db_lifecycle_demo/      # Demostración del ciclo Ingesta -> BD -> Proyección
│   ├── enterprise_showcase/    # Demostración con facturación y días hábiles
│   ├── team_meeting_planner/   # Planificador de reuniones internacionales
│   └── two_database_patterns/  # Demostración de los 2 patrones de base de datos
├── LICENSE                     # Licencia Propietaria Exclusiva
└── README.md                   # Documentación técnica
```

---

## Autor y Aviso Legal

- **Creador y Autor**: Jhonatan
- **Licencia**: Licencia Propietaria. Todos los derechos reservados.

Queda estrictamente prohibida la copia, distribución, modificación o sublicenciamiento no autorizados de este software y sus archivos de documentación asociados. Consulta el archivo [LICENSE](LICENSE) para conocer los términos completos.
