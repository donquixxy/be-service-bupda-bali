package repository

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type BannerRepositoryInterface interface {
	FindBannerByDesa(db *gorm.DB, idDesa string) ([]entity.Banner, error)
	FindBannerAll(db *gorm.DB) ([]entity.Banner, error)
}

type BannerRepositoryImplementation struct {
	DB *config.Database
}

func NewBannerRepository(
	db *config.Database,
) BannerRepositoryInterface {
	return &BannerRepositoryImplementation{
		DB: db,
	}
}

func (service *BannerRepositoryImplementation) FindBannerByDesa(db *gorm.DB, idDesa string) ([]entity.Banner, error) {
	banners := []entity.Banner{}
	results := db.
		Where("id_desa = ?", "650d1076-7818-423f-9a3e-48bb7e6343c2").
		Or("id_desa = ?", idDesa).
		Where("is_active = ?", 1).
		Find(&banners)
	return banners, results.Error
}

func (service *BannerRepositoryImplementation) FindBannerAll(db *gorm.DB) ([]entity.Banner, error) {
	banners := []entity.Banner{}
	results := db.
		Where("id_desa = ?", "650d1076-7818-423f-9a3e-48bb7e6343c2").
		Where("is_active = ?", 1).
		Find(&banners)
	return banners, results.Error
}
