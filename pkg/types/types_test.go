package types_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Jhonatan-Code-dev/timezonermax/pkg/types"
)

func TestTypes_DBTime_Basics(t *testing.T) {
	t0 := time.Date(2026, 9, 1, 15, 0, 0, 123456789, time.UTC)
	dbTime := types.NewDBTime(t0)

	// Accessors controlados
	if !dbTime.Time().Equal(t0) {
		t.Errorf("DBTime.Time() debe ser igual al original")
	}
	if !dbTime.UTC().Equal(t0) {
		t.Errorf("DBTime.UTC() debe ser igual al original")
	}
	if dbTime.IsZero() {
		t.Errorf("DBTime no debe ser zero")
	}
	if dbTime.String() == "" {
		t.Errorf("DBTime.String() no debe ser vacío")
	}

	// Equal y EqualTime
	other := types.NewDBTime(t0.Add(time.Nanosecond))
	if dbTime.Equal(other) {
		t.Errorf("DBTime distintos no deben ser Equal")
	}
	if !dbTime.EqualTime(t0) {
		t.Errorf("EqualTime debe ser true con el mismo instante")
	}

	// NowDBTime
	nowDB := types.NowDBTime()
	if nowDB.IsZero() {
		t.Errorf("NowDBTime no debe ser zero")
	}
	if nowDB.UTC().Location() != time.UTC {
		t.Errorf("NowDBTime debe ser UTC")
	}
}

func TestTypes_DBTime_ValueAndScan(t *testing.T) {
	t0 := time.Date(2026, 9, 1, 15, 0, 0, 0, time.UTC)
	dbTime := types.NewDBTime(t0)

	// Value
	val, err := dbTime.Value()
	if err != nil || val == nil {
		t.Fatalf("Value() falló: %v", err)
	}

	// Value on Zero
	var zeroDB types.DBTime
	zVal, err := zeroDB.Value()
	if err != nil || zVal != nil {
		t.Errorf("Zero Value debe retornar nil")
	}

	// Scan time.Time
	var scanned types.DBTime
	if err := scanned.Scan(t0); err != nil {
		t.Fatalf("Scan(time.Time) falló: %v", err)
	}
	if !scanned.EqualTime(t0) {
		t.Errorf("Scan time.Time no coincide")
	}

	// Scan string RFC3339
	if err := scanned.Scan("2026-09-01T15:00:00Z"); err != nil {
		t.Fatalf("Scan(string) falló: %v", err)
	}

	// Scan []byte
	if err := scanned.Scan([]byte("2026-09-01T15:00:00Z")); err != nil {
		t.Fatalf("Scan([]byte) falló: %v", err)
	}

	// Scan int64
	if err := scanned.Scan(int64(1700000000)); err != nil {
		t.Fatalf("Scan(int64) falló: %v", err)
	}

	// Scan nil
	if err := scanned.Scan(nil); err != nil || !scanned.IsZero() {
		t.Errorf("Scan(nil) debe producir zero value")
	}

	// Scan tipo inválido
	if err := scanned.Scan(true); err == nil {
		t.Errorf("Scan(bool) debe fallar")
	}
}

func TestTypes_DBTime_JSON_RFC3339Nano(t *testing.T) {
	// RFC3339Nano preserva nanosegundos: defecto crítico corregido
	t0 := time.Date(2026, 9, 1, 15, 0, 0, 123456789, time.UTC)
	dbTime := types.NewDBTime(t0)

	data, err := json.Marshal(dbTime)
	if err != nil {
		t.Fatalf("Marshal DBTime falló: %v", err)
	}

	// El JSON serializado debe incluir los nanosegundos
	jsonStr := string(data)
	if jsonStr == "null" {
		t.Fatal("JSON no debe ser null para DBTime no zero")
	}

	var decoded types.DBTime
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal DBTime falló: %v", err)
	}
	if !decoded.EqualTime(t0) {
		t.Errorf("Round-trip JSON/DBTime falló: original=%v decoded=%v", t0, decoded.Time())
	}

	// JSON Zero
	var zeroDB types.DBTime
	zData, _ := json.Marshal(zeroDB)
	if string(zData) != "null" {
		t.Errorf("Zero DBTime en JSON debe ser null, obtenido: %s", zData)
	}
}

func TestTypes_ZonedTime(t *testing.T) {
	t0 := time.Date(2026, 9, 1, 15, 0, 0, 0, time.UTC)

	zoned, err := types.NewZonedTime(t0, "America/Lima")
	if err != nil {
		t.Fatalf("NewZonedTime falló: %v", err)
	}
	if zoned.Zone != "America/Lima" {
		t.Errorf("Zone esperado America/Lima, obtenido: %s", zoned.Zone)
	}

	// Local devuelve hora en Lima (UTC-5 en septiembre)
	locTime, err := zoned.Local()
	if err != nil || locTime.Hour() != 10 {
		t.Errorf("Local() esperado 10, obtenido: %d", locTime.Hour())
	}

	// IsZero
	if zoned.IsZero() {
		t.Errorf("ZonedTime no debe ser zero")
	}
	var z2 types.ZonedTime
	if !z2.IsZero() {
		t.Errorf("Zero ZonedTime debe ser zero")
	}

	// Value
	val, err := zoned.Value()
	if err != nil || val == nil {
		t.Fatalf("ZonedTime Value falló: %v", err)
	}

	// Scan desde Value
	var scannedZ types.ZonedTime
	if err := scannedZ.Scan(val); err != nil || scannedZ.Zone != "America/Lima" {
		t.Errorf("Scan ZonedTime desde Value falló")
	}

	// Scan pipe format
	if err := scannedZ.Scan("2026-09-01T15:00:00Z|America/Lima"); err != nil || scannedZ.Zone != "America/Lima" {
		t.Errorf("Scan ZonedTime con pipe falló")
	}

	// Scan time.Time
	if err := scannedZ.Scan(t0); err != nil || scannedZ.Zone != "UTC" {
		t.Errorf("Scan ZonedTime con time.Time falló")
	}

	// Scan nil
	if err := scannedZ.Scan(nil); err != nil || scannedZ.Zone != "" {
		t.Errorf("Scan ZonedTime nil falló")
	}

	// Error en zona inválida
	if _, err := types.NewZonedTime(t0, "Invalid/Zone"); err == nil {
		t.Errorf("NewZonedTime con zona inválida debe fallar")
	}

	// Scan JSON con campo utc corrupto
	if err := scannedZ.Scan(`{"utc": "NOT_A_DATE", "zone": "America/Lima"}`); err == nil {
		t.Errorf("Scan ZonedTime con utc corrupto debe retornar error")
	}

	// Scan tipo inválido
	if err := scannedZ.Scan(42); err == nil {
		t.Errorf("Scan ZonedTime con int debe fallar")
	}
}
