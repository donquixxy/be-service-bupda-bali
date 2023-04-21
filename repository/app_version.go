package repository

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type AppVersionRepositoryInterface interface {
	FindVersionByOS(db *gorm.DB, os string) ([]entity.AppVersion, error)
}

type AppVersionRepositoryImplementation struct {
	DB *config.Database
}

func NewAppVersionRepository(
	db *config.Database,
) AppVersionRepositoryInterface {
	return &AppVersionRepositoryImplementation{
		DB: db,
	}
}

func (repository *AppVersionRepositoryImplementation) FindVersionByOS(db *gorm.DB, os string) ([]entity.AppVersion, error) {
	appVersion := []entity.AppVersion{}
	results := db.
		Where("app_version.os = ?", os).
		Find(&appVersion)
	return appVersion, results.Error
}
