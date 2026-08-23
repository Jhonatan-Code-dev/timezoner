package main

import (
	"fmt"
	"log"

	timezoner "github.com/Jhonatan-Code-dev/timezoner"
)

func main() {
	fmt.Println("==========================================================")
	fmt.Println("  DEMO: Ciclo Completo de Vida de Fecha (Ingesta -> BD -> Entrega) ")
	fmt.Println("==========================================================")

	// PASO 1: Un usuario en Lima envía una fecha local ("2026-09-01 10:00:00")
	inputDate := "2026-09-01 10:00:00"
	sourceZone := "America/Lima"
	fmt.Printf("\n1. ENTRADA DEL CLIENTE:\n   Fecha ingresada: %s (Zona: %s)\n", inputDate, sourceZone)

	// PASO 2: INGESTA -> Normalización a Hora Mundial (UTC) para la BD
	dbReadyUTC, err := timezoner.IngestFromString(inputDate, sourceZone)
	if err != nil {
		log.Fatalf("Error al normalizar para BD: %v", err)
	}
	fmt.Printf("\n2. NORMALIZACIÓN PARA BASE DE DATOS:\n   Guardado en BD (UTC): %s (Instante universal)\n", dbReadyUTC.Format("2006-01-02 15:04:05 UTC"))

	// PASO 3: PROYECCIÓN -> Lectura desde la BD para diferentes usuarios en el mundo
	fmt.Println("\n3. CONSULTA Y PROYECCIÓN POR USUARIOS:")

	viewers := []string{
		"America/Lima",     // Quien creó el evento
		"America/New_York", // Participante en NY
		"Europe/Madrid",    // Participante en España
		"Asia/Tokyo",       // Participante en Japón
		"Asia/Kolkata",     // Participante en India (Offset +05:30)
	}

	projections, err := timezoner.ProjectBatchForUsers(dbReadyUTC, viewers)
	if err != nil {
		log.Fatalf("Error en proyección: %v", err)
	}

	for _, zone := range viewers {
		u := projections[zone]
		fmt.Printf("   • %-20s: %s | Offset: %s | DST: %-5v | ISO: %s\n",
			zone, u.Formatted, u.OffsetFormatted, u.IsDST, u.ISO8601)
	}

	fmt.Println("\n>> Conclusión: Todos los usuarios ven el evento en su hora exacta y asisten al mismo instante real.")
}
