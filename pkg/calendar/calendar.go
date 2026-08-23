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
// Preserva la hora local original del instante incluso cuando la aritmética atraviesa
// una transición de horario de verano (DST), garantizando correctitud temporal.
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

	// Capturar componentes de hora en la zona del instante original.
	// Esto es crítico: cuando AddDate atraviesa un límite DST, la hora del día puede desplazarse.
	// Recomputamos la hora original en la zona destino para preservar la intención del usuario.
	hour, min, sec := t.Clock()
	ns := t.Nanosecond()
	loc := t.Location()

	current := t
	for remaining > 0 {
		current = current.AddDate(0, 0, step)
		if IsWeekday(current) {
			remaining--
		}
	}

	// Reconstruir con la hora local original en la zona destino para neutralizar el efecto DST.
	year, month, day := current.Date()
	return time.Date(year, month, day, hour, min, sec, ns, loc)
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

// StartOfWeek devuelve el lunes de la semana del instante dado en la zona de la fecha.
func StartOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7 // domingo = 7 en orden ISO (lunes = 1)
	}
	year, month, day := t.Date()
	return time.Date(year, month, day-weekday+1, 0, 0, 0, 0, t.Location())
}

// EndOfWeek devuelve el domingo al final de la semana del instante dado.
func EndOfWeek(t time.Time) time.Time {
	start := StartOfWeek(t)
	year, month, day := start.Date()
	return time.Date(year, month, day+6, 23, 59, 59, 999999999, t.Location())
}

// StartOfYear devuelve el primer instante del año del instante dado.
func StartOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, t.Location())
}

// EndOfYear devuelve el último nanosegundo del último día del año del instante dado.
func EndOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), time.December, 31, 23, 59, 59, 999999999, t.Location())
}
