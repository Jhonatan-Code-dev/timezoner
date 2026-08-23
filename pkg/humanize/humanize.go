package humanize

import (
	"fmt"
	"math"
	"time"
)

// Humanize convierte la diferencia temporal en una cadena legible en español ("hace 5 minutos", "en 2 horas").
func Humanize(t time.Time, relativeTo ...time.Time) string {
	now := time.Now()
	if len(relativeTo) > 0 {
		now = relativeTo[0]
	}

	diff := now.Sub(t)
	isPast := diff >= 0
	absDiff := diff
	if !isPast {
		absDiff = -diff
	}

	seconds := int(math.Round(absDiff.Seconds()))
	minutes := int(math.Round(absDiff.Minutes()))
	hours := int(math.Round(absDiff.Hours()))
	days := hours / 24
	months := days / 30
	years := days / 365

	if isPast {
		switch {
		case seconds < 45:
			return "justo ahora"
		case seconds < 90:
			return "hace 1 minuto"
		case minutes < 45:
			return fmt.Sprintf("hace %d minutos", minutes)
		case minutes < 90:
			return "hace 1 hora"
		case hours < 24:
			return fmt.Sprintf("hace %d horas", hours)
		case hours < 48:
			return "ayer"
		case days < 30:
			return fmt.Sprintf("hace %d días", days)
		case days < 60:
			return "hace 1 mes"
		case months < 12:
			return fmt.Sprintf("hace %d meses", months)
		case months < 24:
			return "hace 1 año"
		default:
			return fmt.Sprintf("hace %d años", years)
		}
	}

	switch {
	case seconds < 45:
		return "en unos momentos"
	case seconds < 90:
		return "en 1 minuto"
	case minutes < 45:
		return fmt.Sprintf("en %d minutos", minutes)
	case minutes < 90:
		return "en 1 hora"
	case hours < 24:
		return fmt.Sprintf("en %d horas", hours)
	case hours < 48:
		return "mañana"
	case days < 30:
		return fmt.Sprintf("en %d días", days)
	case days < 60:
		return "en 1 mes"
	case months < 12:
		return fmt.Sprintf("en %d meses", months)
	case months < 24:
		return "en 1 año"
	default:
		return fmt.Sprintf("en %d años", years)
	}
}

// HumanizeEn convierte la diferencia temporal en una cadena legible en inglés ("5 minutes ago", "in 2 hours").
func HumanizeEn(t time.Time, relativeTo ...time.Time) string {
	now := time.Now()
	if len(relativeTo) > 0 {
		now = relativeTo[0]
	}

	diff := now.Sub(t)
	isPast := diff >= 0
	absDiff := diff
	if !isPast {
		absDiff = -diff
	}

	seconds := int(math.Round(absDiff.Seconds()))
	minutes := int(math.Round(absDiff.Minutes()))
	hours := int(math.Round(absDiff.Hours()))
	days := hours / 24
	months := days / 30
	years := days / 365

	if isPast {
		switch {
		case seconds < 45:
			return "just now"
		case seconds < 90:
			return "1 minute ago"
		case minutes < 45:
			return fmt.Sprintf("%d minutes ago", minutes)
		case minutes < 90:
			return "1 hour ago"
		case hours < 24:
			return fmt.Sprintf("%d hours ago", hours)
		case hours < 48:
			return "yesterday"
		case days < 30:
			return fmt.Sprintf("%d days ago", days)
		case days < 60:
			return "1 month ago"
		case months < 12:
			return fmt.Sprintf("%d months ago", months)
		case months < 24:
			return "1 year ago"
		default:
			return fmt.Sprintf("%d years ago", years)
		}
	}

	switch {
	case seconds < 45:
		return "in a few moments"
	case seconds < 90:
		return "in 1 minute"
	case minutes < 45:
		return fmt.Sprintf("in %d minutes", minutes)
	case minutes < 90:
		return "in 1 hour"
	case hours < 24:
		return fmt.Sprintf("in %d hours", hours)
	case hours < 48:
		return "tomorrow"
	case days < 30:
		return fmt.Sprintf("in %d days", days)
	case days < 60:
		return "in 1 month"
	case months < 12:
		return fmt.Sprintf("in %d months", months)
	case months < 24:
		return "in 1 year"
	default:
		return fmt.Sprintf("in %d years", years)
	}
}
