package repository

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type OrderItemRepositoryInterface interface {
	CreateOrderItem(db *gorm.DB, orderItem []entity.OrderItem) error
	FindOrderItemsByIdOrder(db *gorm.DB, idOrder string) ([]entity.OrderItem, error)
}

type OrderItemRepositoryImplementation struct {
	DB *config.Database
}

func NewOrderItemRepository(
	db *config.Database,
) OrderItemRepositoryInterface {
	return &OrderItemRepositoryImplementation{
		DB: db,
	}
}

func (repository *OrderItemRepositoryImplementation) CreateOrderItem(db *gorm.DB, orderItem []entity.OrderItem) error {
	result := db.Create(&orderItem)
	return result.Error
}

func (repository *OrderItemRepositoryImplementation) FindOrderItemsByIdOrder(db *gorm.DB, idOrder string) ([]entity.OrderItem, error) {
	orderItems := []entity.OrderItem{}
	result := db.
		Find(&orderItems, "id_order = ?", idOrder)
	return orderItems, result.Error
}
