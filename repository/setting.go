package repository

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type SettingRepositoryInterface interface {
	FindSettingShippingCost(db *gorm.DB, idDesa string) (*entity.Setting, error)
}

type SettingRepositoryImplementation struct {
	DB *config.Database
}

func NewSettingRepository(
	db *config.Database,
) SettingRepositoryInterface {
	return &SettingRepositoryImplementation{
		DB: db,
	}
}

func (service *SettingRepositoryImplementation) FindSettingShippingCost(db *gorm.DB, idDesa string) (*entity.Setting, error) {
	setting := &entity.Setting{}
	result := db.
		Where("settings_name = ?", "shipping_cost").
		Find(setting, "id_desa = ?", idDesa)
	return setting, result.Error
}
