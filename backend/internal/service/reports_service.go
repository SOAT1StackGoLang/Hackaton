package service

import (
	"context"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/models"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/persistence"
	"time"
)

type rS struct {
	tR persistence.TimekeepingRepository
}

func (r rS) GetReportByReferenceDateAndUserID(ctx context.Context, userID string, referenceDate time.Time) (*models.Timekeeping, error) {
	//TODO implement me
	panic("implement me")
}

func (r rS) GetReportByRangeAndUserID(ctx context.Context, userID string, start, end time.Time) (*models.RangedTimekeepingReport, error) {
	//TODO implement me
	panic("implement me")
}

func NewReportService(tR persistence.TimekeepingRepository) ReportService {
	return &rS{tR: tR}
}
