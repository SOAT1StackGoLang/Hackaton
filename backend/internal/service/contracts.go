package service

import (
	"context"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/models"
	"github.com/google/uuid"
	"time"
)

type EntriesService interface {
	CreateEntry(ctx context.Context, in *models.Entry) (*models.Entry, error)
}

type ReportsService interface {
	GetDailyReportFromDate(ctx context.Context, userID uuid.UUID, date time.Time) (*models.DailyReport, error)
	GetMonthlyFromDateReport(ctx context.Context, userID uuid.UUID, date time.Time) (*models.Report, error)
}
