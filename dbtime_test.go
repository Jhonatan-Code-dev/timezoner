package timezoner_test

import (
	"encoding/json"
	"testing"
	"time"

	"timezoner"
)

type EventRecord struct {
	ID        int              `json:"id"`
	Title     string           `json:"title"`
	CreatedAt timezoner.DBTime `json:"created_at"`
}

func TestDBTime_ValueAndScan(t *testing.T) {
	t0 := time.Date(2026, 9, 1, 15, 30, 0, 0, time.UTC)
	dbTime := timezoner.NewDBTime(t0)

	// Test driver.Valuer
	val, err := dbTime.Value()
	if err != nil {
		t.Fatalf("Value() error: %v", err)
	}
	tVal, ok := val.(time.Time)
	if !ok || !tVal.Equal(t0) {
		t.Errorf("Value() esperado %v, obtenido: %v", t0, val)
	}

	// Test sql.Scanner con time.Time
	var scanned timezoner.DBTime
	if err := scanned.Scan(t0); err != nil {
		t.Fatalf("Scan(time.Time) error: %v", err)
	}
	if !scanned.Equal(t0) {
		t.Errorf("Scan esperado %v, obtenido %v", t0, scanned.Time)
	}

	// Test sql.Scanner con string ISO
	if err := scanned.Scan("2026-09-01T15:30:00Z"); err != nil {
		t.Fatalf("Scan(string) error: %v", err)
	}
	if !scanned.Equal(t0) {
		t.Errorf("Scan(string) esperado %v, obtenido %v", t0, scanned.Time)
	}

	// Test sql.Scanner con Unix int64
	if err := scanned.Scan(int64(1700000000)); err != nil {
		t.Fatalf("Scan(int64) error: %v", err)
	}

	// Test sql.Scanner con nil
	if err := scanned.Scan(nil); err != nil || !scanned.IsZero() {
		t.Errorf("Scan(nil) debería resetear a zero value")
	}
}

func TestDBTime_JSONRoundTrip(t *testing.T) {
	t0 := time.Date(2026, 9, 1, 15, 30, 0, 0, time.UTC)
	ev := EventRecord{
		ID:        101,
		Title:     "Lanzamiento Oficial",
		CreatedAt: timezoner.NewDBTime(t0),
	}

	data, err := json.Marshal(ev)
	if err != nil {
		t.Fatalf("json.Marshal falló: %v", err)
	}

	var decoded EventRecord
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal falló: %v", err)
	}

	if !decoded.CreatedAt.Equal(t0) {
		t.Errorf("JSON RoundTrip esperado: %v, obtenido: %v", t0, decoded.CreatedAt.Time)
	}
}
