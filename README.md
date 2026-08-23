# Timezoner

Timezoner is a pure, zero-dependency Go package designed for precise IANA timezone conversions, temporal offset calculations, daylight saving time (DST) inspection, and distributed team meeting overlap scheduling.

It includes dedicated sub-packages for handling the complete database datetime lifecycle:
- **`timezoner/ingest`**: Normalizes and sanitizes all incoming timestamps into global UTC before database persistence.
- **`timezoner/project`**: Projects and formats stored UTC timestamps into the exact local time and format required by each user.

---

## Overview

Working with distributed systems across timezones requires robust handling of daylight saving transitions, non-standard offsets, and persistence standardization. Timezoner provides an idiomatic API built on top of the Go standard library with built-in concurrency safety and memory-cached location resolution.

### Key Capabilities

- **Database Ingestion (`ingest`)**: Convert user inputs, timestamps, and date strings from any local zone to standardized UTC.
- **User Projection (`project`)**: Adapt stored UTC records into localized dates with ISO 8601 strings, offsets, and DST flags.
- **Meeting Overlap Calculation**: Calculates matching working hours across global teams for any target calendar date.
- **DST & Offset Inspection**: Determine active daylight saving status and exact duration offsets between two zones.
- **Fluent API**: Chainable methods for building and transforming temporal representations.
- **Zero External Dependencies**: Implemented entirely with the Go standard library (`time`, `sync`, `errors`, `fmt`).
- **Thread-Safe Caching**: Concurrent in-memory cache (`sync.Map`) for `*time.Location` lookups to minimize runtime overhead.

---

## Installation

```bash
go get timezoner
```

Requires Go 1.22 or higher.

---

## Complete Database Lifecycle Workflow

```
[ Incoming User Input ]  -->  [ timezoner/ingest ]  -->  [ Database (UTC) ]  -->  [ timezoner/project ]  -->  [ Local User Response ]
```

### 1. Ingestion: Normalizing to UTC Before Database Storage

```go
package main

import (
	"fmt"
	"timezoner/ingest"
)

func main() {
	// A user in Lima submits a local datetime string
	inputDate := "2026-09-01 10:00:00"
	sourceZone := "America/Lima"

	// Normalize to absolute UTC for database storage
	dbTimeUTC, err := ingest.FromString(inputDate, sourceZone)
	if err != nil {
		panic(err)
	}

	fmt.Println("Save to DB (UTC):", dbTimeUTC.Format("2006-01-02 15:04:05 UTC"))
	// Output: 2026-09-01 15:00:00 UTC
}
```

---

### 2. Projection: Reading from Database for Different Users

```go
package main

import (
	"fmt"
	"time"
	"timezoner/project"
)

func main() {
	// Value retrieved from DB in UTC
	dbRecordUTC := time.Date(2026, 9, 1, 15, 0, 0, 0, time.UTC)

	// Project for a viewer in Madrid
	userMadrid, err := project.ForUser(dbRecordUTC, "Europe/Madrid")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Madrid: %s (Offset: %s, DST: %v)\n", 
		userMadrid.Formatted, userMadrid.OffsetFormatted, userMadrid.IsDST)
	// Output: Madrid: 2026-09-01 17:00:00 (Offset: +02:00, DST: true)
}
```

---

### 3. Core Conversion & Fluent API

```go
package main

import (
	"fmt"
	"timezoner"
)

func main() {
	formatted, err := timezoner.Now().
		In("Asia/Tokyo").
		Format("2006-01-02 15:04:05 MST")

	if err != nil {
		panic(err)
	}

	fmt.Println("Tokyo:", formatted)
}
```

---

### 4. Distributed Team Meeting Overlap

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
		fmt.Printf("  Lima:      %s - %s\n", slot.ZoneTimes["America/Lima"].Format("15:04"), slot.ZoneTimes["America/Lima"].Add(slot.Duration).Format("15:04"))
		fmt.Printf("  New York:  %s - %s\n", slot.ZoneTimes["America/New_York"].Format("15:04"), slot.ZoneTimes["America/New_York"].Add(slot.Duration).Format("15:04"))
		fmt.Printf("  Madrid:    %s - %s\n", slot.ZoneTimes["Europe/Madrid"].Format("15:04"), slot.ZoneTimes["Europe/Madrid"].Add(slot.Duration).Format("15:04"))
	}
}
```

---

## Package Architecture

| Package | Purpose |
| :--- | :--- |
| **`timezoner`** | Core timezone conversions, comparisons, overlap calculations, and Fluent API. |
| **`timezoner/ingest`** | Sanitizes and converts incoming user dates/strings into monotonic-clean UTC for DB storage. |
| **`timezoner/project`** | Projects stored UTC records into user-specific timezones with ISO 8601 formatting and DST details. |

---

## Error Handling

Timezoner exports standard sentinel errors for explicit verification via `errors.Is`:

- `timezoner.ErrEmptyZoneName`: Returned when an empty string is passed as a timezone identifier.
- `timezoner.ErrInvalidZone`: Returned when a timezone name cannot be resolved in the IANA database.
- `timezoner.ErrInvalidTimeFormat`: Returned when a date string fails to parse with the supplied layout.
- `timezoner.ErrNoZonesProvided`: Returned when a multi-zone operation is invoked with an empty list.
- `ingest.ErrEmptyDateString`: Returned when an empty date string is submitted for ingestion.

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
├── zones.go                    # IANA catalog, location cache, and alias dictionary
├── timezoner_test.go           # Unit tests, fuzz testing, concurrency, and benchmarks
├── examples_test.go            # Executable godoc examples
├── ingest/
│   ├── ingest.go               # Ingestion & normalization to UTC for DB
│   └── ingest_test.go          # Unit tests for ingest module
├── project/
│   ├── project.go              # Projection & formatting for destination users
│   └── project_test.go         # Unit tests for project module
├── examples/
│   ├── basic_usage/main.go     # Basic conversion examples
│   ├── db_lifecycle_demo/      # Full Ingestion -> DB -> Projection demo
│   └── team_meeting_planner/   # Meeting planner example
├── LICENSE                     # Proprietary License
└── README.md                   # Technical documentation
```

---

## Author & Legal Notice

- **Author & Creator**: Jhonatan
- **License**: Proprietary License. All rights reserved.

Unauthorized copying, distribution, modification, or sublicensing of this software and associated documentation files is strictly prohibited. Refer to [LICENSE](LICENSE) for full terms.
