package repository

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type PaymentQueueRepositoryInterface interface {
	CreatePaymentQueue(db *gorm.DB, paymentQueue []entity.PaymentQueue) error
	FindPaymentQueueByIdUser(db *gorm.DB, idUser string) ([]entity.PaymentQueue, error)
	UpdatePaymentQueueById(db *gorm.DB, id string, paymentQueueUpdate *entity.PaymentQueue) error
	UpdateFailedPaymentQueueById(db *gorm.DB, idUser string, paymentQueueUpdate *entity.PaymentQueue) error
}

type PaymentQueueRepositoryImplementation struct {
}

func NewPaymentQueueRepository() PaymentQueueRepositoryInterface {
	return &PaymentQueueRepositoryImplementation{}
}

func (repository *PaymentQueueRepositoryImplementation) CreatePaymentQueue(db *gorm.DB, paymentQueue []entity.PaymentQueue) error {
	result := db.Create(&paymentQueue)
	return result.Error
}

func (repository *PaymentQueueRepositoryImplementation) FindPaymentQueueByIdUser(db *gorm.DB, idUser string) ([]entity.PaymentQueue, error) {
	paymentQueue := []entity.PaymentQueue{}
	result := db.
		Find(&paymentQueue, "id_user = ?", idUser).
		Where("status = ?", 0)
	return paymentQueue, result.Error
}

func (repository *PaymentQueueRepositoryImplementation) UpdatePaymentQueueById(db *gorm.DB, id string, paymentQueueUpdate *entity.PaymentQueue) error {
	paymentQueue := &entity.PaymentQueue{}
	result := db.
		Model(paymentQueue).
		Where("id = ?", id).
		Updates(paymentQueueUpdate)
	return result.Error
}

func (repository *PaymentQueueRepositoryImplementation) UpdateFailedPaymentQueueById(db *gorm.DB, idUser string, paymentQueueUpdate *entity.PaymentQueue) error {
	paymentQueue := &entity.PaymentQueue{}
	result := db.
		Model(paymentQueue).
		Where("id_user = ?", idUser).
		Where("status = ?", 0).
		Updates(paymentQueueUpdate)
	return result.Error
}
