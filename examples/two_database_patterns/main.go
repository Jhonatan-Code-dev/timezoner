package main

import (
	"encoding/json"
	"fmt"

	"timezoner"
)

// PATRÓN 1: Pagos, Transacciones y Logs (90% de los casos) -> Solo UTC
type PaymentTransaction struct {
	ID        string           `json:"id"`
	Amount    float64          `json:"amount"`
	Currency  string           `json:"currency"`
	PaidAtUTC timezoner.DBTime `json:"paid_at_utc"` // 1 Columna: TIMESTAMPTZ en UTC
}

// PATRÓN 2: Citas Médicas, Vuelos y Calendario (Google Calendar) -> UTC + Zona IANA
type MedicalAppointment struct {
	ID          int                 `json:"id"`
	Doctor      string              `json:"doctor"`
	Patient     string              `json:"patient"`
	ScheduledAt timezoner.ZonedTime `json:"scheduled_at"` // 2 Columnas (o JSONB): UTC + "America/Lima"
}

func main() {
	fmt.Println("====================================================================")
	fmt.Println("  DEMO: Los 2 Patrones de Base de Datos de las Grandes Empresas     ")
	fmt.Println("====================================================================")

	// -------------------------------------------------------------
	// PATRÓN 1: TRANSACCIONES (Solo UTC)
	// -------------------------------------------------------------
	fmt.Println("\n🟢 PATRÓN 1: Transacción / Pago (Solo UTC)")
	pago := PaymentTransaction{
		ID:        "TX-99881",
		Amount:    250.00,
		Currency:  "USD",
		PaidAtUTC: timezoner.NowDBTime(),
	}

	pagoJSON, _ := json.MarshalIndent(pago, "", "  ")
	fmt.Printf("1. Guardado en Base de Datos (1 Columna UTC):\n%s\n", string(pagoJSON))

	// Al mostrarlo al cajero en Lima y al auditor en Madrid:
	vistaLima, _ := timezoner.ProjectForUser(pago.PaidAtUTC.Time, "America/Lima")
	vistaMadrid, _ := timezoner.ProjectForUser(pago.PaidAtUTC.Time, "Europe/Madrid")
	fmt.Println("2. Visualización adaptada a cada usuario:")
	fmt.Printf("   • Cajero en Lima ve:   %s (%s)\n", vistaLima.Formatted, vistaLima.OffsetFormatted)
	fmt.Printf("   • Auditor en Madrid ve: %s (%s)\n", vistaMadrid.Formatted, vistaMadrid.OffsetFormatted)

	// -------------------------------------------------------------
	// PATRÓN 2: CITAS FUTURAS Y CALENDARIOS (UTC + Zona IANA)
	// -------------------------------------------------------------
	fmt.Println("\n--------------------------------------------------------------------")
	fmt.Println("🟡 PATRÓN 2: Cita Médica / Evento Futuro (UTC + Zona IANA)")

	// El paciente en Lima agenda cita para el 1 de Octubre a las 10:00 AM
	citaZoned, _ := timezoner.ZonedFromLocal("2026-10-01 10:00", "America/Lima")

	cita := MedicalAppointment{
		ID:          1001,
		Doctor:      "Dra. García",
		Patient:     "Roberto Gomez",
		ScheduledAt: citaZoned,
	}

	citaJSON, _ := json.MarshalIndent(cita, "", "  ")
	fmt.Printf("1. Guardado en Base de Datos (UTC + Zona IANA):\n%s\n", string(citaJSON))

	// Consultamos el evento
	horaLocalOriginal, _ := cita.ScheduledAt.Local()
	fmt.Println("2. Reconstrucción de la cita:")
	fmt.Printf("   • Hora local en Lima:   %s (Garantizado siempre a las 10:00 AM)\n", horaLocalOriginal.Format("2006-01-02 15:04"))

	// Médico especialista conectado desde Tokio para teleconsulta:
	medicoTokio, _ := timezoner.ProjectForUser(cita.ScheduledAt.UTC.Time, "Asia/Tokyo")
	fmt.Printf("   • Teleconsulta en Tokio: %s (Hora de Japón: %s)\n", medicoTokio.Formatted, medicoTokio.OffsetFormatted)

	fmt.Println("\n====================================================================")
	fmt.Println("  Ambos patrones resueltos de forma 100% nativa con Timezoner       ")
	fmt.Println("====================================================================")
}
