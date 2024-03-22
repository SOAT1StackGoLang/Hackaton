package endpoints

type (
	// Entries
	InsertEntryRequest struct {
		UserID  string `json:"user_id"`
		EntryAt string `json:"entry_at"`
	}
	InsertEntryResponse struct {
		ID      string `json:"id"`
		UserID  string `json:"user_id"`
		EntryAt string `json:"entry_at"`
	}

	TimekeepingReportRequest struct {
		UserID string `json:"user_id"`
		From   string `json:"from"`
		To     string `json:"to"`
	}

	TimekeepingReportResponse struct {
		UserID      string        `json:"user_id"`
		From        string        `json:"from"`
		To          string        `json:"to"`
		Status      string        `json:"status"`
		WorkedHours string        `json:"worked_hours"`
		Report      []DailyReport `json:"report"`
	}

	DailyReport struct {
		Date        string  `json:"date"`
		Open        bool    `json:"open"`
		WorkedHours string  `json:"worked_hours"`
		Entries     []Entry `json:"entries"`
	}

	Entry struct {
		ID             string `json:"id"`
		UserID         string `json:"user_id"`
		OpeningEntryAt string `json:"opening_entry_at"`
		ClosingEntryAt string `json:"closing_entry_at"`
		WorkedHours    string `json:"worked_hours"`
	}
)
