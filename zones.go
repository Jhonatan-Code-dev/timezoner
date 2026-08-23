package timezoner

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

var (
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

	// CommonZonesList contiene una lista representativa de husos horarios IANA por región.
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

// LoadLocation carga un *time.Location con soporte de caché en memoria y resolución de alias.
func LoadLocation(zoneName string) (*time.Location, error) {
	name := strings.TrimSpace(zoneName)
	if name == "" {
		return nil, fmt.Errorf("timezoner: el nombre de la zona horaria no puede estar vacío")
	}

	// Revisar alias
	if canonical, ok := zoneAliases[strings.ToUpper(name)]; ok {
		name = canonical
	}

	// Buscar en caché
	if loc, ok := locationCache.Load(name); ok {
		return loc.(*time.Location), nil
	}

	loc, err := time.LoadLocation(name)
	if err != nil {
		return nil, fmt.Errorf("timezoner: zona horaria no válida '%s': %w", zoneName, err)
	}

	locationCache.Store(name, loc)
	return loc, nil
}

// IsValid verifica si un nombre de zona o alias es válido y puede ser cargado por Go.
func IsValid(zoneName string) bool {
	_, err := LoadLocation(zoneName)
	return err == nil
}

// NormalizeZone convierte un alias o nombre al identificador canónico correspondiente si existe.
func NormalizeZone(zoneName string) (string, error) {
	loc, err := LoadLocation(zoneName)
	if err != nil {
		return "", err
	}
	return loc.String(), nil
}

// RegisterAlias permite a los proyectos consumidores registrar alias personalizados.
func RegisterAlias(alias, ianaZone string) error {
	if !IsValid(ianaZone) {
		return fmt.Errorf("timezoner: la zona IANA destino '%s' no es válida", ianaZone)
	}
	zoneAliases[strings.ToUpper(strings.TrimSpace(alias))] = ianaZone
	return nil
}

// CommonZones retorna una copia de la lista de zonas horarias comunes.
func CommonZones() []string {
	result := make([]string, len(CommonZonesList))
	copy(result, CommonZonesList)
	return result
}
