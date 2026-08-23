package ingest_test

import (
	"errors"
	"testing"
	"time"

	"github.com/Jhonatan-Code-dev/timezonermax/pkg/ingest"
	"github.com/Jhonatan-Code-dev/timezonermax/pkg/zone"
)

func TestIngest_FromStringAndLocal(t *testing.T) {
	utcTime, err := ingest.FromString("2026-09-01 10:00", "America/Lima")
	if err != nil {
		t.Fatalf("FromString falló: %v", err)
	}
	if utcTime.Hour() != 15 {
		t.Errorf("Hora esperada 15, obtenida: %d", utcTime.Hour())
	}

	// FromString con defaultZone vacío (debe asumir UTC)
	utcDefault, err := ingest.FromString("2026-09-01 10:00:00", "")
	if err != nil || utcDefault.Hour() != 10 {
		t.Errorf("FromString con defaultZone vacío falló, obtenido: %v", utcDefault)
	}

	locTime := time.Date(2026, 9, 1, 10, 0, 0, 0, time.UTC)
	fromLoc, err := ingest.FromLocal(locTime, "America/Lima")
	if err != nil {
		t.Fatalf("FromLocal falló: %v", err)
	}
	if fromLoc.Hour() != 15 {
		t.Errorf("FromLocal hora esperada 15, obtenida: %d", fromLoc.Hour())
	}

	// FromLocal con zona inválida
	if _, err := ingest.FromLocal(locTime, "Invalid/Zone"); err == nil || !errors.Is(err, zone.ErrInvalidZone) {
		t.Errorf("FromLocal con zona inválida debería retornar ErrInvalidZone")
	}

	// Now y FromTime
	now := ingest.Now()
	if now.Location() != time.UTC {
		t.Errorf("Now debe ser UTC")
	}

	fTime := ingest.FromTime(time.Now())
	if fTime.Location() != time.UTC {
		t.Errorf("FromTime debe ser UTC")
	}

	// Unix
	uSec := ingest.FromUnix(1700000000)
	uMilli := ingest.FromUnixMilli(1700000000000)
	if !uSec.Equal(uMilli) {
		t.Errorf("Unix y UnixMilli deben coincidir")
	}

	// Errores
	if _, err := ingest.FromString("", "UTC"); err == nil || !errors.Is(err, ingest.ErrEmptyDateString) {
		t.Errorf("FromString con string vacío debería retornar ErrEmptyDateString")
	}

	if _, err := ingest.FromString("invalid-format", "UTC"); err == nil || !errors.Is(err, ingest.ErrInvalidInput) {
		t.Errorf("FromString con formato inválido debería retornar ErrInvalidInput")
	}

	if _, err := ingest.FromString("2026-09-01 10:00", "Invalid/ZoneName"); err == nil || !errors.Is(err, zone.ErrInvalidZone) {
		t.Errorf("FromString con zona inválida debería retornar ErrInvalidZone")
	}
}

func TestIngest_AllSupportedLayouts(t *testing.T) {
	testDates := []string{
		"2026-09-01T15:04:05.123456789Z",
		"2026-09-01T15:04:05Z",
		"2026-09-01T15:04:05",
		"2026-09-01 15:04:05",
		"2026-09-01 15:04",
		"2026-09-01",
		"01/09/2026 15:04:05",
		"01/09/2026 15:04",
		"01/09/2026",
	}

	for _, d := range testDates {
		parsed, err := ingest.FromString(d, "UTC")
		if err != nil {
			t.Errorf("Layout %q debería parsear con éxito, error: %v", d, err)
		}
		if parsed.Year() != 2026 || parsed.Month() != time.September || parsed.Day() != 1 {
			t.Errorf("Fecha parseada incorrecta para %q: %v", d, parsed)
		}
	}
}

func TestIngest_SupportedLayouts_Immutability(t *testing.T) {
	layouts := ingest.SupportedLayouts()
	if len(layouts) == 0 {
		t.Fatal("SupportedLayouts() no debe estar vacío")
	}

	// Mutar el slice retornado no debe afectar llamadas posteriores
	originalFirst := layouts[0]
	layouts[0] = "MODIFIED_LAYOUT"

	freshLayouts := ingest.SupportedLayouts()
	if freshLayouts[0] != originalFirst {
		t.Errorf("SupportedLayouts() debe retornar copias defensivas inmutables")
	}
}
