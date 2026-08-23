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
		{now.Add(-5 * time.Minute), "hace 5 minutos", false},
		{now.Add(-2 * time.Hour), "hace 2 horas", false},
		{now.Add(10 * time.Minute), "en 10 minutos", false},
		{now.Add(36 * time.Hour), "mañana", false},
		{now.Add(-5 * time.Minute), "5 minutes ago", true},
		{now.Add(10 * time.Minute), "in 10 minutes", true},
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
