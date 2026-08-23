package timezoner_test

import (
	"testing"
	"time"

	"timezoner"
)

func TestCalendar_WeekdaysAndBusinessDays(t *testing.T) {
	// Viernes 4 de Septiembre de 2026
	friday := time.Date(2026, 9, 4, 10, 0, 0, 0, time.UTC)
	if timezoner.IsWeekend(friday) || !timezoner.IsWeekday(friday) {
		t.Errorf("Viernes debería ser Weekday y no Weekend")
	}

	// Sábado 5 de Septiembre de 2026
	saturday := time.Date(2026, 9, 5, 10, 0, 0, 0, time.UTC)
	if !timezoner.IsWeekend(saturday) || timezoner.IsWeekday(saturday) {
		t.Errorf("Sábado debería ser Weekend y no Weekday")
	}

	// Sumar 1 día hábil a un Viernes debe dar Lunes (7 de Septiembre)
	nextBusinessDay := timezoner.AddBusinessDays(friday, 1)
	if nextBusinessDay.Day() != 7 || nextBusinessDay.Weekday() != time.Monday {
		t.Errorf("AddBusinessDays(Viernes, 1) esperado Lunes 7, obtenido día %d (%v)", nextBusinessDay.Day(), nextBusinessDay.Weekday())
	}

	// Restar 1 día hábil a un Lunes debe dar Viernes (4 de Septiembre)
	monday := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	prevBusinessDay := timezoner.AddBusinessDays(monday, -1)
	if prevBusinessDay.Day() != 4 || prevBusinessDay.Weekday() != time.Friday {
		t.Errorf("AddBusinessDays(Lunes, -1) esperado Viernes 4, obtenido día %d (%v)", prevBusinessDay.Day(), prevBusinessDay.Weekday())
	}
}

func TestCalendar_Bounds(t *testing.T) {
	t0 := time.Date(2026, 9, 15, 14, 30, 45, 123456789, time.UTC)

	startDay := timezoner.StartOfDay(t0)
	if startDay.Hour() != 0 || startDay.Minute() != 0 || startDay.Second() != 0 || startDay.Nanosecond() != 0 {
		t.Errorf("StartOfDay esperado 00:00:00.0, obtenido %v", startDay)
	}

	endDay := timezoner.EndOfDay(t0)
	if endDay.Hour() != 23 || endDay.Minute() != 59 || endDay.Second() != 59 {
		t.Errorf("EndOfDay esperado 23:59:59, obtenido %v", endDay)
	}

	startMonth := timezoner.StartOfMonth(t0)
	if startMonth.Day() != 1 || startMonth.Hour() != 0 {
		t.Errorf("StartOfMonth esperado día 1 a las 00:00, obtenido %v", startMonth)
	}

	endMonth := timezoner.EndOfMonth(t0)
	if endMonth.Day() != 30 || endMonth.Hour() != 23 {
		t.Errorf("EndOfMonth esperado día 30 a las 23h para septiembre, obtenido %v", endMonth)
	}

	// Comprobar año bisiesto 2028 en febrero
	daysFebLeap := timezoner.DaysInMonth(2028, time.February)
	if daysFebLeap != 29 {
		t.Errorf("Días en Feb 2028 esperado 29, obtenido %d", daysFebLeap)
	}
	daysFebNormal := timezoner.DaysInMonth(2026, time.February)
	if daysFebNormal != 28 {
		t.Errorf("Días en Feb 2026 esperado 28, obtenido %d", daysFebNormal)
	}
}

func TestFluent_CalendarMethods(t *testing.T) {
	friday := time.Date(2026, 9, 4, 10, 0, 0, 0, time.UTC)
	monday := timezoner.At(friday).
		AddBusinessDays(1).
		StartOfDay().
		MustTime()

	if monday.Day() != 7 || monday.Hour() != 0 {
		t.Errorf("Fluent AddBusinessDays + StartOfDay falló, obtenido: %v", monday)
	}
}
