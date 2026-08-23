package main

import (
	"fmt"
	"log"
	"time"

	"timezoner"
)

func main() {
	fmt.Println("=== DEMO 1: Conversión Directa ===")
	now := time.Now()

	// Convertir a varias zonas
	zones := []string{"America/Lima", "America/New_York", "Europe/Madrid", "Asia/Tokyo"}
	for _, z := range zones {
		localTime, err := timezoner.Convert(now, z)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		fmt.Printf("• %-20s: %s\n", z, localTime.Format("2006-01-02 15:04:05 MST (Z07:00)"))
	}

	fmt.Println("\n=== DEMO 2: API Fluida (Fluent API) ===")
	tp := timezoner.Now().In("Europe/Madrid")
	formatted, err := tp.Format("15:04:05 MST")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Hora actual en Madrid: %s\n", formatted)

	info, err := tp.Info()
	if err == nil {
		fmt.Printf("Detalles de Madrid -> Offset: %s, DST activo: %v\n", info.OffsetFormatted, info.IsDST)
	}

	fmt.Println("\n=== DEMO 3: Diferencia Horaria ===")
	diff, err := timezoner.Difference("Asia/Tokyo", "America/Lima", now)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Diferencia entre Tokio y Lima: %v horas\n", diff.Hours())

	fmt.Println("\n=== DEMO 4: Parsear string con zona y convertir ===")
	meetingTime, err := timezoner.ConvertBetween("2026-09-01 10:00", "2006-01-02 15:04", "America/Lima", "Europe/London")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Reunión de 10:00 (Lima) en Londres será a las: %s\n", meetingTime.Format("15:04 MST"))
}
