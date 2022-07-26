package repository

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type KecamatanRepositoryInterface interface {
	FindKecamatan(db *gorm.DB) ([]entity.Kecamatan, error)
}

type KecamatanRepositoryImplementation struct {
	DB *config.Database
}

func NewKecamatanRepository(
	db *config.Database,
) KecamatanRepositoryInterface {
	return &KecamatanRepositoryImplementation{
		DB: db,
	}
}

func (service *KecamatanRepositoryImplementation) FindKecamatan(db *gorm.DB) ([]entity.Kecamatan, error) {
	kecamatans := []entity.Kecamatan{}
	results := db.
		Where("idprop = ?", 17).
		Find(&kecamatans)
	return kecamatans, results.Error
}
