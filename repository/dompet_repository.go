package repository

import (
	"context"
	"dompet-api/entity"

	"gorm.io/gorm"
)

type dompetRepository struct {
	db *gorm.DB
}

type DompetRepository interface {
	// functional
	GetMyDompet(tx *gorm.DB, id uint64) (entity.User, error)
	InsertDompet(ctx context.Context, dompet entity.Dompet) (entity.Dompet, error)
}

func NewDompetRepository(db *gorm.DB) DompetRepository {
	return &dompetRepository{
		db: db,
	}
}

func (r *dompetRepository) GetMyDompet(tx *gorm.DB, id uint64) (entity.User, error) {
	var user entity.User
	var err error
	if tx == nil {
		tx = r.db.Where("id = ?", id).Preload("ListDompet").Take(&user)
		err = tx.Error
	} else {
		err = tx.Where("id = ?", id).Preload("ListDompet").Take(&user).Error
	}

	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *dompetRepository) InsertDompet(ctx context.Context, dompet entity.Dompet) (entity.Dompet, error) {
	if err := r.db.Create(&dompet).Error; err != nil {
		return entity.Dompet{}, err
	}
	return dompet, nil
}
