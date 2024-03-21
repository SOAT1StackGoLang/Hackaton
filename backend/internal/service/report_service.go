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
	persistence persistence.TimekeepingRepository
}

func (r *rS) GetDailyFromDate(ctx context.Context, userID uuid.UUID, date time.Time) (*models.DailyReport, error) {
	//TODO implement me
	panic("implement me")
}

func (r *rS) GetMonthlyReport(ctx context.Context, userID uuid.UUID, date time.Time) (*models.Report, error) {
	//TODO implement me
	panic("implement me")
}

func NewReportsService(persistence persistence.TimekeepingRepository, log kitlog.Logger) ReportsService {
	return &rS{
		log:         log,
		persistence: persistence,
	}
}
