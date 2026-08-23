package humanize_test

import (
	"testing"
	"time"

	"timezoner/pkg/humanize"
)

func TestHumanize(t *testing.T) {
	now := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		target   time.Time
		expected string
		isEn     bool
	}{
		{now.Add(-10 * time.Second), "justo ahora", false},
		{now.Add(-60 * time.Second), "hace 1 minuto", false},
		{now.Add(-5 * time.Minute), "hace 5 minutos", false},
		{now.Add(-60 * time.Minute), "hace 1 hora", false},
		{now.Add(-2 * time.Hour), "hace 2 horas", false},
		{now.Add(-30 * time.Hour), "ayer", false},
		{now.Add(-5 * 24 * time.Hour), "hace 5 días", false},
		{now.Add(-40 * 24 * time.Hour), "hace 1 mes", false},
		{now.Add(-100 * 24 * time.Hour), "hace 3 meses", false},
		{now.Add(-400 * 24 * time.Hour), "hace 1 año", false},
		{now.Add(-800 * 24 * time.Hour), "hace 2 años", false},

		// Futuro
		{now.Add(10 * time.Second), "en unos momentos", false},
		{now.Add(60 * time.Second), "en 1 minuto", false},
		{now.Add(10 * time.Minute), "en 10 minutos", false},
		{now.Add(60 * time.Minute), "en 1 hora", false},
		{now.Add(3 * time.Hour), "en 3 horas", false},
		{now.Add(36 * time.Hour), "mañana", false},
		{now.Add(5 * 24 * time.Hour), "en 5 días", false},
		{now.Add(40 * 24 * time.Hour), "en 1 mes", false},
		{now.Add(100 * 24 * time.Hour), "en 3 meses", false},
		{now.Add(400 * 24 * time.Hour), "en 1 año", false},
		{now.Add(800 * 24 * time.Hour), "en 2 años", false},

		// English
		{now.Add(-10 * time.Second), "just now", true},
		{now.Add(-60 * time.Second), "1 minute ago", true},
		{now.Add(-5 * time.Minute), "5 minutes ago", true},
		{now.Add(-60 * time.Minute), "1 hour ago", true},
		{now.Add(-2 * time.Hour), "2 hours ago", true},
		{now.Add(-30 * time.Hour), "yesterday", true},
		{now.Add(-5 * 24 * time.Hour), "5 days ago", true},
		{now.Add(-40 * 24 * time.Hour), "1 month ago", true},
		{now.Add(-100 * 24 * time.Hour), "3 months ago", true},
		{now.Add(-400 * 24 * time.Hour), "1 year ago", true},
		{now.Add(-800 * 24 * time.Hour), "2 years ago", true},

		{now.Add(10 * time.Second), "in a few moments", true},
		{now.Add(60 * time.Second), "in 1 minute", true},
		{now.Add(10 * time.Minute), "in 10 minutes", true},
		{now.Add(60 * time.Minute), "in 1 hour", true},
		{now.Add(3 * time.Hour), "in 3 hours", true},
		{now.Add(36 * time.Hour), "tomorrow", true},
		{now.Add(5 * 24 * time.Hour), "in 5 days", true},
		{now.Add(40 * 24 * time.Hour), "in 1 month", true},
		{now.Add(100 * 24 * time.Hour), "in 3 months", true},
		{now.Add(400 * 24 * time.Hour), "in 1 year", true},
		{now.Add(800 * 24 * time.Hour), "in 2 years", true},
	}

	for _, tc := range tests {
		var res string
		if tc.isEn {
			res = humanize.HumanizeEn(tc.target, now)
		} else {
			res = humanize.Humanize(tc.target, now)
		}

		if res != tc.expected {
			t.Errorf("Humanize(%v) = %q, esperado %q", tc.target, res, tc.expected)
		}
	}
}
