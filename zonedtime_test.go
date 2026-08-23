package timezoner_test

import (
	"encoding/json"
	"testing"
	"time"

	"timezoner"
)

type AppointmentRecord struct {
	ID          int                 `json:"id"`
	Patient     string              `json:"patient"`
	ScheduledAt timezoner.ZonedTime `json:"scheduled_at"`
}

func TestZonedTime_CreationAndProjection(t *testing.T) {
	// Cita programada en Lima: 2026-09-01 10:00 AM
	zoned, err := timezoner.ZonedFromLocal("2026-09-01 10:00", "America/Lima")
	if err != nil {
		t.Fatalf("ZonedFromLocal falló: %v", err)
	}

	if zoned.Zone != "America/Lima" {
		t.Errorf("Zona esperada 'America/Lima', obtenida: '%s'", zoned.Zone)
	}
	if zoned.UTC.Hour() != 15 {
		t.Errorf("Hora UTC esperada 15, obtenida: %d", zoned.UTC.Hour())
	}

	// Recalcular hora local original
	localTime, err := zoned.Local()
	if err != nil {
		t.Fatalf("Local() falló: %v", err)
	}
	if localTime.Hour() != 10 {
		t.Errorf("Hora local esperada 10, obtenida: %d", localTime.Hour())
	}

	// Proyectar para un médico o colega en Madrid (17:00 CEST)
	madridView, err := zoned.ForViewer("Europe/Madrid")
	if err != nil {
		t.Fatalf("ForViewer falló: %v", err)
	}
	if madridView.LocalTime.Hour() != 17 {
		t.Errorf("Hora en Madrid esperada 17, obtenida: %d", madridView.LocalTime.Hour())
	}
}

func TestZonedTime_DatabaseAndJSON(t *testing.T) {
	zoned, err := timezoner.ZonedFromLocal("2026-09-01 10:00", "PET")
	if err != nil {
		t.Fatalf("ZonedFromLocal falló: %v", err)
	}

	app := AppointmentRecord{
		ID:          501,
		Patient:     "Carlos Perez",
		ScheduledAt: zoned,
	}

	// JSON Round-trip
	data, err := json.Marshal(app)
	if err != nil {
		t.Fatalf("json.Marshal falló: %v", err)
	}

	var decoded AppointmentRecord
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal falló: %v", err)
	}

	if decoded.ScheduledAt.Zone != "America/Lima" {
		t.Errorf("Zona decodificada esperada America/Lima, obtenida: %s", decoded.ScheduledAt.Zone)
	}
	if !decoded.ScheduledAt.UTC.Equal(zoned.UTC.Time) {
		t.Errorf("Instante UTC no coincide tras JSON Roundtrip")
	}

	// SQL Valuer & Scanner (formato JSON)
	val, err := zoned.Value()
	if err != nil {
		t.Fatalf("Value() falló: %v", err)
	}

	var scanned timezoner.ZonedTime
	if err := scanned.Scan(val); err != nil {
		t.Fatalf("Scan(string JSON) falló: %v", err)
	}
	if scanned.Zone != "America/Lima" || !scanned.UTC.Equal(zoned.UTC.Time) {
		t.Errorf("Scan falló en reconstruir ZonedTime correctamente")
	}

	// SQL Scanner con formato compuesto "2026-09-01T15:00:00Z|America/Lima"
	if err := scanned.Scan("2026-09-01T15:00:00Z|America/Lima"); err != nil {
		t.Fatalf("Scan(formato compuesto) falló: %v", err)
	}
	if scanned.Zone != "America/Lima" || scanned.UTC.Hour() != 15 {
		t.Errorf("Scan con formato compuesto falló")
	}

	// SQL Scanner con nil
	if err := scanned.Scan(nil); err != nil || scanned.Zone != "" {
		t.Errorf("Scan(nil) debería dejar ZonedTime en zero-value")
	}
}

func TestNewZonedTime_Direct(t *testing.T) {
	t0 := time.Date(2026, 9, 1, 15, 0, 0, 0, time.UTC)
	zoned, err := timezoner.NewZonedTime(t0, "Asia/Tokyo")
	if err != nil {
		t.Fatalf("NewZonedTime falló: %v", err)
	}

	if zoned.Zone != "Asia/Tokyo" {
		t.Errorf("Zona esperada Asia/Tokyo, obtenida: %s", zoned.Zone)
	}

	// Error con zona inválida
	_, err = timezoner.NewZonedTime(t0, "Invalid/Zone")
	if err == nil {
		t.Errorf("Se esperaba error al pasar zona inválida a NewZonedTime")
	}
}
