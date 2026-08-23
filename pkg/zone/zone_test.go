package zone_test

import (
	"errors"
	"testing"
	"time"

	"timezoner/pkg/zone"
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
}
