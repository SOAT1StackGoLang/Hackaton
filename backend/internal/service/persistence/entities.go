package persistence

import (
	"database/sql"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/models"
	"github.com/google/uuid"
	"time"
)

type (
	Timekeeping struct {
		ID            uuid.UUID    `gorm:"id,primaryKey" json:"id"`
		UserID        string       `json:"user_id"`
		CreatedAt     time.Time    `json:"created_at"`
		UpdatedAt     sql.NullTime `json:"updated_at"`
		WorkedMinutes int64        `json:"worked_minutes"`
		Open          bool         `json:"open"`
		Details       []Period     `gorm:"type:jsonb" json:"details"`
	}

	Period struct {
		WorkedMinutes int64 `json:"worked_minutes"`
		StartingEntry Entry `json:"starting_entry"`
		EndingEntry   Entry `json:"ending_entry"`
	}

	Entry struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
	}
)

func timekeepingFromModels(in *models.Timekeeping) *Timekeeping {
	var details []Period
	out := &Timekeeping{
		ID:            in.ID,
		UserID:        in.UserID,
		CreatedAt:     in.CreatedAt,
		WorkedMinutes: in.WorkedMinutes,
		Open:          in.Open,
	}

	if in.UpdatedAt.IsZero() {
		out.UpdatedAt = sql.NullTime{}
	} else {
		out.UpdatedAt = sql.NullTime{Time: in.UpdatedAt, Valid: true}
	}

	for _, detail := range in.Details {
		details = append(details, periodFromModel(detail))
	}

	out.Details = details

	return out
}

func periodFromModel(in models.Period) Period {
	out := Period{
		WorkedMinutes: in.WorkedMinutes,
		StartingEntry: Entry{
			ID:        in.StartingEntry.ID,
			CreatedAt: in.StartingEntry.CreatedAt,
		},
		EndingEntry: Entry{
			ID:        in.EndingEntry.ID,
			CreatedAt: in.EndingEntry.CreatedAt,
		},
	}

	return out
}

func (p *Period) toModel() *models.Period {
	out := models.Period{
		WorkedMinutes: p.WorkedMinutes,
		StartingEntry: models.Entry{
			ID:        p.StartingEntry.ID,
			CreatedAt: p.StartingEntry.CreatedAt,
		},
		EndingEntry: models.Entry{
			ID:        p.EndingEntry.ID,
			CreatedAt: p.EndingEntry.CreatedAt,
		},
	}

	return &out
}

func (t *Timekeeping) toModel() *models.Timekeeping {
	var details []models.Period
	out := &models.Timekeeping{
		ID:        t.ID,
		UserID:    t.UserID,
		CreatedAt: t.CreatedAt,
		Open:      t.Open,
	}

	if t.UpdatedAt.Valid {
		out.UpdatedAt = t.UpdatedAt.Time
	}

	for _, detail := range t.Details {
		details = append(details, models.Period{
			WorkedMinutes: detail.WorkedMinutes,
			StartingEntry: models.Entry{
				ID:        detail.StartingEntry.ID,
				CreatedAt: detail.StartingEntry.CreatedAt,
			},
			EndingEntry: models.Entry{
				ID:        detail.EndingEntry.ID,
				CreatedAt: detail.EndingEntry.CreatedAt,
			},
		})
	}

	out.Details = details

	return out
}
