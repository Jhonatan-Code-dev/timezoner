package timezoner

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	// ErrInvalidInput se retorna cuando el tipo o valor de entrada no puede ser procesado.
	ErrInvalidInput = errors.New("timezoner: tipo de entrada no soportado o inválido")
	// ErrEmptyDateString se retorna cuando se intenta procesar un texto de fecha vacío.
	ErrEmptyDateString = errors.New("timezoner: la cadena de fecha no puede estar vacía")
)

// SupportedIngestLayouts contiene los formatos de fecha más comunes aceptados por APIs y frontends.
var SupportedIngestLayouts = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02 15:04",
	"2006-01-02",
	"02/01/2006 15:04:05",
	"02/01/2006 15:04",
	"02/01/2006",
}

// IngestNow devuelve el instante actual del sistema en UTC absoluto, limpio de reloj monotónico.
func IngestNow() time.Time {
	return time.Now().UTC().Round(0)
}

// IngestTime normaliza un time.Time a UTC absoluto, eliminando residuos monotónicos para almacenamiento en BD.
func IngestTime(t time.Time) time.Time {
	return t.UTC().Round(0)
}

// IngestFromLocal toma una fecha/hora expresada en una zona horaria local específica
// y la convierte a su instante UTC equivalente para guardarse en la base de datos.
func IngestFromLocal(localTime time.Time, sourceZone string) (time.Time, error) {
	loc, err := LoadLocation(sourceZone)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: zona de origen inválida '%s'", ErrInvalidZone, sourceZone)
	}

	tInLoc := time.Date(
		localTime.Year(),
		localTime.Month(),
		localTime.Day(),
		localTime.Hour(),
		localTime.Minute(),
		localTime.Second(),
		localTime.Nanosecond(),
		loc,
	)

	return tInLoc.UTC().Round(0), nil
}

// IngestFromString parsea una cadena de fecha intentando múltiples formatos estándar.
// Si no contiene información de zona, se asume defaultZone (o UTC si defaultZone es "").
func IngestFromString(dateStr, defaultZone string) (time.Time, error) {
	cleanStr := strings.TrimSpace(dateStr)
	if cleanStr == "" {
		return time.Time{}, ErrEmptyDateString
	}

	zone := strings.TrimSpace(defaultZone)
	if zone == "" {
		zone = "UTC"
	}

	loc, err := LoadLocation(zone)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: zona por defecto inválida '%s'", ErrInvalidZone, defaultZone)
	}

	for _, layout := range SupportedIngestLayouts {
		if parsed, err := time.ParseInLocation(layout, cleanStr, loc); err == nil {
			return parsed.UTC().Round(0), nil
		}
	}

	return time.Time{}, fmt.Errorf("%w: no se pudo parsear '%s' con los layouts soportados", ErrInvalidInput, dateStr)
}

// IngestFromUnix convierte una marca de tiempo Unix (segundos) a tiempo UTC para la BD.
func IngestFromUnix(seconds int64) time.Time {
	return time.Unix(seconds, 0).UTC().Round(0)
}

// IngestFromUnixMilli convierte una marca de tiempo Unix (milisegundos) a tiempo UTC para la BD.
func IngestFromUnixMilli(milli int64) time.Time {
	return time.UnixMilli(milli).UTC().Round(0)
}
