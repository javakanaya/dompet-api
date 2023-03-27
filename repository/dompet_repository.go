package repository

import (
	"context"
	"errors"
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
	GetDetailDompet(tx *gorm.DB, id uint64) (entity.Dompet, error)
	InviteToDompet(tx *gorm.DB, idDompet uint64, emailUser string) (entity.User, error)
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

	newDetail := entity.DetailUserDompet{
		UserID:   dompet.UserID,
		DompetID: dompet.ID,
	}
	r.db.Debug().Create(&newDetail)

	return dompet, nil
}

func (r *dompetRepository) GetDetailDompet(tx *gorm.DB, id uint64) (entity.Dompet, error) {
	var dompet entity.Dompet
	var err error
	if tx == nil {
		tx = r.db.Where("id = ?", id).Preload("ListUser").Preload("ListCatatanKeuangan").Take(&dompet)
		err = tx.Error
	} else {
		err = tx.Where("id = ?", id).Preload("ListUser").Preload("ListCatatanKeuangan").Take(&dompet).Error
	}

	if err != nil {
		return entity.Dompet{}, err
	}

	return dompet, nil

}

func (r *dompetRepository) InviteToDompet(tx *gorm.DB, idDompet uint64, emailUser string) (entity.User, error) {
	var dompet entity.Dompet
	var newUser entity.User

	checkDompet := r.db.Where("id = ?", idDompet).Take(&dompet) // cek dompet apakah ada
	if checkDompet.Error != nil {
		return entity.User{}, errors.New("dompet tidak tersedia")
	}

	checkUser := r.db.Where("email = ?", emailUser).Take(&newUser) // cek user dengan email tersebut apakah ada
	if checkUser.Error != nil {
		return entity.User{}, errors.New("email user tidak valid")
	}

	var detail []entity.DetailUserDompet
	r.db.Where("dompet_id = ?", idDompet).Find(&detail) // ambil dompet yang ingin ditambah user baru

	var UserID []uint64
	for _, cek := range detail { // dari dompet yang sebelumnya diambil, extract id user siapa saja yang ada pada dompet.
		UserID = append(UserID, cek.UserID) // ambil seluruh id user dan simpan pada array, untuk pengecekan duplikat
	}

	for _, cek := range UserID {
		if cek == newUser.ID {
			return entity.User{}, errors.New("user sudah ada pada dompet")
		}
	}

	r.db.Model(&dompet).Association("ListUser").Append(&newUser)

	return newUser, nil
}
