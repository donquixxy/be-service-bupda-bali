package repository

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type ProductDesaStockHistoryRepositoryInterface interface {
	CreateProductDesaStockHistory(db *gorm.DB, productDesaStockHistory *entity.ProductDesaStockHistory) error
}

type ProductDesaStockHistoryRepositoryImplementation struct {
	DB *config.Database
}

func NewProductDesaStockHistoryRepository(
	db *config.Database,
) ProductDesaStockHistoryRepositoryInterface {
	return &ProductDesaStockHistoryRepositoryImplementation{
		DB: db,
	}
}

func (repository *ProductDesaStockHistoryRepositoryImplementation) CreateProductDesaStockHistory(db *gorm.DB, productDesaStockHistory *entity.ProductDesaStockHistory) error {
	result := db.Create(&productDesaStockHistory)
	return result.Error
}
