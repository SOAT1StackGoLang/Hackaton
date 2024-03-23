package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/models"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/persistence"
	logger "github.com/SOAT1StackGoLang/Hackaton/pkg/middleware"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

const referenceLayout = "2006-01-02"

type tS struct {
	tR persistence.TimekeepingRepository
}

var location, _ = time.LoadLocation("America/Sao_Paulo") // Fuso horário de Brasília

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
			ID:            uuid.New(),
			UserID:        userID,
			ReferenceDate: referenceDate.In(location),
			CreatedAt:     referenceDate,
			Open:          true,
			Details: []*models.Details{
				{
					StartingEntry: &models.Entry{
						ID:        uuid.New(),
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
		currentDetail := tK.Details[len(tK.Details)-1]
		currentDetail.WorkedMinutes = t.processWorkedMinutes(tK, referenceDate)
		currentDetail.EndingEntry = &models.Entry{
			ID:        uuid.New(),
			CreatedAt: referenceDate,
		}

		tK.Open = false
	} else {
		tK.Details = append(tK.Details, &models.Details{
			WorkedMinutes: 0,
			StartingEntry: &models.Entry{
				ID:        uuid.New(),
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

func (t *tS) processWorkedMinutes(tK *models.Timekeeping, referenceDate time.Time) int64 {
	points := len(tK.Details)
	currentDetail := tK.Details[points-1]

	minutesWorked := int64(referenceDate.Sub(currentDetail.StartingEntry.CreatedAt).Minutes())
	currentDetail.WorkedMinutes = minutesWorked

	tK.WorkedMinutes += minutesWorked
	return minutesWorked
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
