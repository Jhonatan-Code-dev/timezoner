package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Jhonatan-Code-dev/timezonermax/pkg/zone"
)

// ZonedTime preserva tanto el instante UTC como la zona IANA de origen para calendarios y eventos futuros.
// Patrón 2 de persistencia: almacena el instante universal Y la intención local original.
type ZonedTime struct {
	UTC  DBTime `json:"utc"`
	Zone string `json:"zone"`
}

// NewZonedTime crea una instancia de ZonedTime normalizada.
func NewZonedTime(t time.Time, zoneName string) (ZonedTime, error) {
	loc, err := zone.LoadLocation(zoneName)
	if err != nil {
		return ZonedTime{}, fmt.Errorf("%w: zona '%s'", zone.ErrInvalidZone, zoneName)
	}

	canonicalZone := loc.String()
	return ZonedTime{
		UTC:  NewDBTime(t),
		Zone: canonicalZone,
	}, nil
}

// Local devuelve la fecha recalculada en la zona de origen según las reglas vigentes.
func (z ZonedTime) Local() (time.Time, error) {
	if z.Zone == "" {
		return z.UTC.Time(), nil
	}
	loc, err := zone.LoadLocation(z.Zone)
	if err != nil {
		return z.UTC.Time(), err
	}
	return z.UTC.Time().In(loc), nil
}

// IsZero informa si el ZonedTime es el instante cero.
func (z ZonedTime) IsZero() bool {
	return z.UTC.IsZero()
}

// Value implementa driver.Valuer serializando ZonedTime como JSON string para SQL.
func (z ZonedTime) Value() (driver.Value, error) {
	if z.UTC.IsZero() {
		return nil, nil
	}
	data, err := json.Marshal(z)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

// Scan implementa sql.Scanner para reconstruir ZonedTime desde la base de datos.
func (z *ZonedTime) Scan(value any) error {
	if value == nil {
		*z = ZonedTime{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return z.unmarshalString(string(v))
	case string:
		return z.unmarshalString(v)
	case time.Time:
		z.UTC = NewDBTime(v)
		z.Zone = "UTC"
		return nil
	default:
		return fmt.Errorf("types: tipo incompatible para escanear ZonedTime: %T", value)
	}
}

func (z *ZonedTime) unmarshalString(s string) error {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" || trimmed == "null" {
		*z = ZonedTime{}
		return nil
	}

	// JSON object: {"utc": "...", "zone": "..."}
	if strings.HasPrefix(trimmed, "{") {
		// Usar struct auxiliar temporal con time.Time para evitar recursión en DBTime.
		var aux struct {
			UTC  string `json:"utc"`
			Zone string `json:"zone"`
		}
		if err := json.Unmarshal([]byte(trimmed), &aux); err != nil {
			return fmt.Errorf("types: JSON inválido para ZonedTime: %w", err)
		}
		var dbT DBTime
		if err := dbT.parseString(aux.UTC); err != nil {
			return fmt.Errorf("types: campo utc inválido en ZonedTime: %w", err)
		}
		// Validación explícita: si el campo utc está vacío es un error, no silencio.
		if dbT.IsZero() && aux.UTC != "" {
			return fmt.Errorf("types: campo utc resultó en zero value para ZonedTime: '%s'", aux.UTC)
		}
		z.UTC = dbT
		z.Zone = aux.Zone
		return nil
	}

	// Formato pipe: "2026-09-01T15:00:00Z|America/Lima"
	if parts := strings.Split(trimmed, "|"); len(parts) == 2 {
		if t, err := time.Parse(time.RFC3339, parts[0]); err == nil {
			z.UTC = NewDBTime(t)
			z.Zone = parts[1]
			return nil
		}
	}

	// Fallback: RFC3339 sin zona explícita
	if t, err := time.Parse(time.RFC3339, trimmed); err == nil {
		z.UTC = NewDBTime(t)
		z.Zone = "UTC"
		return nil
	}

	return fmt.Errorf("types: no se pudo parsear '%s' como ZonedTime", s)
}
