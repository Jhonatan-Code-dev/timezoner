// Package timezoner proporciona utilidades robustas e idiomáticas para la manipulación,
// conversión, comparación y planificación de horarios entre diferentes zonas horarias IANA.
//
// Diseñado para ser importado como librería independiente en cualquier proyecto Go.
package timezoner

import (
	"fmt"
	"time"
)

// ZoneDetail describe la información completa de una zona horaria en un instante de tiempo.
type ZoneDetail struct {
	ZoneName        string        `json:"zone_name"`
	Abbreviation    string        `json:"abbreviation"`
	OffsetSeconds   int           `json:"offset_seconds"`
	OffsetFormatted string        `json:"offset_formatted"`
	IsDST           bool          `json:"is_dst"`
	LocalTime       time.Time     `json:"local_time"`
	DifferenceToUTC time.Duration `json:"difference_to_utc"`
}

// ZoneSnapshot representa la hora correspondiente en una zona específica para una fecha base.
type ZoneSnapshot struct {
	Zone            string    `json:"zone"`
	Abbreviation    string    `json:"abbreviation"`
	Time            time.Time `json:"time"`
	Formatted       string    `json:"formatted"`
	OffsetFormatted string    `json:"offset_formatted"`
	IsDST           bool      `json:"is_dst"`
}

// WorkingHours define la jornada laboral (en formato 24h, ej: 9 a 18).
type WorkingHours struct {
	StartHour int `json:"start_hour"` // 0-23
	EndHour   int `json:"end_hour"`   // 0-23
}

// DefaultWorkingHours devuelve el horario laboral habitual (09:00 a 17:00).
func DefaultWorkingHours() WorkingHours {
	return WorkingHours{StartHour: 9, EndHour: 17}
}

// OverlapRequest define los parámetros para encontrar solapamiento entre zonas.
type OverlapRequest struct {
	Date         time.Time               // Fecha base a analizar (el año, mes y día)
	Zones        []string                // Lista de zonas a comparar
	CustomHours  map[string]WorkingHours // Horarios laborales personalizados por zona (opcional)
	DefaultHours WorkingHours            // Horario laboral por defecto si no está en CustomHours
	SlotDuration time.Duration           // Duración mínima de cada slot (ej: 1 hora o 30 minutos)
}

// OverlapSlot representa una ventana horaria donde todos los participantes están en horario laboral.
type OverlapSlot struct {
	StartTimeUTC time.Time            `json:"start_time_utc"`
	EndTimeUTC   time.Time            `json:"end_time_utc"`
	Duration     time.Duration        `json:"duration"`
	ZoneTimes    map[string]time.Time `json:"zone_times"`
}

// Convert convierte un time.Time a la zona horaria destino especificada por su nombre IANA o alias.
func Convert(t time.Time, toZone string) (time.Time, error) {
	loc, err := LoadLocation(toZone)
	if err != nil {
		return time.Time{}, err
	}
	return t.In(loc), nil
}

// ConvertBetween parsea un string de fecha/hora en una zona de origen y lo convierte a la zona de destino.
func ConvertBetween(timeStr, layout, fromZone, toZone string) (time.Time, error) {
	fromLoc, err := LoadLocation(fromZone)
	if err != nil {
		return time.Time{}, fmt.Errorf("fromZone error: %w", err)
	}

	toLoc, err := LoadLocation(toZone)
	if err != nil {
		return time.Time{}, fmt.Errorf("toZone error: %w", err)
	}

	parsed, err := time.ParseInLocation(layout, timeStr, fromLoc)
	if err != nil {
		return time.Time{}, fmt.Errorf("timezoner: error al parsear tiempo '%s' con layout '%s': %w", timeStr, layout, err)
	}

	return parsed.In(toLoc), nil
}

// NowIn devuelve la hora actual en la zona horaria especificada.
func NowIn(zoneName string) (time.Time, error) {
	return Convert(time.Now(), zoneName)
}

// FormatIn convierte la hora dada a la zona y la formatea según el layout proporcionado.
// Si layout es "", se utiliza time.RFC3339 por defecto.
func FormatIn(t time.Time, zoneName, layout string) (string, error) {
	targetTime, err := Convert(t, zoneName)
	if err != nil {
		return "", err
	}
	if layout == "" {
		layout = time.RFC3339
	}
	return targetTime.Format(layout), nil
}

// FormatOffset devuelve el offset en formato legible como "+02:00" o "-05:00".
func FormatOffset(offsetSeconds int) string {
	sign := "+"
	if offsetSeconds < 0 {
		sign = "-"
		offsetSeconds = -offsetSeconds
	}
	hours := offsetSeconds / 3600
	minutes := (offsetSeconds % 3600) / 60
	return fmt.Sprintf("%s%02d:%02d", sign, hours, minutes)
}

// GetZoneInfo devuelve información detallada de una zona horaria en un instante específico.
func GetZoneInfo(zoneName string, at time.Time) (ZoneDetail, error) {
	loc, err := LoadLocation(zoneName)
	if err != nil {
		return ZoneDetail{}, err
	}

	localTime := at.In(loc)
	abbr, offset := localTime.Zone()

	// Detección de DST comparando con el offset estándar (calculado en solsticio/invierno)
	// Comparamos el offset actual con enero y julio para verificar si hay DST activo
	jan := time.Date(at.Year(), time.January, 1, 12, 0, 0, 0, loc)
	jul := time.Date(at.Year(), time.July, 1, 12, 0, 0, 0, loc)
	_, janOff := jan.Zone()
	_, julOff := jul.Zone()

	isDST := false
	stdOffset := janOff
	if julOff < janOff {
		stdOffset = julOff
	}
	if offset > stdOffset {
		isDST = true
	}

	return ZoneDetail{
		ZoneName:        loc.String(),
		Abbreviation:    abbr,
		OffsetSeconds:   offset,
		OffsetFormatted: FormatOffset(offset),
		IsDST:           isDST,
		LocalTime:       localTime,
		DifferenceToUTC: time.Duration(offset) * time.Second,
	}, nil
}

// IsDST indica si una zona horaria se encuentra en horario de verano (DST) en un instante dado.
func IsDST(zoneName string, at time.Time) (bool, error) {
	info, err := GetZoneInfo(zoneName, at)
	if err != nil {
		return false, err
	}
	return info.IsDST, nil
}

// Difference calcula la diferencia horaria (duración) entre zoneA y zoneB en un instante determinado.
// Si zoneA está en UTC+2 y zoneB está en UTC-5, la diferencia (zoneA - zoneB) será de +7 horas.
func Difference(zoneA, zoneB string, at time.Time) (time.Duration, error) {
	locA, err := LoadLocation(zoneA)
	if err != nil {
		return 0, fmt.Errorf("zoneA inválida: %w", err)
	}
	locB, err := LoadLocation(zoneB)
	if err != nil {
		return 0, fmt.Errorf("zoneB inválida: %w", err)
	}

	_, offsetA := at.In(locA).Zone()
	_, offsetB := at.In(locB).Zone()

	diffSec := offsetA - offsetB
	return time.Duration(diffSec) * time.Second, nil
}

// Compare genera una instantánea de cómo se refleja un instante 'at' en una lista de zonas horarias.
func Compare(at time.Time, zones ...string) ([]ZoneSnapshot, error) {
	var snapshots []ZoneSnapshot

	for _, z := range zones {
		loc, err := LoadLocation(z)
		if err != nil {
			return nil, err
		}

		tInLoc := at.In(loc)
		abbr, offset := tInLoc.Zone()

		isDst, _ := IsDST(z, at)

		snapshots = append(snapshots, ZoneSnapshot{
			Zone:            loc.String(),
			Abbreviation:    abbr,
			Time:            tInLoc,
			Formatted:       tInLoc.Format("2006-01-02 15:04:05 MST"),
			OffsetFormatted: FormatOffset(offset),
			IsDST:           isDst,
		})
	}

	return snapshots, nil
}

// FindOverlap busca intervalos de tiempo durante un día en los que todas las zonas
// indicadas se encuentran dentro de sus respectivas jornadas laborales.
func FindOverlap(req OverlapRequest) ([]OverlapSlot, error) {
	if len(req.Zones) == 0 {
		return nil, fmt.Errorf("timezoner: se requiere al menos una zona horaria")
	}

	slotDur := req.SlotDuration
	if slotDur <= 0 {
		slotDur = time.Hour
	}

	defHours := req.DefaultHours
	if defHours.StartHour == 0 && defHours.EndHour == 0 {
		defHours = DefaultWorkingHours()
	}

	// Cargar ubicaciones y resolver horas de cada zona
	type zoneRule struct {
		name string
		loc  *time.Location
		work WorkingHours
	}

	rules := make([]zoneRule, 0, len(req.Zones))
	for _, z := range req.Zones {
		loc, err := LoadLocation(z)
		if err != nil {
			return nil, err
		}
		wh := defHours
		if custom, ok := req.CustomHours[z]; ok {
			wh = custom
		}
		rules = append(rules, zoneRule{
			name: z,
			loc:  loc,
			work: wh,
		})
	}

	// Anclamos la ventana de 24 horas a la fecha objetivo de la primera zona
	year, month, day := req.Date.Date()
	baseLoc := rules[0].loc
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, baseLoc)
	endOfDay := startOfDay.Add(24 * time.Hour)

	var validSlots []OverlapSlot
	step := 30 * time.Minute

	for current := startOfDay; current.Before(endOfDay); current = current.Add(step) {
		allMatch := true
		zoneTimes := make(map[string]time.Time)

		for _, r := range rules {
			local := current.In(r.loc)
			hour := local.Hour()
			if hour < r.work.StartHour || hour >= r.work.EndHour {
				allMatch = false
				break
			}
			zoneTimes[r.name] = local
		}

		if allMatch {
			validSlots = append(validSlots, OverlapSlot{
				StartTimeUTC: current.UTC(),
				EndTimeUTC:   current.Add(step).UTC(),
				Duration:     step,
				ZoneTimes:    zoneTimes,
			})
		}
	}

	return mergeConsecutiveSlots(validSlots, slotDur), nil
}

func mergeConsecutiveSlots(slots []OverlapSlot, minDuration time.Duration) []OverlapSlot {
	if len(slots) == 0 {
		return nil
	}

	var merged []OverlapSlot
	current := slots[0]

	for i := 1; i < len(slots); i++ {
		next := slots[i]
		if current.EndTimeUTC.Equal(next.StartTimeUTC) {
			current.EndTimeUTC = next.EndTimeUTC
			current.Duration = current.EndTimeUTC.Sub(current.StartTimeUTC)
		} else {
			if current.Duration >= minDuration {
				merged = append(merged, current)
			}
			current = next
		}
	}

	if current.Duration >= minDuration {
		merged = append(merged, current)
	}

	return merged
}

// --- Fluent API / Builder ---

// TimePoint encapsula un instante y permite encadenar operaciones fluidas.
type TimePoint struct {
	t   time.Time
	err error
}

// At crea un nuevo TimePoint a partir de un time.Time.
func At(t time.Time) *TimePoint {
	return &TimePoint{t: t}
}

// Now crea un nuevo TimePoint con la hora actual.
func Now() *TimePoint {
	return &TimePoint{t: time.Now()}
}

// In convierte el instante a la zona especificada.
func (tp *TimePoint) In(zoneName string) *TimePoint {
	if tp.err != nil {
		return tp
	}
	converted, err := Convert(tp.t, zoneName)
	if err != nil {
		tp.err = err
		return tp
	}
	tp.t = converted
	return tp
}

// Time devuelve el time.Time subyacente y el error acumulado si lo hubiera.
func (tp *TimePoint) Time() (time.Time, error) {
	return tp.t, tp.err
}

// MustTime devuelve el time.Time o produce un pánico si ocurrió un error.
func (tp *TimePoint) MustTime() time.Time {
	if tp.err != nil {
		panic(tp.err)
	}
	return tp.t
}

// Format formatea la hora en la zona actual.
func (tp *TimePoint) Format(layout string) (string, error) {
	if tp.err != nil {
		return "", tp.err
	}
	if layout == "" {
		layout = time.RFC3339
	}
	return tp.t.Format(layout), nil
}

// Info devuelve el detalle de la zona horaria actual del TimePoint.
func (tp *TimePoint) Info() (ZoneDetail, error) {
	if tp.err != nil {
		return ZoneDetail{}, tp.err
	}
	return GetZoneInfo(tp.t.Location().String(), tp.t)
}
