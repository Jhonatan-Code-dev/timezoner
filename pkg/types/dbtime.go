package types

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// DBTime es un tipo de persistencia temporal que garantiza UTC limpio sin reloj monotónico.
// Diseño deliberadamente restrictivo: NO embebe time.Time para evitar exponer métodos
// peligrosos como Local(), Add(), AddDate() que pueden corromper la invariante UTC de la BD.
type DBTime struct {
	t time.Time // campo privado: solo accesible a través de accesores controlados
}

// NewDBTime crea una instancia de DBTime normalizada a UTC absoluto sin reloj monotónico.
func NewDBTime(t time.Time) DBTime {
	return DBTime{t: t.UTC().Round(0)}
}

// NowDBTime devuelve el instante actual como DBTime en UTC.
func NowDBTime() DBTime {
	return DBTime{t: time.Now().UTC().Round(0)}
}

// Time retorna el time.Time subyacente en UTC. Es el único punto de acceso controlado.
func (d DBTime) Time() time.Time {
	return d.t
}

// UTC retorna el instante en UTC (equivalente a Time(), para compatibilidad con time.Time).
func (d DBTime) UTC() time.Time {
	return d.t
}

// IsZero informa si el valor es el instante cero.
func (d DBTime) IsZero() bool {
	return d.t.IsZero()
}

// Equal compara si dos DBTime representan el mismo instante.
func (d DBTime) Equal(other DBTime) bool {
	return d.t.Equal(other.t)
}

// EqualTime compara DBTime con un time.Time arbitrario.
func (d DBTime) EqualTime(t time.Time) bool {
	return d.t.Equal(t.UTC().Round(0))
}

// String retorna la representación ISO 8601 del instante.
func (d DBTime) String() string {
	if d.t.IsZero() {
		return ""
	}
	return d.t.Format(time.RFC3339Nano)
}

// Value implementa driver.Valuer para persistencia en SQL.
func (d DBTime) Value() (driver.Value, error) {
	if d.t.IsZero() {
		return nil, nil
	}
	return d.t, nil
}

// Scan implementa sql.Scanner para leer valores de la base de datos de forma segura.
func (d *DBTime) Scan(value any) error {
	if value == nil {
		d.t = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		d.t = v.UTC().Round(0)
		return nil
	case []byte:
		return d.parseString(string(v))
	case string:
		return d.parseString(v)
	case int64:
		d.t = time.Unix(v, 0).UTC().Round(0)
		return nil
	default:
		return fmt.Errorf("types: tipo no soportado para escanear DBTime: %T", value)
	}
}

func (d *DBTime) parseString(str string) error {
	s := strings.TrimSpace(str)
	if s == "" {
		d.t = time.Time{}
		return nil
	}

	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05.999999999-07:00",
		"2006-01-02 15:04:05.999999-07:00",
		"2006-01-02 15:04:05-07:00",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02 15:04:05.999999",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			d.t = t.UTC().Round(0)
			return nil
		}
	}

	return fmt.Errorf("types: formato inválido para DBTime: '%s'", str)
}

// MarshalJSON implementa json.Marshaler formateando en RFC3339Nano UTC (preserva nanosegundos).
func (d DBTime) MarshalJSON() ([]byte, error) {
	if d.t.IsZero() {
		return []byte("null"), nil
	}
	// RFC3339Nano conserva la precisión de nanosegundos que RFC3339 silenciosamente elimina.
	return []byte(fmt.Sprintf("%q", d.t.Format(time.RFC3339Nano))), nil
}

// UnmarshalJSON implementa json.Unmarshaler aceptando cadenas ISO 8601 / RFC 3339 o null.
func (d *DBTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" || str == `""` {
		d.t = time.Time{}
		return nil
	}

	str = strings.Trim(str, `"`)
	return d.parseString(str)
}
