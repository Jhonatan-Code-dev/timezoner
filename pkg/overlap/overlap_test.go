package overlap_test

import (
	"testing"
	"time"

	"timezoner/pkg/overlap"
)

func TestOverlap_Find(t *testing.T) {
	req := overlap.Request{
		Date:         time.Date(2026, 10, 15, 0, 0, 0, 0, time.UTC),
		Zones:        []string{"America/Lima", "America/New_York", "Europe/Madrid"},
		DefaultHours: overlap.WorkingHours{StartHour: 9, EndHour: 18},
		SlotDuration: 1 * time.Hour,
	}

	slots, err := overlap.Find(req)
	if err != nil {
		t.Fatalf("Find falló: %v", err)
	}
	if len(slots) == 0 {
		t.Fatalf("Se esperaba encontrar solapamiento")
	}
}
