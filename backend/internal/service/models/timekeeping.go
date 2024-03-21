package models

import "time"
import "github.com/google/uuid"

type (
	Entry struct {
		ID        int
		UserID    uuid.UUID
		CreatedAt time.Time
	}

	Period struct {
		RegistryID uuid.UUID
		EntryTime  time.Time
		ExitTime   time.Time
		TotalTime  time.Duration
	}

	DailyReport struct {
		UserID     uuid.UUID
		ReportDate time.Time
		TotalTime  time.Duration
		Registries []Entry
	}

	Report struct {
		UserID          uuid.UUID
		ReportDate      time.Time
		TotalTime       time.Duration
		DetailedEntries []Entry
	}
)
