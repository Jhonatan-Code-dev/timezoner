package calendar

import "time"

// IsWeekend indica si la fecha corresponde a sábado o domingo.
func IsWeekend(t time.Time) bool {
	w := t.Weekday()
	return w == time.Saturday || w == time.Sunday
}

// IsWeekday indica si la fecha corresponde a un día laborable de lunes a viernes.
func IsWeekday(t time.Time) bool {
	return !IsWeekend(t)
}

// AddBusinessDays añade o resta días hábiles de lunes a viernes ignorando fines de semana.
func AddBusinessDays(t time.Time, days int) time.Time {
	if days == 0 {
		return t
	}

	step := 1
	remaining := days
	if days < 0 {
		step = -1
		remaining = -days
	}

	current := t
	for remaining > 0 {
		current = current.AddDate(0, 0, step)
		if IsWeekday(current) {
			remaining--
		}
	}

	return current
}

// StartOfDay devuelve el inicio del día (00:00:00.000000000) en la zona de la fecha.
func StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// EndOfDay devuelve el final del día (23:59:59.999999999) en la zona de la fecha.
func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
}

// StartOfMonth devuelve el primer instante del primer día del mes (00:00:00).
func StartOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth devuelve el último nanosegundo del último día del mes correspondiente.
func EndOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	firstOfNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, t.Location())
	return firstOfNextMonth.Add(-time.Nanosecond)
}

// DaysInMonth calcula la cantidad exacta de días que contiene un mes considerando bisiestos.
func DaysInMonth(year int, month time.Month) int {
	firstOfNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	lastDayOfMonth := firstOfNextMonth.AddDate(0, 0, -1)
	return lastDayOfMonth.Day()
}
