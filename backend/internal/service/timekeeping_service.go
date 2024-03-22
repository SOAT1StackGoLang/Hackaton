package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/models"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/persistence"
	logger "github.com/SOAT1StackGoLang/Hackaton/pkg/middleware"
	"gorm.io/gorm"
	"time"
)

const referenceLayout = "2006-01-02"

type tS struct {
	tR persistence.TimekeepingRepository
}

func (t *tS) InsertEntry(ctx context.Context, userID string, instant time.Time) (*models.Timekeeping, error) {
	var referenceDate time.Time
	if instant.IsZero() {
		referenceDate = time.Now().UTC()
	} else {
		referenceDate = instant.UTC()
	}

	tK, err := t.tR.GetTimekeepingByReferenceDateAndUserID(ctx, userID, referenceDate)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		tK = &models.Timekeeping{
			UserID:    userID,
			CreatedAt: referenceDate,
			Open:      true,
			Details: []models.Period{
				{
					StartingEntry: models.Entry{
						CreatedAt: instant,
					},
				},
			},
		}
		if err := t.tR.CreateTimekeeping(ctx, tK); err != nil {
			logger.Error(fmt.Sprintf("failed creating timekeeping: %v", err))
			return nil, err
		}
		return tK, nil
	} else if err != nil {
		logger.Error(fmt.Sprintf("failed getting timekeeping: %v", err))
		return nil, err
	}

	if tK.Open {
		t.processWorkedMinutes(tK, referenceDate)

		tK.Open = false
	} else {
		tK.Details = append(tK.Details, models.Period{
			WorkedMinutes: 0,
			StartingEntry: models.Entry{
				CreatedAt: referenceDate,
			},
		})
		tK.Open = true
	}

	if tK, err = t.tR.UpdateTimekeeping(ctx, tK); err != nil {
		logger.Error(fmt.Sprintf("failed updating timekeeping: %v", err))
		return nil, err
	}

	return tK, nil
}

func (t *tS) processWorkedMinutes(tK *models.Timekeeping, referenceDate time.Time) {
	points := len(tK.Details)
	currentDetail := tK.Details[points-1]

	minutesWorked := int64(referenceDate.Sub(currentDetail.StartingEntry.CreatedAt).Minutes())
	currentDetail.WorkedMinutes = minutesWorked

	tK.WorkedMinutes += minutesWorked
}

func (t *tS) GetTimekeepingByReferenceDateAndUserID(ctx context.Context, userID string, referenceDate time.Time) (*models.Timekeeping, error) {
	//TODO implement me
	panic("implement me")
}

func (t *tS) GetTimekeepingByRangeAndUserID(ctx context.Context, userID string, start, end time.Time) (*models.RangedTimekeepingReport, error) {
	//TODO implement me
	panic("implement me")
}

func NewTimekeepingService(persistence persistence.TimekeepingRepository) TimekeepingService {
	return &tS{
		tR: persistence,
	}
}
