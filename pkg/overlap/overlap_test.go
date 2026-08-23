package overlap_test

import (
	"errors"
	"testing"
	"time"

	"github.com/Jhonatan-Code-dev/timezoner/pkg/overlap"
	"github.com/Jhonatan-Code-dev/timezoner/pkg/zone"
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

	for _, s := range slots {
		if s.Duration < 1*time.Hour {
			t.Errorf("Duración del slot menor que SlotDuration: %v", s.Duration)
		}
	}
}

func TestOverlap_CustomHours(t *testing.T) {
	custom := map[string]overlap.WorkingHours{
		"America/Lima":  {StartHour: 8, EndHour: 17},
		"Europe/Madrid": {StartHour: 10, EndHour: 19},
	}

	req := overlap.Request{
		Date:         time.Date(2026, 10, 15, 0, 0, 0, 0, time.UTC),
		Zones:        []string{"America/Lima", "Europe/Madrid"},
		CustomHours:  custom,
		SlotDuration: 30 * time.Minute,
	}

	slots, err := overlap.Find(req)
	if err != nil {
		t.Fatalf("Find con custom hours falló: %v", err)
	}
	if len(slots) == 0 {
		t.Errorf("Se esperaba solapamiento con custom hours")
	}
}

func TestOverlap_NoMatch(t *testing.T) {
	// Horarios estrictos sin solapamiento
	custom := map[string]overlap.WorkingHours{
		"Pacific/Honolulu": {StartHour: 9, EndHour: 10},
		"Asia/Tokyo":       {StartHour: 9, EndHour: 10},
	}

	req := overlap.Request{
		Date:         time.Date(2026, 10, 15, 0, 0, 0, 0, time.UTC),
		Zones:        []string{"Pacific/Honolulu", "Asia/Tokyo"},
		CustomHours:  custom,
		SlotDuration: 1 * time.Hour,
	}

	slots, err := overlap.Find(req)
	if err != nil {
		t.Fatalf("Find falló: %v", err)
	}
	if len(slots) != 0 {
		t.Errorf("No debería haber solapamiento entre estos horarios, encontrados %d slots", len(slots))
	}
}

func TestOverlap_ErrorsAndDefaults(t *testing.T) {
	// Sin zonas
	if _, err := overlap.Find(overlap.Request{}); err == nil || !errors.Is(err, overlap.ErrNoZonesProvided) {
		t.Errorf("Find sin zonas debería retornar ErrNoZonesProvided")
	}

	// Zona inválida
	reqInvalid := overlap.Request{
		Date:  time.Now(),
		Zones: []string{"Invalid/ZoneName"},
	}
	if _, err := overlap.Find(reqInvalid); err == nil || !errors.Is(err, zone.ErrInvalidZone) {
		t.Errorf("Find con zona inválida debería retornar ErrInvalidZone")
	}

	// Defaults (SlotDuration=0 y DefaultHours={0,0})
	reqDefaults := overlap.Request{
		Date:  time.Date(2026, 10, 15, 0, 0, 0, 0, time.UTC),
		Zones: []string{"America/Lima", "America/New_York"},
	}
	slots, err := overlap.Find(reqDefaults)
	if err != nil || len(slots) == 0 {
		t.Errorf("Find con defaults falló")
	}
}
