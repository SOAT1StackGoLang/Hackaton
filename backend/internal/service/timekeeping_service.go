package service

import (
	"context"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/models"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/persistence"
)

type tS struct {
	tR *persistence.TimekeepingRepository
}

func (t tS) InsertEntry(ctx context.Context, in *models.Entry) (*models.Entry, error) {
	//TODO implement me
	panic("implement me")
}

func NewTimekeepingService(persistence *persistence.TimekeepingRepository) TimekeepingService {
	return &tS{
		tR: persistence,
	}
}
