package persistence

import (
	"context"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/models"
	kitlog "github.com/go-kit/log"
	"gorm.io/gorm"
	"time"
)

const timekeepingTable = "timekeeping"

type tP struct {
	db  *gorm.DB
	log kitlog.Logger
}

func (tP *tP) UpdateTimekeeping(ctx context.Context, in *models.Timekeeping) (*models.Timekeeping, error) {
	reg := timekeepingFromModels(in)
	if err := tP.db.WithContext(ctx).Table(timekeepingTable).Save(reg).Error; err != nil {
		tP.log.Log(
			"failed updating timekeeping",
			err,
		)
		return nil, err
	}
	return reg.toModel(), nil
}

func (tP *tP) GetTimekeepingByReferenceDateAndUserID(ctx context.Context, userID string, referenceDate time.Time) (*models.Timekeeping, error) {
	var (
		reg *Timekeeping
	)

	begginningOfDay := referenceDate.UTC().Truncate(24 * time.Hour)
	endOfDay := begginningOfDay.Add(24 * time.Hour)

	if err := tP.db.WithContext(ctx).Table(timekeepingTable).
		Where("user_id = ? and created_at >= ? and created_at < ?", userID, begginningOfDay, endOfDay).
		Find(reg).Error; err != nil {
		tP.log.Log(
			"failed getting timekeeping",
			err,
		)
		return nil, err
	}

	return reg.toModel(), nil
}

func (tP *tP) CreateTimekeeping(ctx context.Context, in *models.Timekeeping) error {
	var err error
	if err := tP.db.WithContext(ctx).Table(timekeepingTable).Create(in); err != nil {
		tP.log.Log(
			"failed creating timekeeping",
			err,
		)
	}

	return err
}

func (tP *tP) ListTimekeepingByRangeAndUserID(ctx context.Context, userID string, start, end time.Time) ([]*models.Timekeeping, error) {
	//var (
	//	err error
	//
	//	entries     []*models.Entry
	//	entriesDB   []*Entry
	//	endClausule time.Time
	//)
	//
	//endClausule = end.Add(24 * time.Hour).UTC()
	//
	//if err := tP.db.WithContext(ctx).
	//	Table(timekeepingTable).
	//	Where("user_id = ? and entry_at >= ? and entry_at < ?", userID, start, endClausule).
	//	Find(&entriesDB).Error; err != nil {
	//	tP.log.Log(
	//		"failed listing products",
	//		err,
	//	)
	//}
	//
	//for _, entry := range entriesDB {
	//	entries = append(entries, &models.Entry{
	//		ID:        entry.ID,
	//		CreatedAt: entry.CreatedAt,
	//	})
	//}
	//
	//return entries, err
	return nil, nil
}

func NewTimekeepingRepository(db *gorm.DB, log kitlog.Logger) TimekeepingRepository {
	return &tP{
		db:  db,
		log: log,
	}
}
