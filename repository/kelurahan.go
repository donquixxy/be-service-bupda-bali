package repository

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type KelurahanRepositoryInterface interface {
	FindKelurahanByIdKeca(db *gorm.DB, idKeca int) ([]entity.Kelurahan, error)
}

type KelurahanRepositoryImplementation struct {
	DB *config.Database
}

func NewKelurahanRepository(
	db *config.Database,
) KelurahanRepositoryInterface {
	return &KelurahanRepositoryImplementation{
		DB: db,
	}
}

func (service *KelurahanRepositoryImplementation) FindKelurahanByIdKeca(db *gorm.DB, idKeca int) ([]entity.Kelurahan, error) {
	kelurahans := []entity.Kelurahan{}
	results := db.
		Where("idkeca = ?", idKeca).
		Find(&kelurahans)
	return kelurahans, results.Error
}
