package ingest_test

import (
	"errors"
	"testing"
	"time"

	"timezoner/pkg/ingest"
)

func TestIngest_FromStringAndLocal(t *testing.T) {
	utcTime, err := ingest.FromString("2026-09-01 10:00", "America/Lima")
	if err != nil {
		t.Fatalf("FromString falló: %v", err)
	}
	if utcTime.Hour() != 15 {
		t.Errorf("Hora esperada 15, obtenida: %d", utcTime.Hour())
	}

	locTime := time.Date(2026, 9, 1, 10, 0, 0, 0, time.UTC)
	fromLoc, err := ingest.FromLocal(locTime, "America/Lima")
	if err != nil {
		t.Fatalf("FromLocal falló: %v", err)
	}
	if fromLoc.Hour() != 15 {
		t.Errorf("FromLocal hora esperada 15, obtenida: %d", fromLoc.Hour())
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
}
