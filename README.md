# Timezoner

Timezoner is a pure, high-performance Go package designed for precise IANA timezone conversions, temporal offset calculations, daylight saving time (DST) inspection, business calendar arithmetic, human-readable relative time, and distributed team meeting overlap scheduling.

It provides embedded tzdata support for zero-configuration portability, native `database/sql` driver integration, and two enterprise database persistence patterns (`DBTime` and `ZonedTime`).

---

## The Two Enterprise Database Patterns

Timezoner supports both industry-standard persistence patterns out of the box:

| Pattern | Type | Best For | Storage Strategy |
| :--- | :--- | :--- | :--- |
| **Pattern 1: Transactional / Audit** | `timezoner.DBTime` | Payments, orders, logs, sensor data, chat messages. | Single column stored as pure UTC (`TIMESTAMPTZ`). |
| **Pattern 2: Future Schedules & Calendars** | `timezoner.ZonedTime` | Medical appointments, flights, webinars, alarms. | Dual columns (or JSONB) storing UTC instant + origin IANA zone name (`"America/Lima"`). Preserves original wall-clock time across DST law revisions. |

---

## Key Capabilities

- **Universal Portability (`time/tzdata`)**: Bundles the official IANA database within the binary. Works reliably on Windows, Alpine Linux, AWS Lambda, and Scratch Docker containers without OS zoneinfo dependencies.
- **Native SQL & JSON Types (`DBTime`, `ZonedTime`)**: Implements `driver.Valuer`, `sql.Scanner`, `json.Marshaler`, and `json.Unmarshaler` to guarantee clean UTC storage without monotonic drift.
- **Business Calendar Arithmetic**: Add or subtract business days (`AddBusinessDays`) skipping weekends, and compute day/month boundaries (`StartOfDay`, `EndOfMonth`).
- **Human-Readable Relative Time (`Humanize`)**: Converts durations to natural language strings in Spanish and English (`"hace 5 minutos"`, `"in 2 hours"`).
- **Database Ingestion & User Projection**: Complete pipelines for converting local user inputs to UTC and projecting UTC records back to any timezone.
- **Meeting Overlap Calculation**: Calculates matching working hours across global teams for any target calendar date.
- **Thread-Safe Caching**: Concurrent in-memory cache (`sync.Map`) for `*time.Location` lookups to minimize runtime overhead.
- **Zero External Dependencies**: Implemented entirely with the Go standard library (`time`, `sync`, `errors`, `database/sql/driver`, `fmt`).

---

## Installation

```bash
go get timezoner
```

Requires Go 1.22 or higher.

---

## Usage Examples

### 1. Pattern 1: Payments & Transactions (Pure UTC via `DBTime`)

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

	// Project for a user in Madrid
	madridView, _ := timezoner.ProjectForUser(p.PaidAtUTC.Time, "Europe/Madrid")
	fmt.Printf("Madrid User sees: %s (%s)\n", madridView.Formatted, madridView.OffsetFormatted)
}
```

---

### 2. Pattern 2: Future Appointments & Calendars (`ZonedTime`)

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
	// A patient in Lima books for 10:00 AM
	zoned, err := timezoner.ZonedFromLocal("2026-10-01 10:00", "America/Lima")
	if err != nil {
		panic(err)
	}

	app := Appointment{
		ID:          1,
		Doctor:      "Dr. Smith",
		ScheduledAt: zoned,
	}

	// 1. Guaranteed original local time in Lima
	localTime, _ := app.ScheduledAt.Local()
	fmt.Println("Lima:", localTime.Format("2006-01-02 15:04")) // 2026-10-01 10:00

	// 2. Projected for a specialist in Tokyo
	tokyoView, _ := app.ScheduledAt.ForViewer("Asia/Tokyo")
	fmt.Printf("Tokyo: %s (%s)\n", tokyoView.Formatted, tokyoView.OffsetFormatted)
}
```

---

### 3. Business Days and Calendar Bounds

```go
package main

import (
	"fmt"
	"time"
	"timezoner"
)

func main() {
	friday := time.Date(2026, 9, 4, 10, 0, 0, 0, time.UTC)

	dueDate := timezoner.At(friday).
		AddBusinessDays(5). // Skips Saturday/Sunday -> Following Friday
		EndOfDay().         // Sets 23:59:59
		MustTime()

	fmt.Println("Due date:", dueDate.Format("2006-01-02 15:04:05 MST"))
}
```

---

### 4. Human-Readable Relative Time

```go
package main

import (
	"fmt"
	"time"
	"timezoner"
)

func main() {
	t := time.Now().Add(-2 * time.Hour)

	fmt.Println("Spanish:", timezoner.Humanize(t))   // "hace 2 horas"
	fmt.Println("English:", timezoner.HumanizeEn(t)) // "2 hours ago"
}
```

---

### 5. Distributed Team Meeting Overlap

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
		fmt.Printf("Window #%d (Duration: %v):\n", i+1, slot.Duration)
		fmt.Printf("  UTC:       %s - %s\n", slot.StartTimeUTC.Format("15:04"), slot.EndTimeUTC.Format("15:04"))
		fmt.Printf("  Lima:      %s\n", slot.ZoneTimes["America/Lima"].Format("15:04"))
		fmt.Printf("  New York:  %s\n", slot.ZoneTimes["America/New_York"].Format("15:04"))
		fmt.Printf("  Madrid:    %s\n", slot.ZoneTimes["Europe/Madrid"].Format("15:04"))
	}
}
```

---

## Package Architecture

| Package / Type | Purpose |
| :--- | :--- |
| **`timezoner.DBTime`** | Single-column SQL/JSON UTC type for transactions, payments, logs. |
| **`timezoner.ZonedTime`** | Dual-column SQL/JSON type (UTC + IANA zone) for appointments, calendars, flights. |
| **`timezoner.Ingest*`** | Sanitizes and converts incoming user dates/strings into monotonic-clean UTC for DB storage. |
| **`timezoner.Project*`** | Projects stored UTC records into user-specific timezones with ISO 8601 formatting and DST details. |

---

## Testing & Benchmarks

Run unit tests across all packages:

```bash
go test -v -cover ./...
```

Run benchmarks:

```bash
go test -bench=. -benchmem ./...
```

---

## Repository Layout

```
.
├── go.mod                      # Module definition
├── timezoner.go                # Core library and Fluent API
├── zones.go                    # IANA catalog, location cache, embedded tzdata, and alias dictionary
├── calendar.go                 # Business days arithmetic, week day check, and time bounds
├── humanize.go                 # Natural language relative time (ES/EN)
├── dbtime.go                   # Single-column SQL driver.Valuer / sql.Scanner (Pure UTC)
├── zonedtime.go                # Dual-column SQL/JSON type (UTC + IANA zone for calendars)
├── ingest.go                   # Ingestion & normalization to UTC for DB
├── project.go                  # Projection & formatting for destination users
├── timezoner_test.go           # Unit tests, fuzz testing, concurrency, and benchmarks
├── calendar_test.go            # Unit tests for calendar and business days
├── humanize_test.go            # Unit tests for relative time formatting
├── dbtime_test.go              # Unit tests for SQL driver and JSON marshaling
├── zonedtime_test.go           # Unit tests for ZonedTime future schedules
├── examples_test.go            # Executable godoc examples
├── examples/
│   ├── basic_usage/            # Basic conversion examples
│   ├── db_lifecycle_demo/      # Ingestion -> DB -> Projection demo
│   ├── enterprise_showcase/    # Complete enterprise capabilities showcase
│   ├── team_meeting_planner/   # Meeting planner example
│   └── two_database_patterns/  # The 2 enterprise DB persistence patterns demo
├── LICENSE                     # Proprietary License
└── README.md                   # Technical documentation
```

---

## Author & Legal Notice

- **Author & Creator**: Jhonatan
- **License**: Proprietary License. All rights reserved.

Unauthorized copying, distribution, modification, or sublicensing of this software and associated documentation files is strictly prohibited. Refer to [LICENSE](LICENSE) for full terms.
