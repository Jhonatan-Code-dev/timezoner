package types_test

import (
	"encoding/json"
	"testing"
	"time"

	"timezoner/pkg/types"
)

type Event struct {
	ID        int             `json:"id"`
	CreatedAt types.DBTime    `json:"created_at"`
	Scheduled types.ZonedTime `json:"scheduled"`
}

func TestTypes_DBTimeAndZonedTime(t *testing.T) {
	t0 := time.Date(2026, 9, 1, 15, 0, 0, 0, time.UTC)
	dbTime := types.NewDBTime(t0)

	// NowDBTime
	nowDB := types.NowDBTime()
	if nowDB.Location() != time.UTC {
		t.Errorf("NowDBTime debe ser UTC")
	}

	// Value
	val, err := dbTime.Value()
	if err != nil || val == nil {
		t.Fatalf("Value falló: %v", err)
	}

	// Value on Zero
	var zeroDB types.DBTime
	zVal, err := zeroDB.Value()
	if err != nil || zVal != nil {
		t.Errorf("Zero value debería retornar nil")
	}

	// Scan
	var scanned types.DBTime
	if err := scanned.Scan(t0); err != nil {
		t.Fatalf("Scan(time.Time) falló: %v", err)
	}
	if !scanned.Equal(t0) {
		t.Errorf("Scan time.Time no coincide")
	}

	if err := scanned.Scan("2026-09-01T15:00:00Z"); err != nil {
		t.Fatalf("Scan(string) falló: %v", err)
	}

	if err := scanned.Scan(int64(1700000000)); err != nil {
		t.Fatalf("Scan(int64) falló: %v", err)
	}

	if err := scanned.Scan(nil); err != nil || !scanned.IsZero() {
		t.Errorf("Scan(nil) falló")
	}

	if err := scanned.Scan(true); err == nil {
		t.Errorf("Scan(bool) debería fallar")
	}

	// JSON
	data, err := json.Marshal(dbTime)
	if err != nil {
		t.Fatalf("Marshal DBTime falló: %v", err)
	}
	var unmarshaled types.DBTime
	if err := json.Unmarshal(data, &unmarshaled); err != nil || !unmarshaled.Equal(t0) {
		t.Errorf("Unmarshal DBTime falló")
	}

	// JSON Zero
	zData, _ := json.Marshal(zeroDB)
	if string(zData) != "null" {
		t.Errorf("Zero DBTime en JSON debería ser null, obtenido: %s", zData)
	}

	// ZonedTime
	zoned, err := types.NewZonedTime(t0, "America/Lima")
	if err != nil {
		t.Fatalf("NewZonedTime falló: %v", err)
	}
	locTime, err := zoned.Local()
	if err != nil || locTime.Hour() != 10 {
		t.Errorf("Local() esperado 10, obtenido: %d", locTime.Hour())
	}

	zVal2, err := zoned.Value()
	if err != nil || zVal2 == nil {
		t.Fatalf("ZonedTime Value falló: %v", err)
	}

	var scannedZ types.ZonedTime
	if err := scannedZ.Scan(zVal2); err != nil || scannedZ.Zone != "America/Lima" {
		t.Errorf("Scan ZonedTime falló")
	}

	if err := scannedZ.Scan("2026-09-01T15:00:00Z|America/Lima"); err != nil || scannedZ.Zone != "America/Lima" {
		t.Errorf("Scan ZonedTime con pipe falló")
	}

	if err := scannedZ.Scan(t0); err != nil || scannedZ.Zone != "UTC" {
		t.Errorf("Scan ZonedTime time.Time falló")
	}

	if err := scannedZ.Scan(nil); err != nil || scannedZ.Zone != "" {
		t.Errorf("Scan ZonedTime nil falló")
	}

	// Error NewZonedTime
	if _, err := types.NewZonedTime(t0, "Invalid/Zone"); err == nil {
		t.Errorf("NewZonedTime con zona inválida debería fallar")
	}
}
