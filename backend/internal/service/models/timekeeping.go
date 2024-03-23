package models

import "time"
import "github.com/google/uuid"

type (
	DailyReport struct {
		UserID     string
		ReportDate time.Time
		TotalTime  time.Duration
		Registries []Entry
	}

	Report struct {
		UserID          string
		ReportDate      time.Time
		TotalTime       time.Duration
		DetailedEntries []Entry
	}

	RangedTimekeepingReport struct {
		UserID        string
		Start         time.Time
		End           time.Time
		WorkedMinutes int64
		Open          bool
		Details       []Timekeeping
	}

	Timekeeping struct {
		ID            uuid.UUID
		UserID        string
		CreatedAt     time.Time
		ReferenceDate time.Time
		UpdatedAt     time.Time
		WorkedMinutes int64
		Open          bool
		Details       []*Details
	}

	Details struct {
		WorkedMinutes int64
		StartingEntry *Entry
		EndingEntry   *Entry
	}

	Entry struct {
		ID        uuid.UUID
		CreatedAt time.Time
	}
)
