package timezoner

import (
	"fmt"
	"time"
)

// UserTime contiene la información de fecha y hora proyectada para la zona específica de un usuario.
type UserTime struct {
	LocalTime       time.Time     `json:"local_time"`        // 24 bytes
	DifferenceToUTC time.Duration `json:"difference_to_utc"` // 8 bytes
	Zone            string        `json:"zone"`              // 16 bytes
	Abbreviation    string        `json:"abbreviation"`      // 16 bytes
	Formatted       string        `json:"formatted"`         // 16 bytes
	ISO8601         string        `json:"iso8601"`           // 16 bytes
	OffsetFormatted string        `json:"offset_formatted"`  // 16 bytes
	OffsetSeconds   int           `json:"offset_seconds"`    // 8 bytes
	IsDST           bool          `json:"is_dst"`            // 1 byte
}

// DefaultDisplayLayout es el formato legible por defecto (ej: "2006-01-02 15:04:05").
const DefaultDisplayLayout = "2006-01-02 15:04:05"

// ProjectForUser proyecta un instante UTC proveniente de la base de datos a la zona horaria del usuario.
// Si customLayout es provisto, se utiliza para el campo Formatted; de lo contrario se usa DefaultDisplayLayout.
func ProjectForUser(utcTime time.Time, userZone string, customLayout ...string) (UserTime, error) {
	loc, err := LoadLocation(userZone)
	if err != nil {
		return UserTime{}, fmt.Errorf("%w: zona del usuario '%s'", ErrInvalidZone, userZone)
	}

	layout := DefaultDisplayLayout
	if len(customLayout) > 0 && customLayout[0] != "" {
		layout = customLayout[0]
	}

	local := utcTime.In(loc)
	abbr, offset := local.Zone()

	isDST, _ := IsDST(loc.String(), utcTime)

	return UserTime{
		LocalTime:       local,
		Zone:            loc.String(),
		Abbreviation:    abbr,
		OffsetSeconds:   offset,
		OffsetFormatted: FormatOffset(offset),
		IsDST:           isDST,
		Formatted:       local.Format(layout),
		ISO8601:         local.Format(time.RFC3339),
		DifferenceToUTC: time.Duration(offset) * time.Second,
	}, nil
}

// ProjectFormat convierte un instante UTC a la zona del usuario y devuelve únicamente la cadena formateada.
func ProjectFormat(utcTime time.Time, userZone, layout string) (string, error) {
	uTime, err := ProjectForUser(utcTime, userZone, layout)
	if err != nil {
		return "", err
	}
	return uTime.Formatted, nil
}

// ProjectBatchForUsers proyecta simultáneamente un instante UTC a múltiples zonas de usuarios distintos.
func ProjectBatchForUsers(utcTime time.Time, userZones []string) (map[string]UserTime, error) {
	if len(userZones) == 0 {
		return nil, ErrNoZonesProvided
	}

	results := make(map[string]UserTime, len(userZones))
	for _, z := range userZones {
		uTime, err := ProjectForUser(utcTime, z)
		if err != nil {
			return nil, err
		}
		results[z] = uTime
	}

	return results, nil
}
