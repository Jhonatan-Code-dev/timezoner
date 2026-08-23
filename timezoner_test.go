package timezoner

import (
	"errors"
	"sync"
	"testing"
	"time"
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
		{"", true, "", ErrEmptyZoneName},
		{"Invalid/NonExistent_Zone", true, "", ErrInvalidZone},
	}

	for _, tc := range tests {
		loc, err := LoadLocation(tc.input)
		if tc.expectError {
			if err == nil {
				t.Errorf("LoadLocation(%q): se esperaba error y no ocurrió", tc.input)
			} else if tc.expectedErr != nil && !errors.Is(err, tc.expectedErr) {
				t.Errorf("LoadLocation(%q) err = %v, se esperaba errors.Is(%v)", tc.input, err, tc.expectedErr)
			}
		} else {
			if err != nil {
				t.Errorf("LoadLocation(%q): error inesperado: %v", tc.input, err)
			} else if loc.String() != tc.expectedLoc {
				t.Errorf("LoadLocation(%q) = %v; se esperaba %v", tc.input, loc.String(), tc.expectedLoc)
			}
		}
	}
}

func TestIsValidAndNormalize(t *testing.T) {
	if !IsValid("America/Bogota") {
		t.Errorf("IsValid(America/Bogota) debería ser true")
	}
	if !IsValid("UTC") {
		t.Errorf("IsValid(UTC) debería ser true")
	}
	if IsValid("Mars/Olympus_Mons") {
		t.Errorf("IsValid(Mars/Olympus_Mons) debería ser false")
	}

	norm, err := NormalizeZone("COT")
	if err != nil || norm != "America/Bogota" {
		t.Errorf("NormalizeZone(COT) = %s, err: %v; esperado America/Bogota", norm, err)
	}
}

func TestRegisterAlias(t *testing.T) {
	err := RegisterAlias("MIZONA", "America/Lima")
	if err != nil {
		t.Fatalf("RegisterAlias falló: %v", err)
	}

	loc, err := LoadLocation("MIZONA")
	if err != nil || loc.String() != "America/Lima" {
		t.Errorf("LoadLocation(MIZONA) = %v, %v; esperado America/Lima", loc, err)
	}

	// Error con zona destino inválida
	err = RegisterAlias("BAD", "Zone/Fake")
	if err == nil || !errors.Is(err, ErrInvalidZone) {
		t.Errorf("RegisterAlias con zona inválida debería haber fallado con ErrInvalidZone")
	}
}

func TestCommonZones(t *testing.T) {
	zones := CommonZones()
	if len(zones) == 0 {
		t.Errorf("CommonZones no debería estar vacío")
	}
}

func TestConvert(t *testing.T) {
	utcTime := time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC)

	// Convertir a Lima (UTC-5 sin DST)
	limaTime, err := Convert(utcTime, "America/Lima")
	if err != nil {
		t.Fatalf("Convert falló: %v", err)
	}
	if limaTime.Hour() != 7 {
		t.Errorf("Hora esperada en Lima: 7, obtenida: %d", limaTime.Hour())
	}

	// Convertir a Tokio (UTC+9)
	tokyoTime, err := Convert(utcTime, "Asia/Tokyo")
	if err != nil {
		t.Fatalf("Convert a Tokio falló: %v", err)
	}
	if tokyoTime.Hour() != 21 {
		t.Errorf("Hora esperada en Tokio: 21, obtenida: %d", tokyoTime.Hour())
	}
}

func TestConvertBetween(t *testing.T) {
	layout := "2006-01-02 15:04"
	timeStr := "2026-05-10 10:00"

	res, err := ConvertBetween(timeStr, layout, "America/Lima", "Europe/Madrid")
	if err != nil {
		t.Fatalf("ConvertBetween falló: %v", err)
	}

	if res.Hour() != 17 {
		t.Errorf("ConvertBetween hora esperada: 17, obtenida: %d", res.Hour())
	}

	// Probar error con layout incompatible
	_, err = ConvertBetween("invalid-date", layout, "UTC", "America/Lima")
	if err == nil || !errors.Is(err, ErrInvalidTimeFormat) {
		t.Errorf("Se esperaba ErrInvalidTimeFormat al parsear fecha inválida, obtenido: %v", err)
	}
}

func TestFormatIn(t *testing.T) {
	utcTime := time.Date(2026, 1, 1, 15, 30, 0, 0, time.UTC)
	formatted, err := FormatIn(utcTime, "America/Lima", "15:04 MST")
	if err != nil {
		t.Fatalf("FormatIn falló: %v", err)
	}
	if formatted != "10:30 -05" && formatted != "10:30 PET" {
		t.Logf("FormatIn retornó: %s", formatted)
	}
}

func TestGetZoneInfoAndDifference(t *testing.T) {
	baseTime := time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)

	infoLima, err := GetZoneInfo("America/Lima", baseTime)
	if err != nil {
		t.Fatalf("GetZoneInfo(Lima) falló: %v", err)
	}
	if infoLima.OffsetSeconds != -5*3600 {
		t.Errorf("Offset de Lima esperado: %d, obtenido: %d", -5*3600, infoLima.OffsetSeconds)
	}
	if infoLima.OffsetFormatted != "-05:00" {
		t.Errorf("OffsetFormatted esperado -05:00, obtenido: %s", infoLima.OffsetFormatted)
	}

	infoMadrid, err := GetZoneInfo("Europe/Madrid", baseTime)
	if err != nil {
		t.Fatalf("GetZoneInfo(Madrid) falló: %v", err)
	}
	if !infoMadrid.IsDST {
		t.Errorf("Madrid en julio debería tener IsDST = true")
	}

	diff, err := Difference("Europe/Madrid", "America/Lima", baseTime)
	if err != nil {
		t.Fatalf("Difference falló: %v", err)
	}
	if diff != 7*time.Hour {
		t.Errorf("Diferencia esperada: 7h, obtenida: %v", diff)
	}
}

func TestCompare(t *testing.T) {
	base := time.Date(2026, 3, 1, 14, 0, 0, 0, time.UTC)
	snaps, err := Compare(base, "America/Lima", "Europe/London", "Asia/Tokyo")
	if err != nil {
		t.Fatalf("Compare falló: %v", err)
	}
	if len(snaps) != 3 {
		t.Fatalf("Se esperaban 3 snapshots, obtenidos %d", len(snaps))
	}

	// Test error con lista vacía
	_, err = Compare(base)
	if err == nil || !errors.Is(err, ErrNoZonesProvided) {
		t.Errorf("Compare sin zonas debería fallar con ErrNoZonesProvided")
	}
}

func TestFindOverlap(t *testing.T) {
	date := time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC)
	req := OverlapRequest{
		Date:         date,
		Zones:        []string{"America/Lima", "America/New_York"},
		DefaultHours: WorkingHours{StartHour: 9, EndHour: 17},
		SlotDuration: 1 * time.Hour,
	}

	slots, err := FindOverlap(req)
	if err != nil {
		t.Fatalf("FindOverlap falló: %v", err)
	}
	if len(slots) == 0 {
		t.Fatalf("Se esperaba encontrar solapamiento entre Lima y New York")
	}

	// Test error sin zonas
	_, err = FindOverlap(OverlapRequest{})
	if err == nil || !errors.Is(err, ErrNoZonesProvided) {
		t.Errorf("FindOverlap sin zonas debería fallar con ErrNoZonesProvided")
	}
}

func TestEdgeCasesTimezones(t *testing.T) {
	// Zona con 30 minutos de offset: India (UTC+05:30)
	t0 := time.Date(2028, 2, 29, 12, 0, 0, 0, time.UTC) // Año bisiesto
	indiaTime, err := Convert(t0, "Asia/Kolkata")
	if err != nil {
		t.Fatalf("Convert Kolkata falló: %v", err)
	}
	if indiaTime.Hour() != 17 || indiaTime.Minute() != 30 {
		t.Errorf("Hora esperada en India: 17:30, obtenida: %02d:%02d", indiaTime.Hour(), indiaTime.Minute())
	}

	// Zona con 45 minutos de offset: Nepal (UTC+05:45)
	nepalTime, err := Convert(t0, "Asia/Kathmandu")
	if err != nil {
		t.Fatalf("Convert Kathmandu falló: %v", err)
	}
	if nepalTime.Hour() != 17 || nepalTime.Minute() != 45 {
		t.Errorf("Hora esperada en Nepal: 17:45, obtenida: %02d:%02d", nepalTime.Hour(), nepalTime.Minute())
	}
}

func TestConcurrentLoadAndConvert(t *testing.T) {
	const goroutines = 200
	var wg sync.WaitGroup
	wg.Add(goroutines)

	zones := []string{"America/Lima", "Europe/Madrid", "Asia/Tokyo", "UTC", "America/New_York", "PET", "BST", "JST"}
	now := time.Now()

	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			z := zones[idx%len(zones)]
			_, err := Convert(now, z)
			if err != nil {
				t.Errorf("Concurrent Convert error en goroutine %d para zona %s: %v", idx, z, err)
			}
			_, _ = GetZoneInfo(z, now)
		}(i)
	}

	wg.Wait()
}

func TestFluentAPI(t *testing.T) {
	t0 := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
	tp := At(t0).In("America/Lima")

	gotTime, err := tp.Time()
	if err != nil {
		t.Fatalf("tp.Time() error: %v", err)
	}
	if gotTime.Hour() != 7 {
		t.Errorf("Hora esperada 7, obtenida %d", gotTime.Hour())
	}

	mustTime := tp.MustTime()
	if mustTime.Hour() != 7 {
		t.Errorf("MustTime hora esperada 7, obtenida %d", mustTime.Hour())
	}

	info, err := tp.Info()
	if err != nil || info.OffsetFormatted != "-05:00" {
		t.Errorf("Info offset esperado -05:00, obtenido: %s (err: %v)", info.OffsetFormatted, err)
	}

	formatted, err := tp.Format("15:04")
	if err != nil || formatted != "07:00" {
		t.Errorf("Format esperado 07:00, obtenido: %s", formatted)
	}

	// Test error chaining
	badTp := At(t0).In("Zone/Invalid")
	if _, err := badTp.Time(); err == nil {
		t.Errorf("Se esperaba error al usar zona inválida")
	}
}

// Fuzzing
func FuzzConvert(f *testing.F) {
	// Seed inputs
	f.Add("UTC", int64(1700000000))
	f.Add("America/Lima", int64(1800000000))
	f.Add("Invalid/Zone", int64(0))
	f.Add("", int64(-100000))
	f.Add("Asia/Tokyo", int64(2000000000))

	f.Fuzz(func(t *testing.T, zone string, unixSec int64) {
		tm := time.Unix(unixSec, 0)
		_, _ = Convert(tm, zone)
		_, _ = GetZoneInfo(zone, tm)
	})
}

// Benchmarks
func BenchmarkConvert(b *testing.B) {
	t0 := time.Now()
	for i := 0; i < b.N; i++ {
		_, _ = Convert(t0, "America/Lima")
	}
}

func BenchmarkDifference(b *testing.B) {
	t0 := time.Now()
	for i := 0; i < b.N; i++ {
		_, _ = Difference("Europe/Madrid", "America/Lima", t0)
	}
}

func BenchmarkFindOverlap(b *testing.B) {
	date := time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC)
	req := OverlapRequest{
		Date:         date,
		Zones:        []string{"America/Lima", "Europe/Madrid", "America/New_York"},
		DefaultHours: WorkingHours{StartHour: 9, EndHour: 18},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FindOverlap(req)
	}
}
