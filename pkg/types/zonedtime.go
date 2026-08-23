package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"timezoner/pkg/zone"
)

// ZonedTime preserva tanto el instante UTC como la zona IANA de origen para calendarios y eventos futuros.
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
	tInLoc := t.In(loc)

	return ZonedTime{
		UTC:  NewDBTime(tInLoc),
		Zone: canonicalZone,
	}, nil
}

// Local devuelve la fecha recalculada en la zona de origen según las reglas vigentes.
func (z ZonedTime) Local() (time.Time, error) {
	if z.Zone == "" {
		return z.UTC.Time, nil
	}
	loc, err := zone.LoadLocation(z.Zone)
	if err != nil {
		return z.UTC.Time, err
	}
	return z.UTC.Time.In(loc), nil
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

	if strings.HasPrefix(trimmed, "{") {
		type alias ZonedTime
		var aux alias
		if err := json.Unmarshal([]byte(trimmed), &aux); err == nil {
			*z = ZonedTime(aux)
			return nil
		}
	}

	parts := strings.Split(trimmed, "|")
	if len(parts) == 2 {
		if t, err := time.Parse(time.RFC3339, parts[0]); err == nil {
			z.UTC = NewDBTime(t)
			z.Zone = parts[1]
			return nil
		}
	}

	if t, err := time.Parse(time.RFC3339, trimmed); err == nil {
		z.UTC = NewDBTime(t)
		z.Zone = "UTC"
		return nil
	}

	return fmt.Errorf("types: no se pudo parsear '%s' como ZonedTime", s)
}
