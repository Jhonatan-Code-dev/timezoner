package ingest_test

import (
	"errors"
	"testing"
	"time"

	"timezoner/ingest"
)

func TestNowAndFromTime(t *testing.T) {
	nowUTC := ingest.Now()
	if nowUTC.Location() != time.UTC {
		t.Errorf("ingest.Now() debe retornar time.UTC, obtenido: %v", nowUTC.Location())
	}

	locLima, _ := time.LoadLocation("America/Lima")
	localTime := time.Now().In(locLima)
	normalized := ingest.FromTime(localTime)

	if normalized.Location() != time.UTC {
		t.Errorf("ingest.FromTime() debe convertir a UTC, obtenido: %v", normalized.Location())
	}
}

func TestFromLocal(t *testing.T) {
	// 10:00 en Lima (UTC-5) son 15:00 UTC
	localTime := time.Date(2026, 9, 1, 10, 0, 0, 0, time.UTC)
	utcTime, err := ingest.FromLocal(localTime, "America/Lima")
	if err != nil {
		t.Fatalf("FromLocal falló: %v", err)
	}

	if utcTime.Hour() != 15 {
		t.Errorf("Hora esperada en UTC: 15, obtenida: %d", utcTime.Hour())
	}

	// Error con zona inválida
	_, err = ingest.FromLocal(localTime, "Invalid/Zone")
	if err == nil {
		t.Errorf("Se esperaba error con zona inválida")
	}
}

func TestFromString(t *testing.T) {
	tests := []struct {
		input       string
		defaultZone string
		expectedH   int
		expectErr   bool
	}{
		{"2026-09-01T15:00:00Z", "UTC", 15, false},
		{"2026-09-01 10:00:00", "America/Lima", 15, false}, // 10:00 Lima -> 15:00 UTC
		{"01/09/2026 10:00", "PET", 15, false},
		{"", "UTC", 0, true},
		{"invalid-date-string", "UTC", 0, true},
	}

	for _, tc := range tests {
		res, err := ingest.FromString(tc.input, tc.defaultZone)
		if tc.expectErr {
			if err == nil {
				t.Errorf("FromString(%q) se esperaba error y no ocurrió", tc.input)
			}
		} else {
			if err != nil {
				t.Errorf("FromString(%q) error inesperado: %v", tc.input, err)
			} else if res.Hour() != tc.expectedH {
				t.Errorf("FromString(%q) hora UTC esperada: %d, obtenida: %d", tc.input, tc.expectedH, res.Hour())
			}
		}
	}
}

func TestFromUnix(t *testing.T) {
	sec := int64(1700000000)
	tUnix := ingest.FromUnix(sec)
	if tUnix.Location() != time.UTC {
		t.Errorf("FromUnix debe retornar UTC")
	}

	tMilli := ingest.FromUnixMilli(sec * 1000)
	if tMilli.Location() != time.UTC {
		t.Errorf("FromUnixMilli debe retornar UTC")
	}
	if !tUnix.Equal(tMilli) {
		t.Errorf("FromUnix y FromUnixMilli deberían ser equivalentes")
	}
}

func TestEmptyStringError(t *testing.T) {
	_, err := ingest.FromString("", "UTC")
	if !errors.Is(err, ingest.ErrEmptyDateString) {
		t.Errorf("Se esperaba ErrEmptyDateString, obtenido: %v", err)
	}
}
