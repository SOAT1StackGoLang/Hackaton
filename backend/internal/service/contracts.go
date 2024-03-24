package service

import (
	"context"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/models"
	"time"
)

type TimekeepingService interface {
	InsertEntry(ctx context.Context, userID string, instant time.Time) (*models.Timekeeping, error)
}

type ReportService interface {
	GetReportByReferenceDateAndUserID(ctx context.Context, userID string, referenceDate time.Time) (*models.Timekeeping, error)
	GetReportByRangeAndUserID(ctx context.Context, userID string, start, end time.Time) (*models.RangedTimekeepingReport, error)
	GetReportCSVByRangeAndUserID(ctx context.Context, userID string, start, end time.Time) ([]byte, error)
}
