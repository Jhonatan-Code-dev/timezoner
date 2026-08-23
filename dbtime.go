package timezoner

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// DBTime es un tipo temporal envoltorio diseñado para interoperar limpiamente con bases de datos (SQL)
// y APIs JSON. Garantiza que cualquier fecha se almacene y serialice en UTC estricto sin reloj monotónico.
type DBTime struct {
	time.Time
}

// NewDBTime crea una instancia de DBTime normalizada a UTC absoluto.
func NewDBTime(t time.Time) DBTime {
	return DBTime{Time: t.UTC().Round(0)}
}

// NowDBTime devuelve el instante actual como DBTime en UTC.
func NowDBTime() DBTime {
	return DBTime{Time: time.Now().UTC().Round(0)}
}

// Value implementa la interfaz driver.Valuer para persistencia en SQL (PostgreSQL, MySQL, SQLite).
func (d DBTime) Value() (driver.Value, error) {
	if d.IsZero() {
		return nil, nil
	}
	return d.UTC().Round(0), nil
}

// Scan implementa la interfaz sql.Scanner para leer valores de la base de datos de forma segura.
func (d *DBTime) Scan(value any) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		d.Time = v.UTC().Round(0)
		return nil
	case []byte:
		return d.parseString(string(v))
	case string:
		return d.parseString(v)
	case int64:
		d.Time = time.Unix(v, 0).UTC().Round(0)
		return nil
	default:
		return fmt.Errorf("timezoner: tipo no soportado para escanear DBTime: %T", value)
	}
}

func (d *DBTime) parseString(str string) error {
	s := strings.TrimSpace(str)
	if s == "" {
		d.Time = time.Time{}
		return nil
	}

	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05.999999-07:00",
		"2006-01-02 15:04:05-07:00",
		"2006-01-02 15:04:05.999999",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			d.Time = t.UTC().Round(0)
			return nil
		}
	}

	return fmt.Errorf("%w: '%s' no pudo ser parseado como DBTime", ErrInvalidTimeFormat, str)
}

// MarshalJSON implementa json.Marshaler formateando siempre en RFC3339 UTC.
func (d DBTime) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%q", d.UTC().Round(0).Format(time.RFC3339))), nil
}

// UnmarshalJSON implementa json.Unmarshaler aceptando cadenas ISO 8601 / RFC 3339 o null.
func (d *DBTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" || str == `""` {
		d.Time = time.Time{}
		return nil
	}

	str = strings.Trim(str, `"`)
	return d.parseString(str)
}
