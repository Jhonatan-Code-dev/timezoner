package zone

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
	_ "time/tzdata" // Embebe la base de datos IANA completa para portabilidad universal
)

var (
	// ErrEmptyZoneName se retorna cuando el nombre de la zona está vacío.
	ErrEmptyZoneName = errors.New("zone: el nombre de la zona horaria no puede estar vacío")
	// ErrInvalidZone se retorna cuando la zona IANA o alias no existe.
	ErrInvalidZone = errors.New("zone: zona horaria no válida")

	locationCache sync.Map

	// zoneAliases mapea abreviaciones comunes y sinónimos a identificadores IANA válidos.
	zoneAliases = map[string]string{
		"UTC":    "UTC",
		"GMT":    "Etc/GMT",
		"EST":    "America/New_York",
		"EDT":    "America/New_York",
		"CST":    "America/Chicago",
		"CDT":    "America/Chicago",
		"MST":    "America/Denver",
		"MDT":    "America/Denver",
		"PST":    "America/Los_Angeles",
		"PDT":    "America/Los_Angeles",
		"PET":    "America/Lima",
		"COT":    "America/Bogota",
		"CLT":    "America/Santiago",
		"ART":    "America/Argentina/Buenos_Aires",
		"BRT":    "America/Sao_Paulo",
		"CET":    "Europe/Paris",
		"CEST":   "Europe/Paris",
		"WET":    "Europe/Lisbon",
		"EET":    "Europe/Athens",
		"EEST":   "Europe/Athens",
		"BST":    "Europe/London",
		"JST":    "Asia/Tokyo",
		"KST":    "Asia/Seoul",
		"CST-CN": "Asia/Shanghai",
		"IST":    "Asia/Kolkata",
		"SGT":    "Asia/Singapore",
		"AEST":   "Australia/Sydney",
		"AEDT":   "Australia/Sydney",
		"NZST":   "Pacific/Auckland",
		"NZDT":   "Pacific/Auckland",
	}

	// CommonZonesList lista representativa de husos horarios IANA por región.
	CommonZonesList = []string{
		"UTC",
		"Africa/Cairo",
		"Africa/Johannesburg",
		"Africa/Lagos",
		"America/Argentina/Buenos_Aires",
		"America/Bogota",
		"America/Caracas",
		"America/Chicago",
		"America/Denver",
		"America/Lima",
		"America/Los_Angeles",
		"America/Mexico_City",
		"America/New_York",
		"America/Santiago",
		"America/Sao_Paulo",
		"America/Toronto",
		"Asia/Bangkok",
		"Asia/Dubai",
		"Asia/Hong_Kong",
		"Asia/Jakarta",
		"Asia/Kolkata",
		"Asia/Seoul",
		"Asia/Shanghai",
		"Asia/Singapore",
		"Asia/Tokyo",
		"Australia/Melbourne",
		"Australia/Sydney",
		"Europe/Amsterdam",
		"Europe/Berlin",
		"Europe/London",
		"Europe/Madrid",
		"Europe/Paris",
		"Europe/Rome",
		"Pacific/Auckland",
		"Pacific/Honolulu",
	}
)

// Detail describe la información completa de una zona horaria en un instante de tiempo.
type Detail struct {
	LocalTime       time.Time     `json:"local_time"`
	DifferenceToUTC time.Duration `json:"difference_to_utc"`
	ZoneName        string        `json:"zone_name"`
	Abbreviation    string        `json:"abbreviation"`
	OffsetFormatted string        `json:"offset_formatted"`
	OffsetSeconds   int           `json:"offset_seconds"`
	IsDST           bool          `json:"is_dst"`
}

// Snapshot representa la hora correspondiente en una zona específica para una fecha base.
type Snapshot struct {
	Time            time.Time `json:"time"`
	Zone            string    `json:"zone"`
	Abbreviation    string    `json:"abbreviation"`
	Formatted       string    `json:"formatted"`
	OffsetFormatted string    `json:"offset_formatted"`
	IsDST           bool      `json:"is_dst"`
}

// LoadLocation carga un *time.Location con soporte de caché en memoria y resolución de alias.
func LoadLocation(zoneName string) (*time.Location, error) {
	name := strings.TrimSpace(zoneName)
	if name == "" {
		return nil, ErrEmptyZoneName
	}

	if canonical, ok := zoneAliases[strings.ToUpper(name)]; ok {
		name = canonical
	}

	if loc, ok := locationCache.Load(name); ok {
		return loc.(*time.Location), nil
	}

	loc, err := time.LoadLocation(name)
	if err != nil {
		return nil, fmt.Errorf("%w: '%s'", ErrInvalidZone, zoneName)
	}

	locationCache.Store(name, loc)
	return loc, nil
}

// IsValid verifica si un nombre de zona o alias es válido.
func IsValid(zoneName string) bool {
	_, err := LoadLocation(zoneName)
	return err == nil
}

// Normalize convierte un alias o nombre al identificador canónico IANA correspondiente.
func Normalize(zoneName string) (string, error) {
	loc, err := LoadLocation(zoneName)
	if err != nil {
		return "", err
	}
	return loc.String(), nil
}

// RegisterAlias permite registrar alias personalizados en tiempo de ejecución.
func RegisterAlias(alias, ianaZone string) error {
	if !IsValid(ianaZone) {
		return fmt.Errorf("%w: zona destino '%s'", ErrInvalidZone, ianaZone)
	}
	zoneAliases[strings.ToUpper(strings.TrimSpace(alias))] = ianaZone
	return nil
}

// CommonZones retorna una copia de la lista de zonas comunes.
func CommonZones() []string {
	result := make([]string, len(CommonZonesList))
	copy(result, CommonZonesList)
	return result
}

// FormatOffset convierte un offset en segundos a string legible como "+02:00" o "-05:00".
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

// GetInfo devuelve el detalle de una zona horaria en un instante determinado.
func GetInfo(zoneName string, at time.Time) (Detail, error) {
	loc, err := LoadLocation(zoneName)
	if err != nil {
		return Detail{}, err
	}

	localTime := at.In(loc)
	abbr, offset := localTime.Zone()

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

	return Detail{
		LocalTime:       localTime,
		DifferenceToUTC: time.Duration(offset) * time.Second,
		ZoneName:        loc.String(),
		Abbreviation:    abbr,
		OffsetFormatted: FormatOffset(offset),
		OffsetSeconds:   offset,
		IsDST:           isDST,
	}, nil
}

// IsDST reporta si una zona se encuentra en horario de verano en el instante dado.
func IsDST(zoneName string, at time.Time) (bool, error) {
	info, err := GetInfo(zoneName, at)
	if err != nil {
		return false, err
	}
	return info.IsDST, nil
}

// Difference calcula la duración diferencial entre zoneA y zoneB en el instante dado.
func Difference(zoneA, zoneB string, at time.Time) (time.Duration, error) {
	locA, err := LoadLocation(zoneA)
	if err != nil {
		return 0, fmt.Errorf("%w: zoneA '%s'", ErrInvalidZone, zoneA)
	}
	locB, err := LoadLocation(zoneB)
	if err != nil {
		return 0, fmt.Errorf("%w: zoneB '%s'", ErrInvalidZone, zoneB)
	}

	_, offsetA := at.In(locA).Zone()
	_, offsetB := at.In(locB).Zone()

	diffSec := offsetA - offsetB
	return time.Duration(diffSec) * time.Second, nil
}
