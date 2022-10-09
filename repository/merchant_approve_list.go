package repository

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type MerchantRepositoryInterface interface {
	CreateMerchantApproveList(db *gorm.DB, cart *entity.MerchantApproveList) error
	FindMerchantStatusApproveByUser(db *gorm.DB, idUser string) (*entity.MerchantApproveList, error)
}

type MerchantRepositoryImplementation struct {
	DB *config.Database
}

func NewMerchantRepository(
	db *config.Database,
) MerchantRepositoryInterface {
	return &MerchantRepositoryImplementation{
		DB: db,
	}
}

func (repository *MerchantRepositoryImplementation) CreateMerchantApproveList(db *gorm.DB, MerchantApproveList *entity.MerchantApproveList) error {
	result := db.Create(&MerchantApproveList)
	return result.Error
}

func (repository *MerchantRepositoryImplementation) FindMerchantStatusApproveByUser(db *gorm.DB, idUser string) (*entity.MerchantApproveList, error) {
	merchantApproveList := &entity.MerchantApproveList{}
	result := db.Where("id_user = ?", idUser).Find(merchantApproveList)
	return merchantApproveList, result.Error
}
