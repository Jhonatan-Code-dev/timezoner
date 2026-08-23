package timezoner_test

import (
	"encoding/json"
	"testing"
	"time"

	"timezoner"
)

func TestTimezoner_FacadeE2E(t *testing.T) {
	// 1. Ingesta
	dbTime, err := timezoner.IngestFromString("2026-09-01 10:00", "America/Lima")
	if err != nil {
		t.Fatalf("IngestFromString falló: %v", err)
	}
	if dbTime.Hour() != 15 {
		t.Errorf("Hora esperada 15, obtenida: %d", dbTime.Hour())
	}

	// 2. Proyección
	userMadrid, err := timezoner.ProjectForUser(dbTime, "Europe/Madrid")
	if err != nil {
		t.Fatalf("ProjectForUser falló: %v", err)
	}
	if userMadrid.LocalTime.Hour() != 17 {
		t.Errorf("Hora en Madrid esperada 17, obtenida: %d", userMadrid.LocalTime.Hour())
	}

	// 3. Aritmética de días hábiles y límites
	dueDate := timezoner.At(dbTime).
		AddBusinessDays(3).
		StartOfMonth().
		EndOfMonth().
		EndOfDay().
		MustTime()

	if dueDate.Hour() != 23 || dueDate.Minute() != 59 {
		t.Errorf("EndOfDay falló, obtenido: %v", dueDate)
	}

	// 4. Humanize
	rel := timezoner.Humanize(dbTime, dbTime.Add(2*time.Hour))
	if rel != "hace 2 horas" {
		t.Errorf("Humanize esperado 'hace 2 horas', obtenido: %q", rel)
	}
	relEn, _ := timezoner.At(dbTime).HumanizeEn(dbTime.Add(2 * time.Hour))
	if relEn != "2 hours ago" {
		t.Errorf("HumanizeEn esperado '2 hours ago', obtenido: %q", relEn)
	}

	// 5. ZonedTime
	zoned, err := timezoner.ZonedFromLocal("2026-10-01 10:00", "America/Lima")
	if err != nil {
		t.Fatalf("ZonedFromLocal falló: %v", err)
	}
	localTime, err := zoned.Local()
	if err != nil || localTime.Hour() != 10 {
		t.Errorf("ZonedTime Local() esperado 10, obtenido: %v", localTime.Hour())
	}

	// JSON Round-trip
	data, err := json.Marshal(zoned)
	if err != nil {
		t.Fatalf("Marshal falló: %v", err)
	}
	var decoded timezoner.ZonedTime
	if err := json.Unmarshal(data, &decoded); err != nil || decoded.Zone != "America/Lima" {
		t.Errorf("Unmarshal ZonedTime falló")
	}

	// 6. Conversiones directas y format
	nowInTokyo, err := timezoner.NowIn("Asia/Tokyo")
	if err != nil || nowInTokyo.IsZero() {
		t.Errorf("NowIn Tokyo falló")
	}

	formatted, err := timezoner.FormatIn(dbTime, "Europe/Paris", "15:04")
	if err != nil || formatted != "17:00" {
		t.Errorf("FormatIn falló, obtenido: %s", formatted)
	}

	convertedBetween, err := timezoner.ConvertBetween("2026-09-01 10:00", "2006-01-02 15:04", "America/Lima", "Europe/Madrid")
	if err != nil || convertedBetween.Hour() != 17 {
		t.Errorf("ConvertBetween falló")
	}

	// 7. Compare
	snapshots, err := timezoner.Compare(dbTime, "America/Lima", "Asia/Tokyo")
	if err != nil || len(snapshots) != 2 {
		t.Errorf("Compare falló")
	}

	// 8. Overlap
	overlapSlots, err := timezoner.FindOverlap(timezoner.OverlapRequest{
		Date:  dbTime,
		Zones: []string{"America/Lima", "Europe/Madrid"},
	})
	if err != nil || len(overlapSlots) == 0 {
		t.Errorf("FindOverlap falló")
	}

	// 9. Calendar bounds y verificación
	if timezoner.DaysInMonth(2028, time.February) != 29 {
		t.Errorf("DaysInMonth en bisiesto falló")
	}
	if !timezoner.IsWeekday(dbTime) || timezoner.IsWeekend(dbTime) {
		t.Errorf("IsWeekday / IsWeekend falló")
	}
	if timezoner.StartOfDay(dbTime).Hour() != 0 {
		t.Errorf("StartOfDay falló")
	}

	// 10. Fluent API Must variants y DBTime
	tp := timezoner.Now().In("America/Lima").ToUTC()
	_ = tp.MustDBTime()
	_ = tp.MustZonedTime("America/Lima")
	_, _ = tp.Format("")
	_, _ = tp.Info()
}

func TestTimezoner_ConcurrentStress(t *testing.T) {
	const goroutines = 100
	done := make(chan bool, goroutines)
	now := time.Now()

	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			_, _ = timezoner.Convert(now, "America/Lima")
			_, _ = timezoner.GetZoneInfo("Europe/Madrid", now)
			_ = timezoner.At(now).AddBusinessDays(1).MustTime()
			done <- true
		}(i)
	}

	for i := 0; i < goroutines; i++ {
		<-done
	}
}

func FuzzConvert(f *testing.F) {
	f.Add("UTC", int64(1700000000))
	f.Add("America/Lima", int64(1800000000))
	f.Add("Invalid/Zone", int64(0))

	f.Fuzz(func(t *testing.T, zone string, unixSec int64) {
		tm := time.Unix(unixSec, 0)
		_, _ = timezoner.Convert(tm, zone)
		_, _ = timezoner.GetZoneInfo(zone, tm)
	})
}
