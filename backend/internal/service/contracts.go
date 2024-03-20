package service

import (
	"context"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/models"
	"github.com/google/uuid"
	"time"
)

type TimekeepingService interface {
	InsertEntry(ctx context.Context, in *models.Entry) (*models.Entry, error)
}

type ReportsService interface {
	GetTimekeepingFromDate(ctx context.Context, userID uuid.UUID, date time.Time) (*models.DailyReport, error)
	GetDailyReport(ctx context.Context, userID uuid.UUID, date time.Time) (Timekeeping, error)
	GetMonthlyReport(ctx context.Context, userID uuid.UUID, date time.Time) (*models.Timekeeping, error)
}
