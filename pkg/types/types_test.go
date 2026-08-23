package types_test

import (
	"encoding/json"
	"testing"
	"time"

	"timezoner/pkg/types"
)

type Event struct {
	ID        int             `json:"id"`
	CreatedAt types.DBTime    `json:"created_at"`
	Scheduled types.ZonedTime `json:"scheduled"`
}

func TestTypes_DBTimeAndZonedTime(t *testing.T) {
	t0 := time.Date(2026, 9, 1, 15, 0, 0, 0, time.UTC)
	dbTime := types.NewDBTime(t0)

	zoned, err := types.NewZonedTime(t0, "America/Lima")
	if err != nil {
		t.Fatalf("NewZonedTime falló: %v", err)
	}

	ev := Event{
		ID:        1,
		CreatedAt: dbTime,
		Scheduled: zoned,
	}

	data, err := json.Marshal(ev)
	if err != nil {
		t.Fatalf("json.Marshal falló: %v", err)
	}

	var decoded Event
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal falló: %v", err)
	}

	if !decoded.CreatedAt.Equal(t0) {
		t.Errorf("DBTime no coincide tras decode")
	}
	if decoded.Scheduled.Zone != "America/Lima" {
		t.Errorf("ZonedTime zone esperada America/Lima, obtenida: %s", decoded.Scheduled.Zone)
	}
}
