package persistence

import (
	"context"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/models"
	kitlog "github.com/go-kit/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

const entriesTable = "entries"

type pP struct {
	db  *gorm.DB
	log kitlog.Logger
}

func (p *pP) ListEntriesByRangeAndUserID(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]*models.Entry, error) {
	var (
		err error

		entries     []*models.Entry
		entriesDB   []*Entry
		endClausule time.Time
	)

	endClausule = end.Add(24 * time.Hour)

	if err := p.db.WithContext(ctx).
		Table(entriesTable).
		Where("user_id = ? and entry_at >= ? and entry_at < ?", userID, start, endClausule).
		Find(&entriesDB).Error; err != nil {
		p.log.Log(
			"failed listing products",
			err,
		)
	}

	for _, entry := range entriesDB {
		entries = append(entries, &models.Entry{
			ID:        entry.ID,
			UserID:    entry.UserID,
			CreatedAt: entry.CreatedAt,
		})
	}

	return entries, err
}

func (p *pP) InsertEntry(ctx context.Context, userID uuid.UUID, entryAt time.Time) (*models.Entry, error) {
	in := &Entry{
		ID:        uuid.New(),
		UserID:    userID,
		CreatedAt: entryAt.UTC(),
	}

	if err := p.db.WithContext(ctx).Table(entriesTable).Create(in).Error; err != nil {
		p.log.Log(
			"failed creating entry",
			err,
		)
		return nil, err
	}

	return &models.Entry{
		ID:        in.ID,
		UserID:    in.UserID,
		CreatedAt: in.CreatedAt,
	}, nil
}

func NewEntriesPersistence(db *gorm.DB, log kitlog.Logger) EntryRepository {
	return &pP{
		db:  db,
		log: log,
	}
}
