package project

import (
	"errors"
	"fmt"
	"time"

	"github.com/Jhonatan-Code-dev/timezoner/pkg/zone"
)

// ErrNoZonesProvided se retorna cuando no se especifican zonas en operaciones batch.
var ErrNoZonesProvided = errors.New("project: se requiere al menos una zona horaria")

// UserTime contiene la información de fecha y hora proyectada para la zona específica de un usuario.
type UserTime struct {
	LocalTime       time.Time     `json:"local_time"`
	DifferenceToUTC time.Duration `json:"difference_to_utc"`
	Zone            string        `json:"zone"`
	Abbreviation    string        `json:"abbreviation"`
	Formatted       string        `json:"formatted"`
	ISO8601         string        `json:"iso8601"`
	OffsetFormatted string        `json:"offset_formatted"`
	OffsetSeconds   int           `json:"offset_seconds"`
	IsDST           bool          `json:"is_dst"`
}

// DefaultDisplayLayout es el formato legible por defecto ("2006-01-02 15:04:05").
const DefaultDisplayLayout = "2006-01-02 15:04:05"

// ForUser proyecta un instante UTC proveniente de la base de datos a la zona horaria del usuario.
func ForUser(utcTime time.Time, userZone string, customLayout ...string) (UserTime, error) {
	loc, err := zone.LoadLocation(userZone)
	if err != nil {
		return UserTime{}, fmt.Errorf("%w: zona del usuario '%s'", zone.ErrInvalidZone, userZone)
	}

	layout := DefaultDisplayLayout
	if len(customLayout) > 0 && customLayout[0] != "" {
		layout = customLayout[0]
	}

	local := utcTime.In(loc)
	abbr, offset := local.Zone()

	isDST, _ := zone.IsDST(loc.String(), utcTime)

	return UserTime{
		LocalTime:       local,
		Zone:            loc.String(),
		Abbreviation:    abbr,
		OffsetSeconds:   offset,
		OffsetFormatted: zone.FormatOffset(offset),
		IsDST:           isDST,
		Formatted:       local.Format(layout),
		ISO8601:         local.Format(time.RFC3339),
		DifferenceToUTC: time.Duration(offset) * time.Second,
	}, nil
}

// Format convierte un instante UTC a la zona del usuario y devuelve únicamente la cadena formateada.
func Format(utcTime time.Time, userZone, layout string) (string, error) {
	uTime, err := ForUser(utcTime, userZone, layout)
	if err != nil {
		return "", err
	}
	return uTime.Formatted, nil
}

// BatchForUsers proyecta simultáneamente un instante UTC a múltiples zonas de usuarios distintos.
func BatchForUsers(utcTime time.Time, userZones []string) (map[string]UserTime, error) {
	if len(userZones) == 0 {
		return nil, ErrNoZonesProvided
	}

	results := make(map[string]UserTime, len(userZones))
	for _, z := range userZones {
		uTime, err := ForUser(utcTime, z)
		if err != nil {
			return nil, err
		}
		results[z] = uTime
	}

	return results, nil
}
