package repository

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type CartRepositoryInterface interface {
	CreateCart(db *gorm.DB, cart *entity.Cart) error
	UpdateCart(db *gorm.DB, idCart string, cart *entity.Cart) error
	FindCartByUser(db *gorm.DB, idUser string) ([]entity.Cart, error)
	FindCartById(db *gorm.DB, idCart string) (*entity.Cart, error)
	FindCartByProductDesa(db *gorm.DB, idUser, idProductDesa string) (*entity.Cart, error)
	DeleteCartById(db *gorm.DB, idCart string) error
	DeleteCartByUser(db *gorm.DB, idUser string, cart []entity.Cart) error
}

type CartRepositoryImplementation struct {
	DB *config.Database
}

func NewCartRepository(
	db *config.Database,
) CartRepositoryInterface {
	return &CartRepositoryImplementation{
		DB: db,
	}
}

func (repository *CartRepositoryImplementation) CreateCart(db *gorm.DB, cart *entity.Cart) error {
	result := db.Create(&cart)
	return result.Error
}

func (repository *CartRepositoryImplementation) UpdateCart(db *gorm.DB, idCart string, cart *entity.Cart) error {
	result := db.
		Model(entity.Cart{}).
		Where("id = ?", idCart).
		Updates(entity.Cart{
			Qty: cart.Qty,
		})
	return result.Error
}

func (repository *CartRepositoryImplementation) FindCartByUser(db *gorm.DB, idUser string) ([]entity.Cart, error) {
	cart := []entity.Cart{}
	results := db.
		Where("cart.id_user = ?", idUser).
		Joins("ProductsDesa").
		Preload("ProductsDesa.ProductsMaster").
		Find(&cart)
	return cart, results.Error
}

func (repository *CartRepositoryImplementation) FindCartById(db *gorm.DB, idCart string) (*entity.Cart, error) {
	cart := &entity.Cart{}
	result := db.Find(cart, "id = ?", idCart)
	return cart, result.Error
}

func (repository *CartRepositoryImplementation) FindCartByProductDesa(db *gorm.DB, idUser, idProductDesa string) (*entity.Cart, error) {
	cart := &entity.Cart{}
	result := db.Find(cart, "id_user = ? AND id_product_desa = ?", idUser, idProductDesa)
	return cart, result.Error
}

func (repository *CartRepositoryImplementation) DeleteCartById(db *gorm.DB, idCart string) error {
	result := db.Where("id = ?", idCart).Delete(&entity.Cart{})
	return result.Error
}

func (repository *CartRepositoryImplementation) DeleteCartByUser(DB *gorm.DB, idUser string, carts []entity.Cart) error {
	result := DB.Where("id_user = ?", idUser).Delete(carts)
	return result.Error
}
