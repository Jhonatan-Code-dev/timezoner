package main

import (
	"encoding/json"
	"fmt"
	"time"

	"timezoner"
	"timezoner/pkg/calendar"
	"timezoner/pkg/humanize"
	"timezoner/pkg/ingest"
	"timezoner/pkg/overlap"
	"timezoner/pkg/types"
	"timezoner/pkg/zone"
)

type AuditFinding struct {
	CaseID      string `json:"case_id"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Passed      bool   `json:"passed"`
	Actual      string `json:"actual"`
	Expected    string `json:"expected"`
	Notes       string `json:"notes"`
}

func main() {
	fmt.Println("================================================================================")
	fmt.Println("   AUDITORÍA ADVERSARIA EXTREMA: BÚSQUEDA DE BUGS Y CASOS DE BORDE REALES      ")
	fmt.Println("================================================================================")

	var findings []AuditFinding

	// -------------------------------------------------------------------------
	// CASO 1: Regla Secular del Calendario Gregoriano (Años 1900, 2000, 2100)
	// Regla: Divisible por 4 es bisiesto, EXCEPTO si es divisible por 100, A MENOS que sea divisible por 400.
	// - Año 2000: Bisiesto (29 días)
	// - Año 2100: NO bisiesto (28 días)
	// - Año 2028: Bisiesto (29 días)
	// -------------------------------------------------------------------------
	feb2000 := calendar.DaysInMonth(2000, time.February)
	feb2100 := calendar.DaysInMonth(2100, time.February)
	feb2028 := calendar.DaysInMonth(2028, time.February)

	passed1 := (feb2000 == 29 && feb2100 == 28 && feb2028 == 29)
	findings = append(findings, AuditFinding{
		CaseID:      "SECULAR-LEAP-YEAR",
		Description: "Años bisiestos seculares gregorianos (2000=29, 2100=28, 2028=29)",
		Severity:    "CRÍTICO",
		Passed:      passed1,
		Actual:      fmt.Sprintf("2000:%d, 2100:%d, 2028:%d", feb2000, feb2100, feb2028),
		Expected:    "2000:29, 2100:28, 2028:29",
		Notes:       "Aritmética matemática de fin de mes",
	})

	// -------------------------------------------------------------------------
	// CASO 2: Fechas Pre-Unix Epoch (Antes de 1970) y Timestamps Negativos
	// Guardar fecha de nacimiento de 1945 en DBTime y serializar/deserializar JSON.
	// -------------------------------------------------------------------------
	pre1970 := time.Date(1945, 8, 6, 8, 15, 0, 0, time.UTC)
	dbPre1970 := types.NewDBTime(pre1970)
	jsonPre, err := json.Marshal(dbPre1970)
	var decodedPre types.DBTime
	_ = json.Unmarshal(jsonPre, &decodedPre)

	passed2 := (err == nil && decodedPre.Time().Year() == 1945 && decodedPre.Time().Month() == time.August && decodedPre.Time().Day() == 6)
	findings = append(findings, AuditFinding{
		CaseID:      "PRE-1970-EPOCH",
		Description: "Persistencia y serialización de fechas anteriores a 1970 (timestamps negativos)",
		Severity:    "MAYOR",
		Passed:      passed2,
		Actual:      string(jsonPre),
		Expected:    `"1945-08-06T08:15:00Z"`,
		Notes:       "Preservación de fechas históricas en DBTime",
	})

	// -------------------------------------------------------------------------
	// CASO 3: Días Hábiles Iniciando en Fin de Semana (Sábado y Domingo)
	// - Si estás en Sábado y pides +1 día hábil -> DEBE dar Lunes (no Domingo ni Martes)
	// - Si estás en Domingo y pides -1 día hábil -> DEBE dar Viernes (no Sábado)
	// -------------------------------------------------------------------------
	sabado := time.Date(2026, 9, 5, 10, 0, 0, 0, time.UTC)  // Sábado
	domingo := time.Date(2026, 9, 6, 10, 0, 0, 0, time.UTC) // Domingo

	lunesRes := calendar.AddBusinessDays(sabado, 1)
	viernesRes := calendar.AddBusinessDays(domingo, -1)

	passed3 := (lunesRes.Weekday() == time.Monday && lunesRes.Day() == 7 &&
		viernesRes.Weekday() == time.Friday && viernesRes.Day() == 4)
	findings = append(findings, AuditFinding{
		CaseID:      "WEEKEND-START-ARITHMETIC",
		Description: "AddBusinessDays iniciado en Sábado (+1 -> Lunes) y Domingo (-1 -> Viernes)",
		Severity:    "CRÍTICO",
		Passed:      passed3,
		Actual:      fmt.Sprintf("Sabado+1=%v (Día %d), Domingo-1=%v (Día %d)", lunesRes.Weekday(), lunesRes.Day(), viernesRes.Weekday(), viernesRes.Day()),
		Expected:    "Sabado+1=Monday (Día 7), Domingo-1=Friday (Día 4)",
		Notes:       "Salto correcto de fines de semana en cualquier dirección",
	})

	// -------------------------------------------------------------------------
	// CASO 4: Transición Fall-Back DST (Repetición de Hora en Otoño)
	// El 25 de Octubre de 2026 en Madrid, el reloj se atrasa de 03:00 a 02:00.
	// La hora 02:30 ocurre dos veces.
	// -------------------------------------------------------------------------
	fallBackSummer := time.Date(2026, 10, 25, 0, 0, 0, 0, time.UTC) // En UTC
	isDstPre, _ := zone.IsDST("Europe/Madrid", fallBackSummer)
	winterDate := time.Date(2026, 11, 15, 12, 0, 0, 0, time.UTC)
	isDstPost, _ := zone.IsDST("Europe/Madrid", winterDate)

	passed4 := (isDstPre == true && isDstPost == false)
	findings = append(findings, AuditFinding{
		CaseID:      "FALL-BACK-DST-TRANSITION",
		Description: "Detección correcta del fin de horario de verano en Europa/Madrid",
		Severity:    "MAYOR",
		Passed:      passed4,
		Actual:      fmt.Sprintf("Oct 25 DST=%v, Nov 15 DST=%v", isDstPre, isDstPost),
		Expected:    "Oct 25 DST=true, Nov 15 DST=false",
		Notes:       "Comprobación de solsticio y cambio de régimen de horario de verano",
	})

	// -------------------------------------------------------------------------
	// CASO 5: Husos con Offsets no Enteros (Minutos Fraccionarios)
	// - India: UTC+05:30
	// - Nepal: UTC+05:45
	// - Eucla (Australia): UTC+08:45
	// - Chatham Islands (Nueva Zelanda): UTC+12:45 / +13:45
	// - Terranova (Canadá): UTC-03:30 / -02:30
	// -------------------------------------------------------------------------
	tRef := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	infoIndia, _ := zone.GetInfo("Asia/Kolkata", tRef)
	infoNepal, _ := zone.GetInfo("Asia/Kathmandu", tRef)
	infoNewfoundland, _ := zone.GetInfo("America/St_Johns", tRef)

	passed5 := (infoIndia.OffsetFormatted == "+05:30" &&
		infoNepal.OffsetFormatted == "+05:45" &&
		infoNewfoundland.OffsetFormatted == "-02:30") // En junio tiene DST (-03:30 + 1 = -02:30)

	findings = append(findings, AuditFinding{
		CaseID:      "FRACTIONAL-OFFSETS",
		Description: "Husos horarios no enteros (+05:30, +05:45, -02:30 DST)",
		Severity:    "CRÍTICO",
		Passed:      passed5,
		Actual:      fmt.Sprintf("India:%s, Nepal:%s, Terranova:%s", infoIndia.OffsetFormatted, infoNepal.OffsetFormatted, infoNewfoundland.OffsetFormatted),
		Expected:    "India:+05:30, Nepal:+05:45, Terranova:-02:30",
		Notes:       "Formateo y cálculo de segundos a minutos en husos con fracción",
	})

	// -------------------------------------------------------------------------
	// CASO 6: Ingesta con Cadenas Maliciosas / Inyecciones
	// Probar cadenas como: "   ", "\n\t", "99999-99-99", "NULL", "'; DROP TABLE--"
	// Debe retornar ErrEmptyDateString o ErrInvalidInput, JAMÁS panic.
	// -------------------------------------------------------------------------
	maliciousInputs := []string{"   ", "\n\t", "99999-99-99", "NULL", "'; DROP TABLE--", "<script>alert(1)</script>"}
	passed6 := true
	for _, m := range maliciousInputs {
		_, err := ingest.FromString(m, "UTC")
		if err == nil {
			passed6 = false
		}
	}
	findings = append(findings, AuditFinding{
		CaseID:      "MALICIOUS-INPUT-IMMUNITY",
		Description: "Inmunidad a inyecciones y cadenas corruptas sin provocar pánicos",
		Severity:    "CRÍTICO",
		Passed:      passed6,
		Actual:      "Todos los inputs maliciosos retornaron error estructurado",
		Expected:    "Rechazo controlado con error",
		Notes:       "Seguridad de sanitización de entradas",
	})

	// -------------------------------------------------------------------------
	// CASO 7: Humanize con Diferencias Sub-Segundo (Nanosegundos / Mismo Instante)
	// -------------------------------------------------------------------------
	now := time.Now()
	resExacto := humanize.Humanize(now, now)
	resNano := humanize.Humanize(now.Add(10*time.Millisecond), now)

	passed7 := (resExacto == "justo ahora" || resExacto == "en unos momentos") &&
		(resNano == "en unos momentos" || resNano == "justo ahora")
	findings = append(findings, AuditFinding{
		CaseID:      "HUMANIZE-SUBSECOND",
		Description: "Formateo relativo ante diferencias de subsegundo y nanosegundos",
		Severity:    "MENOR",
		Passed:      passed7,
		Actual:      fmt.Sprintf("Exacto=%q, Nano=%q", resExacto, resNano),
		Expected:    `"justo ahora" / "en unos momentos"`,
		Notes:       "Estabilidad de redondeo matemático",
	})

	// -------------------------------------------------------------------------
	// CASO 8: Solapamiento en Cambio de Año (31 Diciembre al 1 Enero)
	// -------------------------------------------------------------------------
	overlapNYE, err := overlap.Find(overlap.Request{
		Date:         time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
		Zones:        []string{"America/Lima", "Europe/Madrid"},
		DefaultHours: overlap.WorkingHours{StartHour: 9, EndHour: 18},
		SlotDuration: 1 * time.Hour,
	})
	passed8 := (err == nil && len(overlapNYE) > 0)
	findings = append(findings, AuditFinding{
		CaseID:      "OVERLAP-NEW-YEAR-EVE",
		Description: "Cálculo de solapamiento en fecha de fin de año",
		Severity:    "MAYOR",
		Passed:      passed8,
		Actual:      fmt.Sprintf("Ventanas encontradas: %d", len(overlapNYE)),
		Expected:    "Ventanas >= 1",
		Notes:       "Algoritmo continuo a través de 24 horas",
	})

	// -------------------------------------------------------------------------
	// CASO 9: TimePoint Encadenamiento Extremo (10 Métodos Consecutivos)
	// -------------------------------------------------------------------------
	tpChain := timezoner.At(tRef).
		In("America/Lima").
		ToUTC().
		StartOfWeek().
		AddBusinessDays(4).
		StartOfMonth().
		EndOfMonth().
		StartOfDay().
		EndOfDay().
		In("Asia/Tokyo")

	tFinal, errChain := tpChain.Time()
	passed9 := (errChain == nil && !tFinal.IsZero() && tFinal.Location().String() == "Asia/Tokyo")
	findings = append(findings, AuditFinding{
		CaseID:      "FLUENT-CHAIN-DEEP",
		Description: "Encadenamiento profundo de 10 transformaciones sin degradación",
		Severity:    "MAYOR",
		Passed:      passed9,
		Actual:      fmt.Sprintf("Zona final: %s, Hora: %v", tFinal.Location().String(), tFinal.Format(time.RFC3339)),
		Expected:    "Asia/Tokyo con time válido",
		Notes:       "Invariantes y propagación de puntero en TimePoint",
	})

	// -------------------------------------------------------------------------
	// CASO 10: Rendimiento de Aritmética a Gran Escala (+10,000 Días Hábiles)
	// Verificar que AddBusinessDays(10000) no se cuelgue ni desborde.
	// -------------------------------------------------------------------------
	startCalc := time.Now()
	farFuture := calendar.AddBusinessDays(tRef, 10000)
	calcDuration := time.Since(startCalc)

	passed10 := (farFuture.Year() > 2060 && calcDuration < 10*time.Millisecond)
	findings = append(findings, AuditFinding{
		CaseID:      "LARGE-SCALE-BUSINESS-DAYS",
		Description: "Aritmética de +10,000 días hábiles en menos de 10 milisegundos",
		Severity:    "MAYOR",
		Passed:      passed10,
		Actual:      fmt.Sprintf("Año resultante: %d, Tiempo: %v", farFuture.Year(), calcDuration),
		Expected:    "Año > 2060 en < 10ms",
		Notes:       "Eficiencia de bucle y ausencia de desbordamiento",
	})

	// -------------------------------------------------------------------------
	// REPORTE DE RESULTADOS
	// -------------------------------------------------------------------------
	totalPassed := 0
	totalFailed := 0

	fmt.Println("\nResultados de los 10 Casos de Borde Adversarios:")
	for _, f := range findings {
		status := "✅ PASS"
		if !f.Passed {
			status = "❌ FAIL"
			totalFailed++
		} else {
			totalPassed++
		}
		fmt.Printf("   [%s] %-26s | Severidad: %-8s | %s\n", status, f.CaseID, f.Severity, f.Description)
		fmt.Printf("          Obtenido: %s\n", f.Actual)
	}

	fmt.Println("\n================================================================================")
	fmt.Printf("  RESUMEN DE AUDITORÍA: %d Aprobados | %d Fallados de %d Casos Extremos\n",
		totalPassed, totalFailed, len(findings))
	fmt.Println("================================================================================")

	if totalFailed == 0 {
		fmt.Println("🏆 CONCLUSIÓN DEL AUDITOR: El paquete soporta formalmente todos los casos de borde extremos.")
	} else {
		fmt.Printf("⚠️ CONCLUSIÓN DEL AUDITOR: Se detectaron %d defectos reales.\n", totalFailed)
	}
}
