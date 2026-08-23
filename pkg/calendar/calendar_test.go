package calendar_test

import (
	"testing"
	"time"

	"github.com/Jhonatan-Code-dev/timezoner/pkg/calendar"
)

func TestCalendar_WeekdaysAndBusinessDays(t *testing.T) {
	friday := time.Date(2026, 9, 4, 10, 0, 0, 0, time.UTC)
	if calendar.IsWeekend(friday) || !calendar.IsWeekday(friday) {
		t.Errorf("Viernes debería ser Weekday")
	}
	saturday := friday.AddDate(0, 0, 1)
	if !calendar.IsWeekend(saturday) || calendar.IsWeekday(saturday) {
		t.Errorf("Sábado debería ser Weekend")
	}

	// AddBusinessDays positivo: Viernes + 1 = Lunes
	monday := calendar.AddBusinessDays(friday, 1)
	if monday.Weekday() != time.Monday {
		t.Errorf("AddBusinessDays(viernes, 1) esperado lunes, obtenido %v", monday.Weekday())
	}

	// AddBusinessDays negativo: Lunes - 1 = Viernes
	prevFriday := calendar.AddBusinessDays(monday, -1)
	if prevFriday.Weekday() != time.Friday {
		t.Errorf("AddBusinessDays(lunes, -1) esperado viernes, obtenido %v", prevFriday.Weekday())
	}

	// AddBusinessDays = 0
	same := calendar.AddBusinessDays(friday, 0)
	if !same.Equal(friday) {
		t.Errorf("AddBusinessDays(t, 0) debe retornar el mismo instante")
	}

	// Saltar fin de semana
	next5 := calendar.AddBusinessDays(friday, 5) // Viernes + 5 = viernes siguiente
	if next5.Weekday() != time.Friday {
		t.Errorf("Viernes + 5 días hábiles debería ser otro viernes, obtenido %v", next5.Weekday())
	}
}

func TestCalendar_DST_Preservation(t *testing.T) {
	// Transición DST en New York: 8 de marzo de 2026 a las 02:00 → 03:00
	loc, _ := time.LoadLocation("America/New_York")
	// Un viernes a las 09:00 antes del cambio
	preSwitch := time.Date(2026, 3, 6, 9, 0, 0, 0, loc) // viernes

	// 2 días hábiles después debe ser el martes (9 marzo lunes, 10 marzo martes)
	// y la hora local debe seguir siendo las 09:00, no 10:00 por el cambio DST
	result := calendar.AddBusinessDays(preSwitch, 2)
	if result.Hour() != 9 {
		t.Errorf("AddBusinessDays debe preservar hora local tras DST, esperado 9, obtenido %d", result.Hour())
	}
}

func TestCalendar_Bounds(t *testing.T) {
	t0 := time.Date(2026, 9, 15, 14, 30, 45, 0, time.UTC)

	startDay := calendar.StartOfDay(t0)
	if startDay.Hour() != 0 || startDay.Minute() != 0 || startDay.Second() != 0 {
		t.Errorf("StartOfDay esperado 00:00:00, obtenido %v", startDay)
	}

	endDay := calendar.EndOfDay(t0)
	if endDay.Hour() != 23 || endDay.Minute() != 59 || endDay.Nanosecond() != 999999999 {
		t.Errorf("EndOfDay esperado 23:59:59.999999999, obtenido %v", endDay)
	}

	startMonth := calendar.StartOfMonth(t0)
	if startMonth.Day() != 1 || startMonth.Hour() != 0 {
		t.Errorf("StartOfMonth esperado día 1 00:00, obtenido %v", startMonth)
	}

	endMonth := calendar.EndOfMonth(t0)
	if endMonth.Day() != 30 { // Septiembre tiene 30 días
		t.Errorf("EndOfMonth esperado día 30, obtenido %d", endMonth.Day())
	}

	// DaysInMonth
	if calendar.DaysInMonth(2028, time.February) != 29 {
		t.Errorf("Febrero 2028 (bisiesto) debe tener 29 días")
	}
	if calendar.DaysInMonth(2026, time.February) != 28 {
		t.Errorf("Febrero 2026 debe tener 28 días")
	}
	if calendar.DaysInMonth(2026, time.January) != 31 {
		t.Errorf("Enero 2026 debe tener 31 días")
	}
}

func TestCalendar_WeekAndYear(t *testing.T) {
	// Miércoles 15 de septiembre de 2026
	t0 := time.Date(2026, 9, 15, 14, 30, 0, 0, time.UTC)

	startWeek := calendar.StartOfWeek(t0)
	if startWeek.Weekday() != time.Monday {
		t.Errorf("StartOfWeek debe ser lunes, obtenido %v", startWeek.Weekday())
	}
	if startWeek.Day() != 14 { // 14 de septiembre es lunes
		t.Errorf("StartOfWeek debe ser 14, obtenido %d", startWeek.Day())
	}

	endWeek := calendar.EndOfWeek(t0)
	if endWeek.Weekday() != time.Sunday {
		t.Errorf("EndOfWeek debe ser domingo, obtenido %v", endWeek.Weekday())
	}

	startYear := calendar.StartOfYear(t0)
	if startYear.Month() != time.January || startYear.Day() != 1 {
		t.Errorf("StartOfYear esperado 1 de enero")
	}

	endYear := calendar.EndOfYear(t0)
	if endYear.Month() != time.December || endYear.Day() != 31 {
		t.Errorf("EndOfYear esperado 31 de diciembre")
	}
}
