package service

import (
	"context"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/models"
	"time"
)

type TimekeepingService interface {
	InsertEntry(ctx context.Context, userID string, instant time.Time) (*models.Timekeeping, error)
	GetTimekeepingByReferenceDateAndUserID(ctx context.Context, userID string, referenceDate time.Time) (*models.Timekeeping, error)
	GetTimekeepingByRangeAndUserID(ctx context.Context, userID string, start, end time.Time) (*models.RangedTimekeepingReport, error)
}
