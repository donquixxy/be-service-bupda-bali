package repository

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type ProductDesaRepositoryInterface interface {
	FindProductsDesa(db *gorm.DB, IdDesa string) ([]entity.ProductsDesa, error)
	FindProductsDesaByCategory(db *gorm.DB, IdDesa string, IdCategory int) ([]entity.ProductsDesa, error)
	FindProductsDesaBySubCategory(db *gorm.DB, IdDesa string, IdSubCategory int) ([]entity.ProductsDesa, error)
	FindProductsDesaByPromo(db *gorm.DB, IdDesa string, IdPromo string) ([]entity.ProductsDesa, error)
	FindProductDesaById(db *gorm.DB, IdProductDesa string) (*entity.ProductsDesa, error)
	UpdateProductStock(db *gorm.DB, idProductDesa string, productDesa *entity.ProductsDesa) error
}

type ProductDesaRepositoryImplementation struct {
	DB *config.Database
}

func NewProductDesaRepository(
	db *config.Database,
) ProductDesaRepositoryInterface {
	return &ProductDesaRepositoryImplementation{
		DB: db,
	}
}

func (repository *ProductDesaRepositoryImplementation) FindProductsDesa(db *gorm.DB, IdDesa string) ([]entity.ProductsDesa, error) {
	productsDesa := []entity.ProductsDesa{}
	result := db.
		Joins("ProductsMaster").
		Find(&productsDesa, "products_desa.id_desa = ?", IdDesa)
	return productsDesa, result.Error
}

func (repository *ProductDesaRepositoryImplementation) FindProductsDesaByCategory(db *gorm.DB, IdDesa string, IdCategory int) ([]entity.ProductsDesa, error) {
	productsDesa := []entity.ProductsDesa{}
	result := db.
		Joins("ProductsMaster").
		Where("ProductsMaster.id_category = ?", IdCategory).
		Find(&productsDesa, "products_desa.id_desa = ?", IdDesa)
	return productsDesa, result.Error
}

func (repository *ProductDesaRepositoryImplementation) FindProductsDesaBySubCategory(db *gorm.DB, IdDesa string, IdSubCategory int) ([]entity.ProductsDesa, error) {
	productsDesa := []entity.ProductsDesa{}
	result := db.
		Joins("ProductsMaster").
		Where("ProductsMaster.id_sub_category = ?", IdSubCategory).
		Find(&productsDesa, "products_desa.id_desa = ?", IdDesa)
	return productsDesa, result.Error
}

func (repository *ProductDesaRepositoryImplementation) FindProductsDesaByPromo(db *gorm.DB, IdDesa string, IdPromo string) ([]entity.ProductsDesa, error) {
	productsDesa := []entity.ProductsDesa{}
	result := db.
		Joins("ProductsMaster").
		Find(&productsDesa, "products_desa.id_desa = ? AND products_desa.id_promo = ? AND products_desa.is_promo = ?", IdDesa, IdPromo, 1)
	return productsDesa, result.Error
}

func (repository *ProductDesaRepositoryImplementation) FindProductDesaById(db *gorm.DB, IdProductDesa string) (*entity.ProductsDesa, error) {
	productsDesa := &entity.ProductsDesa{}
	result := db.
		Joins("ProductsMaster").
		Find(productsDesa, "products_desa.id = ?", IdProductDesa)
	return productsDesa, result.Error
}

func (repository *ProductDesaRepositoryImplementation) UpdateProductStock(db *gorm.DB, idProductDesa string, productDesa *entity.ProductsDesa) error {
	updateProduct := make(map[string]interface{})
	updateProduct["stock_opname"] = productDesa.StockOpname
	result := db.
		Model(entity.ProductsDesa{}).
		Where("id = ?", idProductDesa).
		Updates(&updateProduct)
	return result.Error
}
