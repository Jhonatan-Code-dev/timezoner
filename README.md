# Timezoner

Timezoner is a pure, high-performance Go package designed for precise IANA timezone conversions, temporal offset calculations, daylight saving time (DST) inspection, business calendar arithmetic, human-readable relative time, and distributed team meeting overlap scheduling.

It provides embedded tzdata support for zero-configuration portability, native `database/sql` driver integration, and dedicated sub-packages for database lifecycle handling:
- **`timezoner`**: Core conversions, Fluent API, business days arithmetic, bounds, relative time (`Humanize`), and `DBTime` type.
- **`timezoner/ingest`**: Normalizes and sanitizes all incoming timestamps into global UTC before database persistence.
- **`timezoner/project`**: Projects and formats stored UTC timestamps into the exact local time and format required by each user.

---

## Key Capabilities

- **Universal Portability (`time/tzdata`)**: Bundles the official IANA database within the binary. Works reliably on Windows, Alpine Linux, AWS Lambda, and Scratch Docker containers without OS zoneinfo dependencies.
- **Native SQL & JSON Type (`DBTime`)**: Implements `driver.Valuer`, `sql.Scanner`, `json.Marshaler`, and `json.Unmarshaler` to guarantee clean UTC storage without monotonic drift.
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

### 1. Business Days and Calendar Bounds

```go
package main

import (
	"fmt"
	"time"
	"timezoner"
)

func main() {
	// Starting from Friday
	friday := time.Date(2026, 9, 4, 10, 0, 0, 0, time.UTC)

	// Add 5 business days (skips Saturday/Sunday -> following Friday) and move to end of day
	dueDate := timezoner.At(friday).
		AddBusinessDays(5).
		EndOfDay().
		MustTime()

	fmt.Println("Due date:", dueDate.Format("2006-01-02 15:04:05 MST"))
	// Output: Due date: 2026-09-11 23:59:59 UTC
}
```

---

### 2. Native SQL Database & JSON Persistence (`DBTime`)

```go
package main

import (
	"encoding/json"
	"fmt"
	"time"
	"timezoner"
)

type Invoice struct {
	ID        string           `json:"id"`
	IssuedAt  timezoner.DBTime `json:"issued_at"`
	DueAt     timezoner.DBTime `json:"due_at"`
}

func main() {
	inv := Invoice{
		ID:       "INV-1001",
		IssuedAt: timezoner.NowDBTime(),
		DueAt:    timezoner.NewDBTime(time.Now().Add(7 * 24 * time.Hour)),
	}

	data, _ := json.MarshalIndent(inv, "", "  ")
	fmt.Println(string(data))
}
```

---

### 3. Human-Readable Relative Time

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

### 4. Database Ingestion & User Projection Pipeline

```go
package main

import (
	"fmt"
	"timezoner/ingest"
	"timezoner/project"
)

func main() {
	// 1. Ingestion: Convert local input to clean UTC for database storage
	dbUTC, err := ingest.FromString("2026-09-01 10:00:00", "America/Lima")
	if err != nil {
		panic(err)
	}
	fmt.Println("DB UTC:", dbUTC.Format("2006-01-02 15:04:05 UTC")) // 2026-09-01 15:00:00 UTC

	// 2. Projection: Adapt stored UTC record for a user in Madrid
	userMadrid, err := project.ForUser(dbUTC, "Europe/Madrid")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Madrid User: %s (Offset: %s, DST: %v)\n", 
		userMadrid.Formatted, userMadrid.OffsetFormatted, userMadrid.IsDST)
	// Output: Madrid User: 2026-09-01 17:00:00 (Offset: +02:00, DST: true)
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

| Package | Purpose |
| :--- | :--- |
| **`timezoner`** | Core timezone conversions, comparisons, overlap calculations, calendar arithmetic, relative formatting, and Fluent API. |
| **`timezoner/ingest`** | Sanitizes and converts incoming user dates/strings into monotonic-clean UTC for DB storage. |
| **`timezoner/project`** | Projects stored UTC records into user-specific timezones with ISO 8601 formatting and DST details. |

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
├── dbtime.go                   # SQL driver.Valuer, sql.Scanner, and JSON serialization
├── timezoner_test.go           # Unit tests, fuzz testing, concurrency, and benchmarks
├── calendar_test.go            # Unit tests for calendar and business days
├── humanize_test.go            # Unit tests for relative time formatting
├── dbtime_test.go              # Unit tests for SQL driver and JSON marshaling
├── examples_test.go            # Executable godoc examples
├── ingest/
│   ├── ingest.go               # Ingestion & normalization to UTC for DB
│   └── ingest_test.go          # Unit tests for ingest module
├── project/
│   ├── project.go              # Projection & formatting for destination users
│   └── project_test.go         # Unit tests for project module
├── examples/
│   ├── basic_usage/main.go     # Basic conversion examples
│   ├── db_lifecycle_demo/      # Ingestion -> DB -> Projection demo
│   ├── enterprise_showcase/    # Complete enterprise capabilities showcase
│   └── team_meeting_planner/   # Meeting planner example
├── LICENSE                     # Proprietary License
└── README.md                   # Technical documentation
```

---

## Author & Legal Notice

- **Author & Creator**: Jhonatan
- **License**: Proprietary License. All rights reserved.

Unauthorized copying, distribution, modification, or sublicensing of this software and associated documentation files is strictly prohibited. Refer to [LICENSE](LICENSE) for full terms.
