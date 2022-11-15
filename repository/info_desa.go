package repository

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type InfoDesaRepositoryInterface interface {
	FindInfoDesaByIdDesa(db *gorm.DB, idDesa string) ([]entity.InfoDesa, error)
}

type InfoDesaRepositoryImplementation struct {
	DB *config.Database
}

func NewInfoDesaRepository(
	db *config.Database,
) InfoDesaRepositoryInterface {
	return &InfoDesaRepositoryImplementation{
		DB: db,
	}
}

func (service *InfoDesaRepositoryImplementation) FindInfoDesaByIdDesa(db *gorm.DB, idDesa string) ([]entity.InfoDesa, error) {
	infoDesas := []entity.InfoDesa{}
	results := db.
		Where("id_desa = ?", idDesa).
		Where("is_active = ?", "1").
		Order("created_at desc").
		Find(&infoDesas)
	return infoDesas, results.Error
}
