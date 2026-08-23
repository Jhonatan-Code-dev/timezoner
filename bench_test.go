package timezoner_test

import (
	"testing"
	"time"

	"timezoner"
)

var sinkTime time.Time
var sinkString string
var sinkErr error

// BenchmarkConvert mide la conversión de zona con cache en memoria.
func BenchmarkConvert(b *testing.B) {
	t := time.Date(2026, 9, 1, 15, 0, 0, 0, time.UTC)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkTime, sinkErr = timezoner.Convert(t, "America/Lima")
	}
}

// BenchmarkConvert_ColdCache mide la conversión de zonas poco frecuentes.
func BenchmarkConvert_ColdCache(b *testing.B) {
	zones := []string{
		"Pacific/Apia", "America/Caracas", "Asia/Kathmandu",
		"Pacific/Marquesas", "Africa/Monrovia",
	}
	t := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkTime, sinkErr = timezoner.Convert(t, zones[i%len(zones)])
	}
}

// BenchmarkIngestFromString mide el parsing de fechas con zona.
func BenchmarkIngestFromString(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkTime, sinkErr = timezoner.IngestFromString("2026-09-01 10:00", "America/Lima")
	}
}

// BenchmarkProjectForUser mide la proyección de UTC a zona de usuario.
func BenchmarkProjectForUser(b *testing.B) {
	utc := time.Date(2026, 9, 1, 15, 0, 0, 0, time.UTC)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, sinkErr = timezoner.ProjectForUser(utc, "Europe/Madrid")
	}
}

// BenchmarkAddBusinessDays mide la aritmética de días hábiles (incluyendo corrección DST).
func BenchmarkAddBusinessDays(b *testing.B) {
	t := time.Date(2026, 9, 4, 10, 0, 0, 0, time.UTC) // Viernes
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkTime = timezoner.AddBusinessDays(t, 5)
	}
}

// BenchmarkHumanize mide el formateo de tiempo relativo en español.
func BenchmarkHumanize(b *testing.B) {
	past := time.Now().Add(-2 * time.Hour)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkString = timezoner.Humanize(past)
	}
}

// BenchmarkFindOverlap mide el algoritmo de solapamiento de 3 zonas.
func BenchmarkFindOverlap(b *testing.B) {
	req := timezoner.OverlapRequest{
		Date:         time.Date(2026, 10, 15, 0, 0, 0, 0, time.UTC),
		Zones:        []string{"America/Lima", "America/New_York", "Europe/Madrid"},
		DefaultHours: timezoner.WorkingHours{StartHour: 9, EndHour: 18},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, sinkErr = timezoner.FindOverlap(req)
	}
}

// BenchmarkZonedFromLocal mide la ingesta de fecha local a ZonedTime.
func BenchmarkZonedFromLocal(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, sinkErr = timezoner.ZonedFromLocal("2026-10-01 10:00", "America/Lima")
	}
}

// BenchmarkNewDBTime mide la creación de DBTime (debe ser sub-nanosegundo).
func BenchmarkNewDBTime(b *testing.B) {
	t := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = timezoner.NewDBTime(t)
	}
}
