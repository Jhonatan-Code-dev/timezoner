package main

import (
	"encoding/json"
	"fmt"
	"time"

	timezoner "github.com/Jhonatan-Code-dev/timezonermax"
)

type InvoiceRecord struct {
	ID         string           `json:"id"`
	Amount     float64          `json:"amount"`
	IssuedAt   timezoner.DBTime `json:"issued_at"`
	DueAt      timezoner.DBTime `json:"due_at"`
	CustomerTZ string           `json:"customer_tz"`
}

func main() {
	fmt.Println("==================================================================")
	fmt.Println("  TIMEZONER: Enterprise Showcase (Capacidades de Élite)           ")
	fmt.Println("==================================================================")

	// 1. Aritmética de Negocio: Calcular vencimiento a 5 días hábiles a partir de un viernes
	issueDate := time.Date(2026, 9, 4, 10, 30, 0, 0, time.UTC) // Viernes
	dueDate := timezoner.At(issueDate).
		AddBusinessDays(5). // Salta automáticamente sábado y domingo -> Viernes 11
		EndOfDay().         // Establece 23:59:59
		MustTime()

	invoice := InvoiceRecord{
		ID:         "INV-2026-001",
		Amount:     1450.00,
		IssuedAt:   timezoner.NewDBTime(issueDate),
		DueAt:      timezoner.NewDBTime(dueDate),
		CustomerTZ: "America/Lima",
	}

	// 2. Serialización JSON / SQL Driver
	jsonData, _ := json.MarshalIndent(invoice, "", "  ")
	fmt.Printf("\n1. REGISTRO PERSISTIBLE EN BD (JSON / SQL):\n%s\n", string(jsonData))

	// 3. Tiempo Relativo Humano (Humanize)
	relativeIssue := timezoner.Humanize(issueDate, issueDate.Add(2*time.Hour+15*time.Minute))
	fmt.Printf("\n2. TIEMPO RELATIVO HUMANO:\n   • Emisión: %s\n", relativeIssue)

	// 4. Ingesta y Proyección para el Cliente
	customerView, _ := timezoner.ProjectForUser(invoice.DueAt.Time(), invoice.CustomerTZ)
	fmt.Printf("\n3. VISTA PARA EL CLIENTE EN %s:\n", invoice.CustomerTZ)
	fmt.Printf("   • Vence el: %s (%s)\n", customerView.Formatted, customerView.OffsetFormatted)
	fmt.Printf("   • DST activo: %v\n", customerView.IsDST)

	// 5. Ingesta limpia
	parsedInput, _ := timezoner.IngestFromString("11/09/2026 23:59", "PET")
	fmt.Printf("\n4. INGESTA DESDE FORMULARIO (PET -> UTC):\n   • %s\n", parsedInput.Format(time.RFC3339))

	fmt.Println("\n==================================================================")
	fmt.Println("  Demostración completada con éxito. Totalmente portable y robusto.")
	fmt.Println("==================================================================")
}
