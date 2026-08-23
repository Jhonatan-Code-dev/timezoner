package calendar_test

import (
	"testing"
	"time"

	"timezoner/pkg/calendar"
)

func TestCalendar_WeekdaysAndBusinessDays(t *testing.T) {
	friday := time.Date(2026, 9, 4, 10, 0, 0, 0, time.UTC)
	if calendar.IsWeekend(friday) || !calendar.IsWeekday(friday) {
		t.Errorf("Viernes debería ser Weekday")
	}

	monday := calendar.AddBusinessDays(friday, 1)
	if monday.Day() != 7 || monday.Weekday() != time.Monday {
		t.Errorf("AddBusinessDays(Viernes, 1) esperado Lunes 7, obtenido día %d", monday.Day())
	}
}

func TestCalendar_Bounds(t *testing.T) {
	t0 := time.Date(2026, 9, 15, 14, 30, 45, 0, time.UTC)
	startDay := calendar.StartOfDay(t0)
	if startDay.Hour() != 0 || startDay.Minute() != 0 {
		t.Errorf("StartOfDay esperado 00:00, obtenido %v", startDay)
	}

	endDay := calendar.EndOfDay(t0)
	if endDay.Hour() != 23 || endDay.Minute() != 59 {
		t.Errorf("EndOfDay esperado 23:59, obtenido %v", endDay)
	}

	daysFebLeap := calendar.DaysInMonth(2028, time.February)
	if daysFebLeap != 29 {
		t.Errorf("Días en Feb 2028 esperado 29, obtenido %d", daysFebLeap)
	}
}
