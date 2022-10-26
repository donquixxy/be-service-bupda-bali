package repository

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type SettingRepositoryInterface interface {
	FindSettingShippingCost(db *gorm.DB, idDesa string) (*entity.Setting, error)
	FindVerAppByOS(db *gorm.DB, os int) ([]entity.Setting, error)
	FindAndroidVersion(db *gorm.DB) (*entity.Setting, error)
	FindIosVersion(db *gorm.DB) (*entity.Setting, error)
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

func (service *SettingRepositoryImplementation) FindAndroidVersion(db *gorm.DB) (*entity.Setting, error) {
	setting := &entity.Setting{}
	result := db.
		Where("value = ?", 1).
		Find(setting)
	return setting, result.Error
}

func (service *SettingRepositoryImplementation) FindIosVersion(db *gorm.DB) (*entity.Setting, error) {
	setting := &entity.Setting{}
	result := db.
		Where("value = ?", 2).
		Find(setting)
	return setting, result.Error
}

func (service *SettingRepositoryImplementation) FindVerAppByOS(db *gorm.DB, os int) ([]entity.Setting, error) {
	setting := []entity.Setting{}
	result := db.
		Where("value = ?", os).
		Find(&setting)
	return setting, result.Error
}

func (service *SettingRepositoryImplementation) FindSettingShippingCost(db *gorm.DB, idDesa string) (*entity.Setting, error) {
	setting := &entity.Setting{}
	result := db.
		Where("settings_name = ?", "shipping_cost").
		Find(setting, "id_desa = ?", idDesa).
		Order("setting_name ASC")
	return setting, result.Error
}
