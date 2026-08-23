package project_test

import (
	"errors"
	"testing"
	"time"

	"timezoner"
	"timezoner/project"
)

func TestForUser(t *testing.T) {
	// Base UTC: 2026-09-01 15:00:00 UTC
	baseUTC := time.Date(2026, 9, 1, 15, 0, 0, 0, time.UTC)

	// Proyectar para usuario en Lima (UTC-5) -> 10:00
	userLima, err := project.ForUser(baseUTC, "America/Lima")
	if err != nil {
		t.Fatalf("ForUser Lima falló: %v", err)
	}

	if userLima.LocalTime.Hour() != 10 {
		t.Errorf("Hora esperada en Lima: 10, obtenida: %d", userLima.LocalTime.Hour())
	}
	if userLima.OffsetFormatted != "-05:00" {
		t.Errorf("Offset esperado -05:00, obtenido: %s", userLima.OffsetFormatted)
	}

	// Proyectar para usuario en Madrid (UTC+2 CEST en septiembre) -> 17:00
	userMadrid, err := project.ForUser(baseUTC, "Europe/Madrid")
	if err != nil {
		t.Fatalf("ForUser Madrid falló: %v", err)
	}
	if userMadrid.LocalTime.Hour() != 17 {
		t.Errorf("Hora esperada en Madrid: 17, obtenida: %d", userMadrid.LocalTime.Hour())
	}
	if !userMadrid.IsDST {
		t.Errorf("Madrid en septiembre debería tener IsDST = true")
	}

	// Probar error con zona inválida
	_, err = project.ForUser(baseUTC, "Zone/Fake")
	if err == nil {
		t.Errorf("Se esperaba error con zona inválida")
	}
}

func TestFormat(t *testing.T) {
	baseUTC := time.Date(2026, 9, 1, 15, 0, 0, 0, time.UTC)
	formatted, err := project.Format(baseUTC, "America/Lima", "15:04 MST")
	if err != nil {
		t.Fatalf("Format falló: %v", err)
	}

	if formatted != "10:00 -05" && formatted != "10:00 PET" {
		t.Logf("Format retornó: %s", formatted)
	}
}

func TestBatchForUsers(t *testing.T) {
	baseUTC := time.Date(2026, 9, 1, 15, 0, 0, 0, time.UTC)
	zones := []string{"America/Lima", "Europe/Madrid", "Asia/Tokyo"}

	batch, err := project.BatchForUsers(baseUTC, zones)
	if err != nil {
		t.Fatalf("BatchForUsers falló: %v", err)
	}

	if len(batch) != 3 {
		t.Fatalf("Se esperaban 3 proyecciones, obtenidas: %d", len(batch))
	}
	if batch["America/Lima"].LocalTime.Hour() != 10 {
		t.Errorf("Lima hora esperada 10, obtenida: %d", batch["America/Lima"].LocalTime.Hour())
	}
	if batch["Europe/Madrid"].LocalTime.Hour() != 17 {
		t.Errorf("Madrid hora esperada 17, obtenida: %d", batch["Europe/Madrid"].LocalTime.Hour())
	}
	if batch["Asia/Tokyo"].LocalTime.Hour() != 0 { // 15:00 UTC + 9h = 00:00 del día siguiente
		t.Errorf("Tokyo hora esperada 0, obtenida: %d", batch["Asia/Tokyo"].LocalTime.Hour())
	}

	// Error sin zonas
	_, err = project.BatchForUsers(baseUTC, nil)
	if err == nil || !errors.Is(err, timezoner.ErrNoZonesProvided) {
		t.Errorf("Se esperaba ErrNoZonesProvided")
	}
}
