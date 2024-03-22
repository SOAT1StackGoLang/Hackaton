package service

import (
	"context"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/models"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/persistence"
	kitlog "github.com/go-kit/log"
	"github.com/google/uuid"
	"time"
)

type rS struct {
	log         kitlog.Logger
	persistence persistence.EntryRepository
}

func (r *rS) GetDailyReportFromDate(ctx context.Context, userID uuid.UUID, date time.Time) (*models.DailyReport, error) {
	_, err := r.persistence.ListEntriesByRangeAndUserID(ctx, userID, date, date)
	if err != nil {
		return nil, err
	}

	return nil, err
}

func (r *rS) GetMonthlyFromDateReport(ctx context.Context, userID uuid.UUID, date time.Time) (*models.Report, error) {
	//TODO implement me
	panic("implement me")
}

func NewReportsService(persistence persistence.EntryRepository, log kitlog.Logger) ReportsService {
	return &rS{
		log:         log,
		persistence: persistence,
	}
}
