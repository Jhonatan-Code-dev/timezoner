package project_test

import (
	"errors"
	"testing"
	"time"

	"timezoner/pkg/project"
	"timezoner/pkg/zone"
)

func TestProject_ForUser(t *testing.T) {
	utcTime := time.Date(2026, 9, 1, 15, 0, 0, 0, time.UTC)

	userLima, err := project.ForUser(utcTime, "America/Lima")
	if err != nil {
		t.Fatalf("ForUser falló: %v", err)
	}
	if userLima.LocalTime.Hour() != 10 {
		t.Errorf("Hora esperada 10, obtenida: %d", userLima.LocalTime.Hour())
	}
	if userLima.OffsetFormatted != "-05:00" {
		t.Errorf("Offset esperado -05:00")
	}

	userMadrid, err := project.ForUser(utcTime, "Europe/Madrid", "15:04 MST")
	if err != nil {
		t.Fatalf("ForUser Madrid falló: %v", err)
	}
	if userMadrid.LocalTime.Hour() != 17 {
		t.Errorf("Hora esperada en Madrid 17, obtenida: %d", userMadrid.LocalTime.Hour())
	}

	// Error con zona inválida
	if _, err := project.ForUser(utcTime, "Invalid/Zone"); err == nil || !errors.Is(err, zone.ErrInvalidZone) {
		t.Errorf("ForUser con zona inválida debería fallar")
	}
}

func TestProject_FormatAndBatch(t *testing.T) {
	utcTime := time.Date(2026, 9, 1, 15, 0, 0, 0, time.UTC)

	formatted, err := project.Format(utcTime, "America/Lima", "15:04")
	if err != nil || formatted != "10:00" {
		t.Errorf("Format esperado 10:00, obtenido: %s", formatted)
	}

	// Batch
	zones := []string{"America/Lima", "Asia/Tokyo"}
	batch, err := project.BatchForUsers(utcTime, zones)
	if err != nil || len(batch) != 2 {
		t.Fatalf("BatchForUsers falló")
	}

	// Batch vacío
	if _, err := project.BatchForUsers(utcTime, nil); err == nil || !errors.Is(err, project.ErrNoZonesProvided) {
		t.Errorf("BatchForUsers vacío debería retornar ErrNoZonesProvided")
	}
}
