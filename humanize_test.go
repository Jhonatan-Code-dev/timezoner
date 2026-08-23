package timezoner_test

import (
	"testing"
	"time"

	"timezoner"
)

func TestHumanize(t *testing.T) {
	now := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		target   time.Time
		expected string
		isEn     bool
	}{
		{now.Add(-10 * time.Second), "justo ahora", false},
		{now.Add(-5 * time.Minute), "hace 5 minutos", false},
		{now.Add(-2 * time.Hour), "hace 2 horas", false},
		{now.Add(-30 * time.Hour), "ayer", false},
		{now.Add(-5 * 24 * time.Hour), "hace 5 días", false},
		{now.Add(10 * time.Minute), "en 10 minutos", false},
		{now.Add(3 * time.Hour), "en 3 horas", false},
		{now.Add(36 * time.Hour), "mañana", false},
		{now.Add(7 * 24 * time.Hour), "en 7 días", false},

		// English
		{now.Add(-10 * time.Second), "just now", true},
		{now.Add(-5 * time.Minute), "5 minutes ago", true},
		{now.Add(-2 * time.Hour), "2 hours ago", true},
		{now.Add(10 * time.Minute), "in 10 minutes", true},
		{now.Add(36 * time.Hour), "tomorrow", true},
	}

	for _, tc := range tests {
		var res string
		if tc.isEn {
			res = timezoner.HumanizeEn(tc.target, now)
		} else {
			res = timezoner.Humanize(tc.target, now)
		}

		if res != tc.expected {
			t.Errorf("Humanize(%v) = %q, esperado %q (isEn: %v)", tc.target, res, tc.expected, tc.isEn)
		}
	}
}
