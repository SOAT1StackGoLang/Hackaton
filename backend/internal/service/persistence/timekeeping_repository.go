package persistence

import (
	"context"
	"database/sql"
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
	reg.UpdatedAt = sql.NullTime{
		Time:  time.Now().UTC(),
		Valid: true,
	}
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
		reg         Timekeeping
		location, _ = time.LoadLocation("America/Sao_Paulo") // Fuso horário de Brasília

	)

	begginningOfDay := referenceDate.UTC().Truncate(24 * time.Hour).In(location)
	endOfDay := begginningOfDay.Add(24 * time.Hour).In(location)

	if err := tP.db.WithContext(ctx).Table(timekeepingTable).
		Where("user_id = ? and created_at >= ? and created_at < ?", userID, begginningOfDay, endOfDay).
		First(&reg).Error; err != nil {
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
	dbIN := timekeepingFromModels(in)
	if err := tP.db.WithContext(ctx).Table(timekeepingTable).Create(dbIN).Error; err != nil {
		tP.log.Log(
			"failed creating timekeeping",
			err,
		)
	}

	in.ID = dbIN.ID

	return err
}

func (tP *tP) ListTimekeepingByRangeAndUserID(ctx context.Context, userID string, start, end time.Time) ([]*models.Timekeeping, error) {
	var (
		out         []*models.Timekeeping
		regs        []*Timekeeping
		location, _ = time.LoadLocation("America/Sao_Paulo") // Fuso horário de Brasília

	)

	begginningOfDay := start.Truncate(24 * time.Hour).In(location).UTC()
	endOfDay := end.Add(24 * time.Hour).In(location).UTC()

	if err := tP.db.WithContext(ctx).Table(timekeepingTable).
		Where("user_id = ? and created_at >= ? and created_at < ?", userID, begginningOfDay, endOfDay).
		Order("reference_date ASC").
		Find(&regs).Error; err != nil {
		tP.log.Log(
			"failed getting timekeeping",
			err,
		)
		return nil, err
	}

	for _, reg := range regs {
		out = append(out, reg.toModel())
	}

	return out, nil
}

func NewTimekeepingRepository(db *gorm.DB, log kitlog.Logger) TimekeepingRepository {
	return &tP{
		db:  db,
		log: log,
	}
}
