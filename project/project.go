// Package project proporciona acceso directo a las utilidades de proyección de timezoner.
package project

import (
	"time"

	"timezoner"
)

// UserTime alias del tipo de timezoner.
type UserTime = timezoner.UserTime

// DefaultDisplayLayout formato legible por defecto.
const DefaultDisplayLayout = timezoner.DefaultDisplayLayout

func ForUser(utcTime time.Time, userZone string, customLayout ...string) (UserTime, error) {
	return timezoner.ProjectForUser(utcTime, userZone, customLayout...)
}

func Format(utcTime time.Time, userZone, layout string) (string, error) {
	return timezoner.ProjectFormat(utcTime, userZone, layout)
}

func BatchForUsers(utcTime time.Time, userZones []string) (map[string]UserTime, error) {
	return timezoner.ProjectBatchForUsers(utcTime, userZones)
}
