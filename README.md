# Timezoner

Timezoner is a pure, zero-dependency Go package designed for precise IANA timezone conversions, temporal offset calculations, daylight saving time (DST) inspection, and distributed team meeting overlap scheduling.

---

## Overview

Working with distributed systems across timezones often requires complex conversions and edge-case handling around daylight saving transitions and non-standard offsets. Timezoner provides an idiomatic API built on top of the Go standard library with built-in concurrency safety and memory-cached location resolution.

### Key Capabilities

- **IANA Timezone Conversion**: Safe conversion between any valid IANA location or common abbreviation alias (`UTC`, `EST`, `CET`, `PET`, `JST`).
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

## Usage

### 1. Basic Conversion

```go
package main

import (
	"fmt"
	"time"
	"timezoner"
)

func main() {
	now := time.Now()

	// Convert to Tokyo time
	tokyoTime, err := timezoner.Convert(now, "Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	fmt.Println("Tokyo time:", tokyoTime.Format("2006-01-02 15:04:05 MST"))

	// Direct lookup using common alias (PET = America/Lima)
	limaTime, err := timezoner.NowIn("PET")
	if err != nil {
		panic(err)
	}
	fmt.Println("Lima time:", limaTime.Format("15:04:05"))
}
```

---

### 2. Fluent API

```go
package main

import (
	"fmt"
	"timezoner"
)

func main() {
	formatted, err := timezoner.Now().
		In("Europe/Paris").
		Format("2006-01-02 15:04:05 MST")

	if err != nil {
		panic(err)
	}

	fmt.Println("Paris:", formatted)
}
```

---

### 3. Time Difference and DST Detection

```go
package main

import (
	"fmt"
	"time"
	"timezoner"
)

func main() {
	targetDate := time.Date(2026, 7, 15, 12, 0, 0, 0, time.UTC)

	// Calculate exact offset difference between two zones
	diff, err := timezoner.Difference("Europe/Madrid", "America/Lima", targetDate)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Difference: %v hours\n", diff.Hours()) // +7.0 hours

	// Check if daylight saving time is active
	isDST, err := timezoner.IsDST("Europe/Madrid", targetDate)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Madrid in DST: %v\n", isDST) // true in July
}
```

---

### 4. Distributed Team Meeting Overlap

Calculate overlapping business hours for teams across different continents:

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
		fmt.Printf("  UTC:       %s - %s\n", 
			slot.StartTimeUTC.Format("15:04"), 
			slot.EndTimeUTC.Format("15:04"))
		fmt.Printf("  Lima:      %s - %s\n", 
			slot.ZoneTimes["America/Lima"].Format("15:04"), 
			slot.ZoneTimes["America/Lima"].Add(slot.Duration).Format("15:04"))
		fmt.Printf("  New York:  %s - %s\n", 
			slot.ZoneTimes["America/New_York"].Format("15:04"), 
			slot.ZoneTimes["America/New_York"].Add(slot.Duration).Format("15:04"))
		fmt.Printf("  Madrid:    %s - %s\n", 
			slot.ZoneTimes["Europe/Madrid"].Format("15:04"), 
			slot.ZoneTimes["Europe/Madrid"].Add(slot.Duration).Format("15:04"))
	}
}
```

---

### 5. Multi-Zone Comparison

```go
package main

import (
	"fmt"
	"time"
	"timezoner"
)

func main() {
	snapshots, err := timezoner.Compare(time.Now(), "America/Lima", "Europe/London", "Asia/Tokyo")
	if err != nil {
		panic(err)
	}

	for _, s := range snapshots {
		fmt.Printf("%-20s | %s | Offset: %s | DST: %v\n", 
			s.Zone, s.Formatted, s.OffsetFormatted, s.IsDST)
	}
}
```

---

## API Summary

| Function / Type | Description |
| :--- | :--- |
| `Convert(t, toZone)` | Converts a `time.Time` to a destination zone. |
| `ConvertBetween(str, layout, from, to)` | Parses a formatted string in source zone and converts to target zone. |
| `NowIn(zone)` | Returns current time in the specified zone or alias. |
| `FormatIn(t, zone, layout)` | Formats a time in target zone using standard layout string. |
| `GetZoneInfo(zone, at)` | Returns detailed timezone struct (offset, abbreviation, DST, UTC diff). |
| `IsDST(zone, at)` | Reports whether daylight saving time is active for a given zone and instant. |
| `Difference(zoneA, zoneB, at)` | Computes time difference duration between two locations. |
| `Compare(at, zones...)` | Returns simultaneous time snapshots across multiple zones. |
| `FindOverlap(req)` | Calculates intersecting working hour slots for distributed teams. |
| `At(t)` / `Now()` | Initializes a `TimePoint` for fluent method chaining. |

---

## Error Handling

Timezoner exports standard sentinel errors for explicit verification via `errors.Is`:

- `ErrEmptyZoneName`: Returned when an empty string is passed as a timezone identifier.
- `ErrInvalidZone`: Returned when a timezone name cannot be resolved in the IANA database.
- `ErrInvalidTimeFormat`: Returned when a date string fails to parse with the supplied layout.
- `ErrNoZonesProvided`: Returned when a multi-zone operation is invoked with an empty list.

---

## Testing & Benchmarks

Run unit tests and coverage report:

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
├── examples/
│   ├── basic_usage/main.go     # Basic conversion examples
│   └── team_meeting_planner/   # Meeting planner example
├── LICENSE                     # Proprietary License
└── README.md                   # Technical documentation
```

---

## Author & Legal Notice

- **Author & Creator**: Jhonatan
- **License**: Proprietary License. All rights reserved.

Unauthorized copying, distribution, modification, or sublicensing of this software and associated documentation files is strictly prohibited. Refer to [LICENSE](LICENSE) for full terms.
