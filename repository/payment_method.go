package repository

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type PaymentMethodRepositoryInterface interface {
	FindPaymentMethodByCode(db *gorm.DB, code string) (*entity.PaymentMethod, error)
}

type PaymentMethodRepositoryImplementation struct {
	DB *config.Database
}

func NewPaymentMethodRepository(
	db *config.Database,
) PaymentMethodRepositoryInterface {
	return &PaymentMethodRepositoryImplementation{
		DB: db,
	}
}

func (repository *PaymentMethodRepositoryImplementation) FindPaymentMethodByCode(db *gorm.DB, code string) (*entity.PaymentMethod, error) {
	paymentMethod := &entity.PaymentMethod{}
	result := db.
		Find(paymentMethod, "code = ? AND is_active = ?", code, 1)
	return paymentMethod, result.Error
}
