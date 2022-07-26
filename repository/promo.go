package repository

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type PromoRepositoryInterface interface {
	FindPromo(db *gorm.DB, IdDesa string) ([]entity.Promo, error)
}

type PromoRepositoryImplementation struct {
	DB *config.Database
}

func NewPromoRepository(
	db *config.Database,
) PromoRepositoryInterface {
	return &PromoRepositoryImplementation{
		DB: db,
	}
}

func (service *PromoRepositoryImplementation) FindPromo(db *gorm.DB, IdDesa string) ([]entity.Promo, error) {
	promos := []entity.Promo{}
	results := db.
		Where("id_desa = ?", IdDesa).
		Find(&promos)
	return promos, results.Error
}
