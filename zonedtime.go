package timezoner

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ZonedTime es una estructura de nivel empresarial diseñada para manejar citas futuras,
// eventos de calendario, vuelos y recordatorios donde se debe preservar tanto el instante
// universal (UTC) para búsquedas e indexación en base de datos, como la zona horaria IANA de origen
// para garantizar que la hora local original nunca se altere ante cambios de leyes gubernamentales de DST.
type ZonedTime struct {
	UTC  DBTime `json:"utc"`  // Instante universal UTC (para consultas SQL, ordenamiento e indexación)
	Zone string `json:"zone"` // Identificador canónico IANA de origen (ej: "America/Lima")
}

// NewZonedTime crea un ZonedTime a partir de un time.Time y su zona IANA de origen.
func NewZonedTime(t time.Time, zoneName string) (ZonedTime, error) {
	loc, err := LoadLocation(zoneName)
	if err != nil {
		return ZonedTime{}, fmt.Errorf("%w: zona '%s'", ErrInvalidZone, zoneName)
	}

	canonicalZone := loc.String()
	tInLoc := t.In(loc)

	return ZonedTime{
		UTC:  NewDBTime(tInLoc),
		Zone: canonicalZone,
	}, nil
}

// ZonedFromLocal parsea una fecha/hora local (ej: "2026-09-01 10:00") y su zona IANA ("America/Lima"),
// calculando el UTC correspondiente y guardando la zona de origen.
func ZonedFromLocal(dateStr, defaultZone string) (ZonedTime, error) {
	utcTime, err := IngestFromString(dateStr, defaultZone)
	if err != nil {
		return ZonedTime{}, err
	}

	canonicalZone, err := NormalizeZone(defaultZone)
	if err != nil {
		canonicalZone = "UTC"
	}

	return ZonedTime{
		UTC:  NewDBTime(utcTime),
		Zone: canonicalZone,
	}, nil
}

// Local devuelve la fecha y hora recalculada en la zona de origen IANA según las reglas vigentes.
func (z ZonedTime) Local() (time.Time, error) {
	if z.Zone == "" {
		return z.UTC.Time, nil
	}
	loc, err := LoadLocation(z.Zone)
	if err != nil {
		return z.UTC.Time, err
	}
	return z.UTC.Time.In(loc), nil
}

// ForViewer proyecta este evento para un participante o usuario ubicado en cualquier otra parte del mundo.
func (z ZonedTime) ForViewer(viewerZone string, customLayout ...string) (UserTime, error) {
	return ProjectForUser(z.UTC.Time, viewerZone, customLayout...)
}

// Value implementa driver.Valuer serializando ZonedTime como un string compuesto o JSON para bases de datos SQL.
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

// Scan implementa sql.Scanner para reconstruir ZonedTime desde PostgreSQL (JSONB, TEXT) o MySQL.
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
		return fmt.Errorf("timezoner: tipo incompatible para escanear ZonedTime: %T", value)
	}
}

func (z *ZonedTime) unmarshalString(s string) error {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" || trimmed == "null" {
		*z = ZonedTime{}
		return nil
	}

	// Intento 1: Parsear como JSON {"utc": "...", "zone": "..."}
	if strings.HasPrefix(trimmed, "{") {
		type alias ZonedTime
		var aux alias
		if err := json.Unmarshal([]byte(trimmed), &aux); err == nil {
			*z = ZonedTime(aux)
			return nil
		}
	}

	// Intento 2: Parsear como formato compuesto "2026-09-01T15:00:00Z|America/Lima"
	parts := strings.Split(trimmed, "|")
	if len(parts) == 2 {
		utcParsed, err := IngestFromString(parts[0], "UTC")
		if err != nil {
			return err
		}
		z.UTC = NewDBTime(utcParsed)
		z.Zone = parts[1]
		return nil
	}

	// Intento 3: Parsear solo fecha simple
	utcParsed, err := IngestFromString(trimmed, "UTC")
	if err != nil {
		return err
	}
	z.UTC = NewDBTime(utcParsed)
	z.Zone = "UTC"
	return nil
}
