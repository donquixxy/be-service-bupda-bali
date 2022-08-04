package repository

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type OrderItemPpobRepositoryInterface interface {
	CreateOrderItemPpob(db *gorm.DB, orderItemPpob *entity.OrderItemPpob) error
	FindOrderItemsPpobByIdOrder(db *gorm.DB, idOrder string) (*entity.OrderItemPpob, error)
}

type OrderItemPpobRepositoryImplementation struct {
	DB *config.Database
}

func NewOrderItemPpobRepository(
	db *config.Database,
) OrderItemPpobRepositoryInterface {
	return &OrderItemPpobRepositoryImplementation{
		DB: db,
	}
}

func (repository *OrderItemPpobRepositoryImplementation) CreateOrderItemPpob(db *gorm.DB, orderItemPpob *entity.OrderItemPpob) error {
	result := db.Create(&orderItemPpob)
	return result.Error
}

func (repository *OrderItemPpobRepositoryImplementation) FindOrderItemsPpobByIdOrder(db *gorm.DB, idOrder string) (*entity.OrderItemPpob, error) {
	orderItemsPpob := &entity.OrderItemPpob{}
	result := db.
		Find(orderItemsPpob, "id_order = ?", idOrder)
	return orderItemsPpob, result.Error
}
