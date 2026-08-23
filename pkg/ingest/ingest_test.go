package ingest_test

import (
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
}
