// Package timezoner es la fachada pública principal de la librería. Proporciona una Fluent API
// idiomática y de alto rendimiento que integra de forma coherente los módulos de zonas IANA,
// persistencia en base de datos (DBTime, ZonedTime), días hábiles, tiempo relativo y cálculo de solapamientos.
//
// Diseñado y creado por Jhonatan bajo arquitectura limpia y modular.
package timezoner

import (
	"time"

	"github.com/Jhonatan-Code-dev/timezoner/pkg/calendar"
	"github.com/Jhonatan-Code-dev/timezoner/pkg/humanize"
	"github.com/Jhonatan-Code-dev/timezoner/pkg/ingest"
	"github.com/Jhonatan-Code-dev/timezoner/pkg/overlap"
	"github.com/Jhonatan-Code-dev/timezoner/pkg/project"
	"github.com/Jhonatan-Code-dev/timezoner/pkg/types"
	"github.com/Jhonatan-Code-dev/timezoner/pkg/zone"
)

// Re-exportación de tipos principales para conveniencia del consumidor de la librería.
type (
	DBTime         = types.DBTime
	ZonedTime      = types.ZonedTime
	ZoneDetail     = zone.Detail
	ZoneSnapshot   = zone.Snapshot
	UserTime       = project.UserTime
	OverlapRequest = overlap.Request
	OverlapSlot    = overlap.Slot
	WorkingHours   = overlap.WorkingHours
)

// Re-exportación de errores centinela.
var (
	ErrEmptyZoneName     = zone.ErrEmptyZoneName
	ErrInvalidZone       = zone.ErrInvalidZone
	ErrNoZonesProvided   = overlap.ErrNoZonesProvided
	ErrEmptyDateString   = ingest.ErrEmptyDateString
	ErrInvalidTimeFormat = ingest.ErrInvalidInput
)

// Re-exportación de constructores y funciones de módulos internos.
// Todas son funciones nombradas, no variables de tipo func (evita reemplazos externos peligrosos).
var (
	NewDBTime            = types.NewDBTime
	NowDBTime            = types.NowDBTime
	NewZonedTime         = types.NewZonedTime
	IngestNow            = ingest.Now
	IngestTime           = ingest.FromTime
	IngestFromLocal      = ingest.FromLocal
	IngestFromString     = ingest.FromString
	IngestFromUnix       = ingest.FromUnix
	ProjectForUser       = project.ForUser
	ProjectFormat        = project.Format
	ProjectBatchForUsers = project.BatchForUsers
	LoadLocation         = zone.LoadLocation
	IsValid              = zone.IsValid
	NormalizeZone        = zone.Normalize
	RegisterAlias        = zone.RegisterAlias
	CommonZones          = zone.CommonZones
	FormatOffset         = zone.FormatOffset
	GetZoneInfo          = zone.GetInfo
	IsDST                = zone.IsDST
	Difference           = zone.Difference
	IsWeekend            = calendar.IsWeekend
	IsWeekday            = calendar.IsWeekday
	AddBusinessDays      = calendar.AddBusinessDays
	StartOfDay           = calendar.StartOfDay
	EndOfDay             = calendar.EndOfDay
	StartOfMonth         = calendar.StartOfMonth
	EndOfMonth           = calendar.EndOfMonth
	StartOfWeek          = calendar.StartOfWeek
	EndOfWeek            = calendar.EndOfWeek
	StartOfYear          = calendar.StartOfYear
	EndOfYear            = calendar.EndOfYear
	DaysInMonth          = calendar.DaysInMonth
	Humanize             = humanize.Humanize
	HumanizeEn           = humanize.HumanizeEn
	SupportedLayouts     = ingest.SupportedLayouts
)

// ZonedFromLocal parsea una fecha local expresada como string y la convierte en ZonedTime.
// Es una función nombrada (no una variable func) para evitar reemplazos externos peligrosos.
func ZonedFromLocal(dateStr, zoneName string) (ZonedTime, error) {
	utcParsed, err := ingest.FromString(dateStr, zoneName)
	if err != nil {
		return ZonedTime{}, err
	}
	canonical, err := zone.Normalize(zoneName)
	if err != nil {
		canonical = "UTC"
	}
	return ZonedTime{UTC: types.NewDBTime(utcParsed), Zone: canonical}, nil
}

// Convert convierte un time.Time a la zona horaria destino.
func Convert(t time.Time, toZone string) (time.Time, error) {
	loc, err := zone.LoadLocation(toZone)
	if err != nil {
		return time.Time{}, err
	}
	return t.In(loc), nil
}

// ConvertBetween parsea un string en una zona de origen y lo convierte a la de destino.
func ConvertBetween(timeStr, layout, fromZone, toZone string) (time.Time, error) {
	fromLoc, err := zone.LoadLocation(fromZone)
	if err != nil {
		return time.Time{}, err
	}
	toLoc, err := zone.LoadLocation(toZone)
	if err != nil {
		return time.Time{}, err
	}
	parsed, err := time.ParseInLocation(layout, timeStr, fromLoc)
	if err != nil {
		return time.Time{}, err
	}
	return parsed.In(toLoc), nil
}

// NowIn devuelve la hora actual en la zona especificada.
func NowIn(zoneName string) (time.Time, error) {
	return Convert(time.Now(), zoneName)
}

// FormatIn formatea un tiempo en la zona especificada.
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

// Compare genera una instantánea de cómo se refleja un instante en varias zonas.
func Compare(at time.Time, zones ...string) ([]ZoneSnapshot, error) {
	if len(zones) == 0 {
		return nil, ErrNoZonesProvided
	}
	snapshots := make([]ZoneSnapshot, 0, len(zones))
	for _, z := range zones {
		loc, err := zone.LoadLocation(z)
		if err != nil {
			return nil, err
		}
		tInLoc := at.In(loc)
		abbr, offset := tInLoc.Zone()
		isDst, _ := zone.IsDST(z, at)
		snapshots = append(snapshots, ZoneSnapshot{
			Zone:            loc.String(),
			Abbreviation:    abbr,
			Time:            tInLoc,
			TimeUTC:         at.UTC().Format(time.RFC3339Nano),
			Formatted:       tInLoc.Format("2006-01-02 15:04:05 MST"),
			OffsetFormatted: zone.FormatOffset(offset),
			IsDST:           isDst,
		})
	}
	return snapshots, nil
}

// FindOverlap busca intervalos de solapamiento hábil para equipos distribuidos.
func FindOverlap(req OverlapRequest) ([]OverlapSlot, error) {
	return overlap.Find(req)
}

// --- Fluent API / Builder ---

// TimePoint encapsula un instante temporal para encadenar operaciones fluidas.
type TimePoint struct {
	t   time.Time
	err error
}

// At crea un nuevo TimePoint a partir de un time.Time.
func At(t time.Time) *TimePoint {
	return &TimePoint{t: t}
}

// Now crea un nuevo TimePoint con la hora actual del sistema.
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

// ToUTC convierte el TimePoint a la hora global UTC.
func (tp *TimePoint) ToUTC() *TimePoint {
	if tp.err != nil {
		return tp
	}
	tp.t = tp.t.UTC().Round(0)
	return tp
}

// AddBusinessDays añade n días laborables (lunes a viernes), preservando la hora local ante cambios DST.
func (tp *TimePoint) AddBusinessDays(days int) *TimePoint {
	if tp.err != nil {
		return tp
	}
	tp.t = calendar.AddBusinessDays(tp.t, days)
	return tp
}

// StartOfDay mueve el instante al inicio del día (00:00:00).
func (tp *TimePoint) StartOfDay() *TimePoint {
	if tp.err != nil {
		return tp
	}
	tp.t = calendar.StartOfDay(tp.t)
	return tp
}

// EndOfDay mueve el instante al final del día (23:59:59.999999999).
func (tp *TimePoint) EndOfDay() *TimePoint {
	if tp.err != nil {
		return tp
	}
	tp.t = calendar.EndOfDay(tp.t)
	return tp
}

// StartOfMonth mueve el instante al inicio del primer día del mes.
func (tp *TimePoint) StartOfMonth() *TimePoint {
	if tp.err != nil {
		return tp
	}
	tp.t = calendar.StartOfMonth(tp.t)
	return tp
}

// EndOfMonth mueve el instante al final del último día del mes.
func (tp *TimePoint) EndOfMonth() *TimePoint {
	if tp.err != nil {
		return tp
	}
	tp.t = calendar.EndOfMonth(tp.t)
	return tp
}

// StartOfWeek mueve el instante al inicio del lunes de la semana actual.
func (tp *TimePoint) StartOfWeek() *TimePoint {
	if tp.err != nil {
		return tp
	}
	tp.t = calendar.StartOfWeek(tp.t)
	return tp
}

// EndOfWeek mueve el instante al final del domingo de la semana actual.
func (tp *TimePoint) EndOfWeek() *TimePoint {
	if tp.err != nil {
		return tp
	}
	tp.t = calendar.EndOfWeek(tp.t)
	return tp
}

// Humanize devuelve una representación relativa humana en español ("hace 5 minutos", "en 2 días").
func (tp *TimePoint) Humanize(relativeTo ...time.Time) (string, error) {
	if tp.err != nil {
		return "", tp.err
	}
	return humanize.Humanize(tp.t, relativeTo...), nil
}

// HumanizeEn devuelve una representación relativa humana en inglés ("5 minutes ago", "in 2 days").
func (tp *TimePoint) HumanizeEn(relativeTo ...time.Time) (string, error) {
	if tp.err != nil {
		return "", tp.err
	}
	return humanize.HumanizeEn(tp.t, relativeTo...), nil
}

// AsDBTime convierte el TimePoint en una estructura DBTime para persistencia SQL y JSON.
func (tp *TimePoint) AsDBTime() (DBTime, error) {
	if tp.err != nil {
		return DBTime{}, tp.err
	}
	return types.NewDBTime(tp.t), nil
}

// MustDBTime devuelve el DBTime o produce pánico si hubo un error.
func (tp *TimePoint) MustDBTime() DBTime {
	if tp.err != nil {
		panic(tp.err)
	}
	return types.NewDBTime(tp.t)
}

// AsZonedTime convierte el TimePoint en un ZonedTime asociando la zona especificada.
func (tp *TimePoint) AsZonedTime(zoneName string) (ZonedTime, error) {
	if tp.err != nil {
		return ZonedTime{}, tp.err
	}
	return types.NewZonedTime(tp.t, zoneName)
}

// MustZonedTime devuelve ZonedTime o produce pánico si hubo un error.
func (tp *TimePoint) MustZonedTime(zoneName string) ZonedTime {
	if tp.err != nil {
		panic(tp.err)
	}
	z, err := types.NewZonedTime(tp.t, zoneName)
	if err != nil {
		panic(err)
	}
	return z
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

// Format formatea la hora en la zona actual del TimePoint.
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
	return zone.GetInfo(tp.t.Location().String(), tp.t)
}

// Err retorna el error acumulado en la cadena fluida, si lo hubiera.
func (tp *TimePoint) Err() error {
	return tp.err
}
