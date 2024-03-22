package persistence

import (
	"context"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/models"
	"github.com/google/uuid"
	"time"
)

type EntryRepository interface {
	// CreateTimekeeping creates a new timekeeping record
	InsertEntry(ctx context.Context, userID uuid.UUID, entryAt time.Time) (*models.Entry, error)
	// GetTimekeepingFromDate returns a timekeeping record for a given date
	ListEntriesByRangeAndUserID(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]*models.Entry, error)
}
