package endpoints

import (
	"fmt"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/models"
	"time"
)

type (
	// Entries
	InsertTimekeepingEntryRequest struct {
		UserID  string `json:"user_id"`
		EntryAt string `json:"entry_at"`
	}
	InsertEntryResponse struct {
		ID      string `json:"id"`
		UserID  string `json:"user_id"`
		EntryAt string `json:"entry_at"`
	}

	TimekeepingResponse struct {
		ID            string   `json:"id"`
		UserID        string   `json:"user_id,omitempty"`
		ReferenceDate string   `json:"created_at"`
		UpdatedAt     string   `json:"updated_at"`
		WorkedTime    string   `json:"worked_time"`
		Open          bool     `json:"open"`
		Details       []Detail `json:"details"`
	}

	Detail struct {
		WorkedTime    string `json:"worked_time"`
		StartingEntry Entry  `json:"starting_entry"`
		EndingEntry   *Entry `json:"ending_entry,omitempty"`
	}

	Entry struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
	}

	TimekeepingReportRequest struct {
		UserID string `json:"user_id"`
		Start  string `json:"from"`
		End    string `json:"to"`
	}

	TimekeepingReportByReferenceRequest struct {
		UserID        string `json:"user_id"`
		ReferenceDate string `json:"reference_date"`
	}

	TimekeepingReportResponse struct {
		UserID        string                `json:"user_id"`
		ReferenceDate string                `json:"reference_date"`
		Start         string                `json:"from"`
		End           string                `json:"to"`
		Open          string                `json:"open"`
		WorkedHours   string                `json:"worked_hours"`
		Report        []TimekeepingResponse `json:"report"`
	}

	//DailyReport struct {
	//	Date        string   `json:"date"`
	//	Open        bool     `json:"open"`
	//	WorkedHours string   `json:"worked_hours"`
	//	Details     []Detail `json:"details"`
	//}
)

func timeKeepingResponseFromModels(in *models.Timekeeping) TimekeepingResponse {
	location, _ := time.LoadLocation("America/Sao_Paulo") // Fuso horário de Brasília

	out := &TimekeepingResponse{
		ID:            in.ID.String(),
		UserID:        in.UserID,
		ReferenceDate: in.CreatedAt.In(location).Format("2006-01-02"),
		UpdatedAt:     in.UpdatedAt.In(location).Format(time.RFC3339),
		WorkedTime:    formatarDiferencaTempo(in.WorkedMinutes),
		Open:          in.Open,
	}

	for _, p := range in.Details {
		parseDetails(p, location, out)
	}

	return *out
}

func parseDetails(p *models.Details, location *time.Location, out *TimekeepingResponse) {
	add := Detail{
		WorkedTime: formatarDiferencaTempo(p.WorkedMinutes),
		StartingEntry: Entry{
			ID:        p.StartingEntry.ID.String(),
			CreatedAt: p.StartingEntry.CreatedAt.In(location).Format(time.RFC3339),
		},
	}
	if p.EndingEntry != nil {
		add.EndingEntry = &Entry{
			ID:        p.EndingEntry.ID.String(),
			CreatedAt: p.EndingEntry.CreatedAt.In(location).Format(time.RFC3339),
		}
	}
	out.Details = append(out.Details, add)
}

func formatarDiferencaTempo(minutos int64) string {
	// Extrair o número de horas e minutos
	horas := minutos / 60
	minutosRestantes := minutos % 60

	// Criar a string formatada
	stringFormatada := ""
	if horas > 0 {
		stringFormatada = fmt.Sprintf("%d hora(s)", horas)
		if minutosRestantes > 0 {
			stringFormatada += " e "
		}
	}
	if minutosRestantes > 0 {
		stringFormatada += fmt.Sprintf("%d minuto(s)", minutosRestantes)
	}

	return stringFormatada
}
