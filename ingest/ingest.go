// Package ingest proporciona acceso directo a las utilidades de ingesta de timezoner.
package ingest

import (
	"time"

	"timezoner"
)

var (
	ErrInvalidInput    = timezoner.ErrInvalidInput
	ErrEmptyDateString = timezoner.ErrEmptyDateString
	SupportedLayouts   = timezoner.SupportedIngestLayouts
)

func Now() time.Time {
	return timezoner.IngestNow()
}

func FromTime(t time.Time) time.Time {
	return timezoner.IngestTime(t)
}

func FromLocal(localTime time.Time, sourceZone string) (time.Time, error) {
	return timezoner.IngestFromLocal(localTime, sourceZone)
}

func FromString(dateStr, defaultZone string) (time.Time, error) {
	return timezoner.IngestFromString(dateStr, defaultZone)
}

func FromUnix(seconds int64) time.Time {
	return timezoner.IngestFromUnix(seconds)
}

func FromUnixMilli(milli int64) time.Time {
	return timezoner.IngestFromUnixMilli(milli)
}
