package main

import (
	"fmt"
	"log"
	"time"

	"timezoner"
)

func main() {
	fmt.Println("==========================================================")
	fmt.Println("  Planificador de Reuniones para Equipos Internacionales  ")
	fmt.Println("==========================================================")

	targetDate := time.Date(2026, 10, 15, 0, 0, 0, 0, time.UTC)
	teamZones := []string{
		"America/Lima",        // América del Sur
		"America/New_York",    // Costa Este EE.UU.
		"Europe/Madrid",       // Europa
	}

	req := timezoner.OverlapRequest{
		Date:         targetDate,
		Zones:        teamZones,
		DefaultHours: timezoner.WorkingHours{StartHour: 9, EndHour: 18}, // 09:00 a 18:00
		SlotDuration: 1 * time.Hour,
	}

	slots, err := timezoner.FindOverlap(req)
	if err != nil {
		log.Fatalf("Error al buscar solapamiento: %v", err)
	}

	if len(slots) == 0 {
		fmt.Println("No se encontraron horas de solapamiento común para este equipo.")
		return
	}

	fmt.Printf("Se encontraron %d ventanas de reunión disponibles:\n\n", len(slots))

	for i, slot := range slots {
		fmt.Printf("Ventana #%d (Duración: %v):\n", i+1, slot.Duration)
		fmt.Printf("  • UTC:           %s - %s\n",
			slot.StartTimeUTC.Format("15:04"),
			slot.EndTimeUTC.Format("15:04"))

		for _, zone := range teamZones {
			startTime := slot.ZoneTimes[zone]
			endTime := startTime.Add(slot.Duration)
			fmt.Printf("  • %-15s: %s - %s\n",
				zone,
				startTime.Format("15:04"),
				endTime.Format("15:04"))
		}
		fmt.Println("----------------------------------------------------------")
	}
}
