package repository

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type PaymentHistoryRepositoryInterface interface {
	CreatePaymentHistory(db *gorm.DB, paymentHistory *entity.PaymentHistory) error
	FindPaymentHistoryById(db *gorm.DB, idUser, indexDate string) (*entity.PaymentHistory, error)
}

type PaymentHistoryRepositoryImplementation struct {
	DB *config.Database
}

func NewPaymentHistoryRepository(
	db *config.Database,
) PaymentHistoryRepositoryInterface {
	return &PaymentHistoryRepositoryImplementation{
		DB: db,
	}
}

func (service *PaymentHistoryRepositoryImplementation) CreatePaymentHistory(db *gorm.DB, paymentHistory *entity.PaymentHistory) error {
	results := db.Create(paymentHistory)
	return results.Error
}

func (service *PaymentHistoryRepositoryImplementation) FindPaymentHistoryById(db *gorm.DB, idUser, indexDate string) (*entity.PaymentHistory, error) {
	paymentHistory := &entity.PaymentHistory{}
	results := db.Where("id_user = ? AND index_date = ?", idUser, indexDate).First(paymentHistory)
	return paymentHistory, results.Error
}
