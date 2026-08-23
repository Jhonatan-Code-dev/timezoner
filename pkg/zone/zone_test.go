package zone_test

import (
	"errors"
	"testing"
	"time"

	"github.com/Jhonatan-Code-dev/timezonermax/pkg/zone"
)

func TestLoadLocationAndAliases(t *testing.T) {
	tests := []struct {
		input       string
		expectError bool
		expectedLoc string
		expectedErr error
	}{
		{"UTC", false, "UTC", nil},
		{"PET", false, "America/Lima", nil},
		{"America/Lima", false, "America/Lima", nil},
		{"EST", false, "America/New_York", nil},
		{"JST", false, "Asia/Tokyo", nil},
		{"", true, "", zone.ErrEmptyZoneName},
		{"Invalid/NonExistent_Zone", true, "", zone.ErrInvalidZone},
	}

	for _, tc := range tests {
		loc, err := zone.LoadLocation(tc.input)
		if tc.expectError {
			if err == nil {
				t.Errorf("LoadLocation(%q): se esperaba error", tc.input)
			} else if tc.expectedErr != nil && !errors.Is(err, tc.expectedErr) {
				t.Errorf("LoadLocation(%q) err = %v, esperado errors.Is(%v)", tc.input, err, tc.expectedErr)
			}
		} else {
			if err != nil {
				t.Errorf("LoadLocation(%q): error inesperado: %v", tc.input, err)
			} else if loc.String() != tc.expectedLoc {
				t.Errorf("LoadLocation(%q) = %v; esperado %v", tc.input, loc.String(), tc.expectedLoc)
			}
		}
	}
}

func TestIsValidAndNormalize(t *testing.T) {
	if !zone.IsValid("America/Bogota") {
		t.Errorf("IsValid(America/Bogota) debería ser true")
	}
	if zone.IsValid("Mars/Zone") {
		t.Errorf("IsValid(Mars/Zone) debería ser false")
	}

	norm, err := zone.Normalize("COT")
	if err != nil || norm != "America/Bogota" {
		t.Errorf("Normalize(COT) = %s, err: %v; esperado America/Bogota", norm, err)
	}

	if _, err := zone.Normalize("Mars/NonExistent"); err == nil {
		t.Errorf("Normalize con zona inválida debería retornar error")
	}
}

func TestRegisterAlias(t *testing.T) {
	err := zone.RegisterAlias("PERU", "America/Lima")
	if err != nil {
		t.Fatalf("RegisterAlias falló: %v", err)
	}

	loc, err := zone.LoadLocation("PERU")
	if err != nil || loc.String() != "America/Lima" {
		t.Errorf("LoadLocation(PERU) = %v, esperado America/Lima", loc)
	}

	// Error al registrar alias para zona inválida
	if err := zone.RegisterAlias("INVALID", "Invalid/ZoneName"); err == nil {
		t.Errorf("RegisterAlias con zona inválida debería retornar error")
	}
}

func TestGetInfoAndDifference(t *testing.T) {
	base := time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)
	info, err := zone.GetInfo("America/Lima", base)
	if err != nil {
		t.Fatalf("GetInfo falló: %v", err)
	}
	if info.OffsetFormatted != "-05:00" {
		t.Errorf("Offset esperado -05:00, obtenido: %s", info.OffsetFormatted)
	}

	diff, err := zone.Difference("Europe/Madrid", "America/Lima", base)
	if err != nil {
		t.Fatalf("Difference falló: %v", err)
	}
	if diff != 7*time.Hour {
		t.Errorf("Diferencia esperada 7h, obtenida: %v", diff)
	}

	// Difference con zona A inválida
	if _, err := zone.Difference("Invalid/ZoneA", "America/Lima", base); err == nil {
		t.Errorf("Difference con zona A inválida debería retornar error")
	}

	// Difference con zona B inválida
	if _, err := zone.Difference("America/Lima", "Invalid/ZoneB", base); err == nil {
		t.Errorf("Difference con zona B inválida debería retornar error")
	}

	// GetInfo con zona inválida
	if _, err := zone.GetInfo("Invalid/Zone", base); err == nil {
		t.Errorf("GetInfo con zona inválida debería retornar error")
	}
}

func TestIsDST(t *testing.T) {
	// Madrid en Julio está en DST (CEST = UTC+2)
	summer := time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)
	isDstSummer, err := zone.IsDST("Europe/Madrid", summer)
	if err != nil || !isDstSummer {
		t.Errorf("Madrid en Julio debería estar en DST (isDST=true)")
	}

	// Madrid en Enero no está en DST (CET = UTC+1)
	winter := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	isDstWinter, err := zone.IsDST("Europe/Madrid", winter)
	if err != nil || isDstWinter {
		t.Errorf("Madrid en Enero NO debería estar en DST (isDST=false)")
	}

	// Hemisferio Sur: Sydney está en DST en Enero (AEDT = UTC+11)
	isDstSydneyJan, err := zone.IsDST("Australia/Sydney", winter)
	if err != nil || !isDstSydneyJan {
		t.Errorf("Sydney en Enero debería estar en DST (isDST=true)")
	}

	// Error con zona inválida
	if _, err := zone.IsDST("Invalid/Zone", summer); err == nil {
		t.Errorf("IsDST con zona inválida debería retornar error")
	}
}

func TestCommonZones_Immutability(t *testing.T) {
	zones := zone.CommonZones()
	if len(zones) == 0 {
		t.Fatal("CommonZones() no debe estar vacío")
	}

	// Mutar el slice retornado no debe afectar llamadas posteriores
	originalFirst := zones[0]
	zones[0] = "MODIFIED_ZONE"

	freshZones := zone.CommonZones()
	if freshZones[0] != originalFirst {
		t.Errorf("CommonZones() debe retornar copias defensivas inmutables")
	}
}
