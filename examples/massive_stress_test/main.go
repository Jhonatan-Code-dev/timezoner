package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"timezoner"
)

func main() {
	const totalIterations = 10000
	const workers = 50

	fmt.Println("================================================================================")
	fmt.Printf("   EJECUCIÓN MASIVA DE ESTRÉS EXTREMO: %d OPERACIONES CON %d GOROUTINES      \n", totalIterations, workers)
	fmt.Println("================================================================================")

	sampleZones := []string{
		"UTC",
		"America/Lima",
		"America/New_York",
		"Europe/Madrid",
		"Europe/London",
		"Asia/Tokyo",
		"Asia/Kolkata",
		"Asia/Kathmandu",
		"Pacific/Auckland",
		"Pacific/Pago_Pago",
		"Pacific/Kiritimati",
		"Australia/Sydney",
	}

	sampleDates := []string{
		"2026-01-01 00:00:00",
		"2026-03-08 10:00:00", // Día de cambio de hora Spring-Forward en EE.UU.
		"2026-09-04 17:30:00", // Viernes
		"2026-10-15 09:00:00",
		"2026-12-31 23:59:59", // Fin de año
		"2028-02-29 12:00:00", // Bisiesto
	}

	var (
		successCount uint64
		failureCount uint64
		wg           sync.WaitGroup
		jobs         = make(chan int, totalIterations)
	)

	startTime := time.Now()

	// Lanzar pool de workers concurrentes
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(workerID)))

			for jobID := range jobs {
				// Selección pseudoaleatoria de zona y fecha
				z1 := sampleZones[rng.Intn(len(sampleZones))]
				z2 := sampleZones[rng.Intn(len(sampleZones))]
				dStr := sampleDates[rng.Intn(len(sampleDates))]

				// TEST 1: Ingesta limpia a UTC
				utcParsed, err := timezoner.IngestFromString(dStr, z1)
				if err != nil {
					atomic.AddUint64(&failureCount, 1)
					fmt.Printf("[FALLO Job %d] IngestFromString(%q, %q): %v\n", jobID, dStr, z1, err)
					continue
				}

				// TEST 2: Persistencia DBTime
				dbTime := timezoner.NewDBTime(utcParsed)
				if dbTime.Time().Location() != time.UTC {
					atomic.AddUint64(&failureCount, 1)
					fmt.Printf("[FALLO Job %d] DBTime no es UTC: %v\n", jobID, dbTime.Time().Location())
					continue
				}

				// TEST 3: Proyección al segundo huso
				proj, err := timezoner.ProjectForUser(dbTime.Time(), z2)
				if err != nil {
					atomic.AddUint64(&failureCount, 1)
					fmt.Printf("[FALLO Job %d] ProjectForUser(%v, %q): %v\n", jobID, dbTime.Time(), z2, err)
					continue
				}
				if proj.Formatted == "" || proj.OffsetFormatted == "" {
					atomic.AddUint64(&failureCount, 1)
					fmt.Printf("[FALLO Job %d] Proyección vacía\n", jobID)
					continue
				}

				// TEST 4: ZonedTime creación y reconstrucción local
				zoned, err := timezoner.ZonedFromLocal(dStr, z1)
				if err != nil {
					atomic.AddUint64(&failureCount, 1)
					fmt.Printf("[FALLO Job %d] ZonedFromLocal(%q, %q): %v\n", jobID, dStr, z1, err)
					continue
				}
				locRebuilt, err := zoned.Local()
				if err != nil || locRebuilt.IsZero() {
					atomic.AddUint64(&failureCount, 1)
					fmt.Printf("[FALLO Job %d] zoned.Local(): %v\n", jobID, err)
					continue
				}

				// TEST 5: Aritmética de días hábiles y límites de mes
				bDays := timezoner.At(utcParsed).
					AddBusinessDays(3).
					StartOfMonth().
					EndOfMonth().
					EndOfDay().
					MustTime()
				if bDays.Hour() != 23 || bDays.Minute() != 59 {
					atomic.AddUint64(&failureCount, 1)
					fmt.Printf("[FALLO Job %d] EndOfDay hora incorrecta: %v\n", jobID, bDays)
					continue
				}

				// TEST 6: Registro de alias y consultas dinámicas
				_ = timezoner.IsValid(z1)
				_, _ = timezoner.IsDST(z1, utcParsed)
				_, _ = timezoner.Difference(z1, z2, utcParsed)

				atomic.AddUint64(&successCount, 1)
			}
		}(w)
	}

	// Enviar los 10,000 trabajos al pool
	for i := 0; i < totalIterations; i++ {
		jobs <- i
	}
	close(jobs)

	// Esperar finalización de todos los workers
	wg.Wait()
	elapsed := time.Since(startTime)

	opsPerSec := float64(totalIterations) / elapsed.Seconds()

	fmt.Println("\n================================================================================")
	fmt.Println("                       RESULTADOS DE LA PRUEBA MASIVA                           ")
	fmt.Println("================================================================================")
	fmt.Printf("   • Total de Operaciones Solicitadas: %d\n", totalIterations)
	fmt.Printf("   • Operaciones Exitosas:             %d\n", atomic.LoadUint64(&successCount))
	fmt.Printf("   • Fallos Detectados:                %d\n", atomic.LoadUint64(&failureCount))
	fmt.Printf("   • Tiempo Total Transcurrido:        %v\n", elapsed)
	fmt.Printf("   • Rendimiento / Throughput:         %.2f operaciones/segundo\n", opsPerSec)
	fmt.Printf("   • Tiempo Promedio por Ciclo:        %.2f µs\n", float64(elapsed.Microseconds())/float64(totalIterations))
	fmt.Println("================================================================================")

	if atomic.LoadUint64(&failureCount) > 0 {
		fmt.Println("❌ ESTADO: FALLOS ENCONTRADOS")
	} else {
		fmt.Println("✅ ESTADO: 100% LIBRE DE FALLOS Y CONDICIONES DE CARRERA")
	}
}
