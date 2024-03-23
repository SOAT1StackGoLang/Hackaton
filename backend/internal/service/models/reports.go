package models

import "time"

type (
	TimekeepingRangedReport struct {
		UserID      string
		Start       time.Time
		End         time.Time
		Open        bool
		WorkedHours int64
		Report      []Timekeeping
	}
)
