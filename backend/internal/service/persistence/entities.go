package persistence

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/models"
	logger "github.com/SOAT1StackGoLang/Hackaton/pkg/middleware"
	"github.com/google/uuid"
	"time"
)

type (
	Timekeeping struct {
		ID            uuid.UUID       `gorm:"id,primaryKey" json:"id"`
		UserID        string          `json:"user_id"`
		CreatedAt     time.Time       `json:"created_at"`
		ReferenceDate time.Time       `json:"reference_date"`
		UpdatedAt     sql.NullTime    `json:"updated_at"`
		WorkedMinutes int64           `json:"worked_minutes"`
		Open          bool            `json:"open"`
		Details       json.RawMessage `json:"periods" gorm:"type:jsonb"`
	}

	Details struct {
		WorkedMinutes int64  `json:"worked_minutes"`
		StartingEntry *Entry `json:"starting_entry"`
		EndingEntry   *Entry `json:"ending_entry,omitempty"`
	}

	Entry struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
	}
)

func timekeepingFromModels(in *models.Timekeeping) *Timekeeping {
	var details []Details
	out := &Timekeeping{
		ID:            in.ID,
		UserID:        in.UserID,
		CreatedAt:     in.CreatedAt,
		WorkedMinutes: in.WorkedMinutes,
		ReferenceDate: in.ReferenceDate,
		Open:          in.Open,
	}

	if in.UpdatedAt.IsZero() {
		out.UpdatedAt = sql.NullTime{Valid: false}
	} else {
		out.UpdatedAt = sql.NullTime{Time: in.UpdatedAt, Valid: true}
	}

	for _, detail := range in.Details {
		details = append(details, detailsFromModel(*detail))
	}

	outB, err := json.Marshal(details)
	if err != nil {
		logger.Error(fmt.Sprintf("failed marshalling details: %v", err))
	}
	out.Details = outB

	return out
}

func detailsFromModel(in models.Details) Details {
	out := Details{
		WorkedMinutes: in.WorkedMinutes,
		StartingEntry: &Entry{
			ID:        in.StartingEntry.ID,
			CreatedAt: in.StartingEntry.CreatedAt,
		},
	}
	if in.EndingEntry != nil {
		out.EndingEntry = &Entry{
			ID:        in.EndingEntry.ID,
			CreatedAt: in.EndingEntry.CreatedAt,
		}
	}

	return out
}

func (p *Details) toModel() *models.Details {
	if p == nil {
		return nil
	}
	out := models.Details{
		WorkedMinutes: p.WorkedMinutes,
		StartingEntry: &models.Entry{
			ID:        p.StartingEntry.ID,
			CreatedAt: p.StartingEntry.CreatedAt,
		},
	}
	if p.EndingEntry != nil {
		out.EndingEntry = &models.Entry{
			ID:        p.EndingEntry.ID,
			CreatedAt: p.EndingEntry.CreatedAt,
		}

	}

	return &out
}

func (t *Timekeeping) toModel() *models.Timekeeping {
	var details []Details
	out := &models.Timekeeping{
		ID:            t.ID,
		UserID:        t.UserID,
		CreatedAt:     t.CreatedAt,
		ReferenceDate: t.ReferenceDate,
		WorkedMinutes: t.WorkedMinutes,
		Open:          t.Open,
	}

	if t.UpdatedAt.Valid {
		out.UpdatedAt = t.UpdatedAt.Time
	}

	if err := json.Unmarshal(t.Details, &details); err != nil {
		logger.Error(fmt.Sprintf("failed unmarshalling details: %v", err))
	}

	pD := make([]*models.Details, 0)
	for _, detail := range details {
		pD = append(pD, detail.toModel())
	}
	out.Details = pD

	return out
}
