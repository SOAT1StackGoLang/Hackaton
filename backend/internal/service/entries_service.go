package service

import (
	"context"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/models"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/persistence"
)

type tS struct {
	tR persistence.EntryRepository
}

func (t tS) CreateEntry(ctx context.Context, in *models.Entry) (*models.Entry, error) {
	return t.tR.InsertEntry(ctx, in.UserID, in.CreatedAt)
}

func NewEntriesService(persistence persistence.EntryRepository) EntriesService {
	return &tS{
		tR: persistence,
	}
}
