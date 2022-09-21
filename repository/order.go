package repository

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type OrderRepositoryInterface interface {
	CreateOrder(db *gorm.DB, order *entity.Order) error
	FindOrderByNumberOrder(db *gorm.DB, numberOrder string) (*entity.Order, error)
	FindOrderByUser(db *gorm.DB, idUser string, orderStatus int) ([]entity.Order, error)
	FindOrderById(db *gorm.DB, idOrder string) (*entity.Order, error)
	FindOrderByRefId(db *gorm.DB, refId string) (*entity.Order, error)
	UpdateOrderByIdOrder(db *gorm.DB, idOrder string, orderUpdate *entity.Order) error
	FindOrderPrepaidPulsaById(db *gorm.DB, idUser string, productType string) (*entity.Order, error)
	FindOrderPrepaidPlnById(db *gorm.DB, idUser string) (*entity.Order, error)
}

type OrderRepositoryImplementation struct {
	DB *config.Database
}

func NewOrderRepository(
	db *config.Database,
) OrderRepositoryInterface {
	return &OrderRepositoryImplementation{
		DB: db,
	}
}

func (repository *OrderRepositoryImplementation) CreateOrder(db *gorm.DB, order *entity.Order) error {
	result := db.Create(&order)
	return result.Error
}

func (repository *OrderRepositoryImplementation) FindOrderByNumberOrder(db *gorm.DB, numberOrder string) (*entity.Order, error) {
	order := &entity.Order{}
	result := db.Find(order, "number_order = ?", numberOrder)
	return order, result.Error
}

func (repository *OrderRepositoryImplementation) FindOrderByUser(db *gorm.DB, idUser string, orderStatus int) ([]entity.Order, error) {
	var result *gorm.DB
	order := []entity.Order{}
	if orderStatus >= 0 {
		result = db.Order("order_date desc").Find(&order, "id_user = ? AND order_status = ?", idUser, orderStatus)
	} else {
		result = db.Order("order_date desc").Find(&order, "id_user = ?", idUser)
	}

	return order, result.Error
}

func (repository *OrderRepositoryImplementation) FindOrderById(db *gorm.DB, idUser string) (*entity.Order, error) {
	orders := &entity.Order{}
	result := db.Find(orders, "id = ?", idUser)
	return orders, result.Error
}

func (repository *OrderRepositoryImplementation) FindOrderByRefId(db *gorm.DB, refId string) (*entity.Order, error) {
	orders := &entity.Order{}
	result := db.Find(orders, "ref_id = ?", refId)
	return orders, result.Error
}

func (repository *OrderRepositoryImplementation) FindOrderPrepaidPulsaById(db *gorm.DB, idUser string, productType string) (*entity.Order, error) {
	orders := &entity.Order{}
	result := db.Where("product_type = ?", productType).Find(orders, "id = ?", idUser)
	return orders, result.Error
}

func (repository *OrderRepositoryImplementation) FindOrderPrepaidPlnById(db *gorm.DB, idUser string) (*entity.Order, error) {
	orders := &entity.Order{}
	result := db.Find(orders, "id = ?", idUser)
	return orders, result.Error
}

func (repository *OrderRepositoryImplementation) UpdateOrderByIdOrder(db *gorm.DB, idOrder string, orderUpdate *entity.Order) error {
	order := &entity.Order{}
	result := db.
		Model(order).
		Where("id = ?", idOrder).
		Updates(orderUpdate)
	return result.Error
}
