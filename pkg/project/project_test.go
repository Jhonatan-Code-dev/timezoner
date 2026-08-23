package project_test

import (
	"testing"
	"time"

	"timezoner/pkg/project"
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

	userMadrid, err := project.ForUser(utcTime, "Europe/Madrid")
	if err != nil {
		t.Fatalf("ForUser Madrid falló: %v", err)
	}
	if userMadrid.LocalTime.Hour() != 17 {
		t.Errorf("Hora esperada en Madrid 17, obtenida: %d", userMadrid.LocalTime.Hour())
	}
}
