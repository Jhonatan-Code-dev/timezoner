// Package ingest proporciona utilidades para normalizar, sanitizar y convertir cualquier
// fecha u hora que ingrese al sistema a la hora global estándar (UTC) antes de ser guardada en la base de datos.
//
// Autor: Jhonatan.
package ingest

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"timezoner"
)

var (
	// ErrInvalidInput se retorna cuando el tipo o valor de entrada no puede ser procesado.
	ErrInvalidInput = errors.New("ingest: tipo de entrada no soportado o inválido")
	// ErrEmptyDateString se retorna cuando se intenta procesar un texto de fecha vacío.
	ErrEmptyDateString = errors.New("ingest: la cadena de fecha no puede estar vacía")
)

// SupportedLayouts contiene los formatos de fecha más comunes aceptados por APIs y frontends.
var SupportedLayouts = []string{
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

// Now devuelve el instante actual del sistema en UTC absoluto, limpio de reloj monotónico.
func Now() time.Time {
	return time.Now().UTC().Round(0)
}

// FromTime normaliza un time.Time a UTC absoluto, eliminando residuos monotónicos para almacenamiento en BD.
func FromTime(t time.Time) time.Time {
	return t.UTC().Round(0)
}

// FromLocal toma una fecha/hora expresada en una zona horaria local específica
// y la convierte a su instante UTC equivalente para guardarse en la base de datos.
func FromLocal(localTime time.Time, sourceZone string) (time.Time, error) {
	loc, err := timezoner.LoadLocation(sourceZone)
	if err != nil {
		return time.Time{}, fmt.Errorf("ingest: zona de origen inválida: %w", err)
	}

	// Reconstruir en la ubicación especificada
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

// FromString parsea una cadena de fecha intentando múltiples formatos estándar.
// Si no contiene información de zona, se asume defaultZone (o UTC si defaultZone es "").
func FromString(dateStr, defaultZone string) (time.Time, error) {
	cleanStr := strings.TrimSpace(dateStr)
	if cleanStr == "" {
		return time.Time{}, ErrEmptyDateString
	}

	// Determinar zona por defecto si la cadena no tiene offset
	zone := strings.TrimSpace(defaultZone)
	if zone == "" {
		zone = "UTC"
	}

	loc, err := timezoner.LoadLocation(zone)
	if err != nil {
		return time.Time{}, fmt.Errorf("ingest: zona por defecto inválida: %w", err)
	}

	for _, layout := range SupportedLayouts {
		if parsed, err := time.ParseInLocation(layout, cleanStr, loc); err == nil {
			return parsed.UTC().Round(0), nil
		}
	}

	return time.Time{}, fmt.Errorf("%w: no se pudo parsear '%s' con los layouts soportados", ErrInvalidInput, dateStr)
}

// FromUnix convierte una marca de tiempo Unix (segundos) a tiempo UTC para la BD.
func FromUnix(seconds int64) time.Time {
	return time.Unix(seconds, 0).UTC().Round(0)
}

// FromUnixMilli convierte una marca de tiempo Unix (milisegundos) a tiempo UTC para la BD.
func FromUnixMilli(milli int64) time.Time {
	return time.UnixMilli(milli).UTC().Round(0)
}
