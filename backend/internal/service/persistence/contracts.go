package persistence

import (
	"context"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/models"
	"time"
)

type TimekeepingRepository interface {
	GetTimekeepingByReferenceDateAndUserID(ctx context.Context, userID string, referenceDate time.Time) (*models.Timekeeping, error)
	// GetTimekeepingFromDate returns a timekeeping record for a given date
	ListTimekeepingByRangeAndUserID(ctx context.Context, userID string, start, end time.Time) ([]*models.Timekeeping, error)
	// CreateTimekeeping
	CreateTimekeeping(ctx context.Context, in *models.Timekeeping) error
	// UpdateTimekeeping
	UpdateTimekeeping(ctx context.Context, in *models.Timekeeping) (*models.Timekeeping, error)
}
