package timezoner_test

import (
	"testing"
	"time"

	"timezoner"
)

func TestTimezoner_FacadeE2E(t *testing.T) {
	// 1. Ingesta: string en zona Lima → UTC
	dbTime, err := timezoner.IngestFromString("2026-09-01 10:00", "America/Lima")
	if err != nil {
		t.Fatalf("IngestFromString falló: %v", err)
	}
	if dbTime.Hour() != 15 {
		t.Errorf("Hora esperada 15 UTC, obtenida: %d", dbTime.Hour())
	}

	// 2. Proyección Madrid (UTC+2 en verano → 17:00)
	userMadrid, err := timezoner.ProjectForUser(dbTime, "Europe/Madrid")
	if err != nil {
		t.Fatalf("ProjectForUser falló: %v", err)
	}
	if userMadrid.LocalTime.Hour() != 17 {
		t.Errorf("Hora Madrid esperada 17, obtenida: %d", userMadrid.LocalTime.Hour())
	}

	// 3. Fluent API: días hábiles, límites de mes y día
	dueDate := timezoner.At(dbTime).
		AddBusinessDays(3).
		StartOfMonth().
		EndOfMonth().
		EndOfDay().
		MustTime()

	if dueDate.Hour() != 23 || dueDate.Minute() != 59 {
		t.Errorf("EndOfDay falló, obtenido: %v", dueDate)
	}

	// 4. StartOfWeek / EndOfWeek
	monday := timezoner.At(dbTime).StartOfWeek().MustTime()
	if monday.Weekday() != time.Monday {
		t.Errorf("StartOfWeek debe ser lunes, obtenido: %v", monday.Weekday())
	}
	sunday := timezoner.At(dbTime).EndOfWeek().MustTime()
	if sunday.Weekday() != time.Sunday {
		t.Errorf("EndOfWeek debe ser domingo, obtenido: %v", sunday.Weekday())
	}

	// 5. Humanize
	rel := timezoner.Humanize(dbTime, dbTime.Add(2*time.Hour))
	if rel != "hace 2 horas" {
		t.Errorf("Humanize esperado 'hace 2 horas', obtenido: %q", rel)
	}
	relEn, _ := timezoner.At(dbTime).HumanizeEn(dbTime.Add(2 * time.Hour))
	if relEn != "2 hours ago" {
		t.Errorf("HumanizeEn esperado '2 hours ago', obtenido: %q", relEn)
	}

	// 6. ZonedFromLocal
	zoned, err := timezoner.ZonedFromLocal("2026-10-01 10:00", "America/Lima")
	if err != nil {
		t.Fatalf("ZonedFromLocal falló: %v", err)
	}
	localTime, err := zoned.Local()
	if err != nil || localTime.Hour() != 10 {
		t.Errorf("ZonedTime Local() esperado 10, obtenido: %d", localTime.Hour())
	}

	// 7. Conversiones directas
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

	// 8. Compare (ZoneSnapshot con TimeUTC determinista)
	snapshots, err := timezoner.Compare(dbTime, "America/Lima", "Asia/Tokyo")
	if err != nil || len(snapshots) != 2 {
		t.Errorf("Compare falló")
	}
	if snapshots[0].TimeUTC == "" {
		t.Errorf("Snapshot.TimeUTC no debe ser vacío")
	}

	// 9. FindOverlap
	overlapSlots, err := timezoner.FindOverlap(timezoner.OverlapRequest{
		Date:  dbTime,
		Zones: []string{"America/Lima", "Europe/Madrid"},
	})
	if err != nil || len(overlapSlots) == 0 {
		t.Errorf("FindOverlap falló")
	}

	// 10. Calendar helpers
	if timezoner.DaysInMonth(2028, time.February) != 29 {
		t.Errorf("DaysInMonth bisiesto falló")
	}
	if !timezoner.IsWeekday(dbTime) || timezoner.IsWeekend(dbTime) {
		t.Errorf("IsWeekday / IsWeekend falló")
	}
	if timezoner.StartOfDay(dbTime).Hour() != 0 {
		t.Errorf("StartOfDay falló")
	}

	// 11. SupportedLayouts es función (no variable mutable)
	layouts := timezoner.SupportedLayouts()
	if len(layouts) == 0 {
		t.Errorf("SupportedLayouts no debe estar vacío")
	}

	// 12. Fluent API — AsDBTime, AsZonedTime, Err()
	tp := timezoner.Now().In("America/Lima").ToUTC()
	if tp.Err() != nil {
		t.Errorf("TimePoint no debe tener error: %v", tp.Err())
	}
	_ = tp.MustDBTime()
	_ = tp.MustZonedTime("America/Lima")
	_, _ = tp.Format("")
	_, _ = tp.Info()
	_, _ = tp.AsDBTime()
	_, _ = tp.AsZonedTime("America/Lima")
}

func TestTimezoner_ConcurrentStress(t *testing.T) {
	const goroutines = 200
	done := make(chan bool, goroutines)
	now := time.Now()

	// Verificación de RegisterAlias concurrente (race condition corregido)
	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			_, _ = timezoner.Convert(now, "America/Lima")
			_, _ = timezoner.GetZoneInfo("Europe/Madrid", now)
			_ = timezoner.At(now).AddBusinessDays(1).MustTime()
			_ = timezoner.RegisterAlias("TESTZONE", "America/New_York")
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
	f.Add("", int64(-1))
	f.Add("PET", int64(9223372036))

	f.Fuzz(func(t *testing.T, zone string, unixSec int64) {
		tm := time.Unix(unixSec, 0)
		_, _ = timezoner.Convert(tm, zone)
		_, _ = timezoner.GetZoneInfo(zone, tm)
		_, _ = timezoner.NowIn(zone)
	})
}

func FuzzIngestFromString(f *testing.F) {
	f.Add("2026-09-01 10:00", "America/Lima")
	f.Add("", "UTC")
	f.Add("not-a-date", "Invalid/Zone")

	f.Fuzz(func(t *testing.T, dateStr, zone string) {
		// Nunca debe producir pánico, siempre retorna error o resultado válido
		_, _ = timezoner.IngestFromString(dateStr, zone)
	})
}
