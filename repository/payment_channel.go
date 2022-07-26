package repository

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type PaymentChannelRepositoryInterface interface {
	FindPaymentChannelByCode(db *gorm.DB, code string) (*entity.PaymentChannel, error)
	FindPaymentChannel(db *gorm.DB) ([]entity.PaymentChannel, error)
}

type PaymentChannelRepositoryImplementation struct {
	DB *config.Database
}

func NewPaymentChannelRepository(
	db *config.Database,
) PaymentChannelRepositoryInterface {
	return &PaymentChannelRepositoryImplementation{
		DB: db,
	}
}

func (repository *PaymentChannelRepositoryImplementation) FindPaymentChannelByCode(db *gorm.DB, code string) (*entity.PaymentChannel, error) {
	paymentChannel := &entity.PaymentChannel{}
	result := db.
		Find(paymentChannel, "code = ? AND is_active = ?", code, 1)
	return paymentChannel, result.Error
}

func (repository *PaymentChannelRepositoryImplementation) FindPaymentChannel(db *gorm.DB) ([]entity.PaymentChannel, error) {
	paymentChannel := []entity.PaymentChannel{}
	result := db.
		Joins("PaymentMethod").
		Find(&paymentChannel, "payment_channel.is_active = ?", 1)
	return paymentChannel, result.Error
}
