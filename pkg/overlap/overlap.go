package overlap

import (
	"errors"
	"time"

	"timezoner/pkg/zone"
)

// ErrNoZonesProvided se retorna cuando no se proporcionan zonas horarias.
var ErrNoZonesProvided = errors.New("overlap: se requiere al menos una zona horaria")

// WorkingHours define la jornada laboral (ej: 9 a 18).
type WorkingHours struct {
	StartHour int `json:"start_hour"`
	EndHour   int `json:"end_hour"`
}

// DefaultWorkingHours devuelve el horario laboral habitual (09:00 a 17:00).
func DefaultWorkingHours() WorkingHours {
	return WorkingHours{StartHour: 9, EndHour: 17}
}

// Request define los parámetros para encontrar solapamiento de horarios laborales.
type Request struct {
	Date         time.Time               // Fecha base
	Zones        []string                // Lista de zonas a comparar
	CustomHours  map[string]WorkingHours // Horarios laborales personalizados por zona
	DefaultHours WorkingHours            // Horario por defecto
	SlotDuration time.Duration           // Duración mínima de cada slot
}

// Slot representa una ventana horaria donde todos los participantes están en horario laboral.
type Slot struct {
	StartTimeUTC time.Time            `json:"start_time_utc"`
	EndTimeUTC   time.Time            `json:"end_time_utc"`
	Duration     time.Duration        `json:"duration"`
	ZoneTimes    map[string]time.Time `json:"zone_times"`
}

// Find busca intervalos de tiempo durante un día donde todas las zonas coinciden en horario hábil.
func Find(req Request) ([]Slot, error) {
	if len(req.Zones) == 0 {
		return nil, ErrNoZonesProvided
	}

	slotDur := req.SlotDuration
	if slotDur <= 0 {
		slotDur = time.Hour
	}

	defHours := req.DefaultHours
	if defHours.StartHour == 0 && defHours.EndHour == 0 {
		defHours = DefaultWorkingHours()
	}

	type zoneRule struct {
		name string
		loc  *time.Location
		work WorkingHours
	}

	rules := make([]zoneRule, 0, len(req.Zones))
	for _, z := range req.Zones {
		loc, err := zone.LoadLocation(z)
		if err != nil {
			return nil, err
		}
		wh := defHours
		if custom, ok := req.CustomHours[z]; ok {
			wh = custom
		}
		rules = append(rules, zoneRule{
			name: z,
			loc:  loc,
			work: wh,
		})
	}

	year, month, day := req.Date.Date()
	baseLoc := rules[0].loc
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, baseLoc)
	endOfDay := startOfDay.Add(24 * time.Hour)

	var validSlots []Slot
	step := 30 * time.Minute

	for current := startOfDay; current.Before(endOfDay); current = current.Add(step) {
		allMatch := true
		zoneTimes := make(map[string]time.Time, len(rules))

		for _, r := range rules {
			local := current.In(r.loc)
			hour := local.Hour()
			if hour < r.work.StartHour || hour >= r.work.EndHour {
				allMatch = false
				break
			}
			zoneTimes[r.name] = local
		}

		if allMatch {
			validSlots = append(validSlots, Slot{
				StartTimeUTC: current.UTC(),
				EndTimeUTC:   current.Add(step).UTC(),
				Duration:     step,
				ZoneTimes:    zoneTimes,
			})
		}
	}

	return mergeConsecutiveSlots(validSlots, slotDur), nil
}

func mergeConsecutiveSlots(slots []Slot, minDuration time.Duration) []Slot {
	if len(slots) == 0 {
		return nil
	}

	var merged []Slot
	current := slots[0]

	for i := 1; i < len(slots); i++ {
		next := slots[i]
		if current.EndTimeUTC.Equal(next.StartTimeUTC) {
			current.EndTimeUTC = next.EndTimeUTC
			current.Duration = current.EndTimeUTC.Sub(current.StartTimeUTC)
		} else {
			if current.Duration >= minDuration {
				merged = append(merged, current)
			}
			current = next
		}
	}

	if current.Duration >= minDuration {
		merged = append(merged, current)
	}

	return merged
}
