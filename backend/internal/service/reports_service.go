package service

import (
	"context"
	"fmt"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/models"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/persistence"
	logger "github.com/SOAT1StackGoLang/Hackaton/pkg/middleware"
	"time"
)

type rS struct {
	tR persistence.TimekeepingRepository
}

func (r rS) GetReportByReferenceDateAndUserID(ctx context.Context, userID string, referenceDate time.Time) (*models.Timekeeping, error) {
	out, err := r.tR.GetTimekeepingByReferenceDateAndUserID(ctx, userID, referenceDate)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s", "failed getting timekeeping", err.Error()))
		return nil, err
	}
	return out, nil
}

func (r rS) GetReportByRangeAndUserID(ctx context.Context, userID string, start, end time.Time) (*models.RangedTimekeepingReport, error) {
	//TODO implement me
	panic("implement me")
}

func NewReportService(tR persistence.TimekeepingRepository) ReportService {
	return &rS{tR: tR}
}
