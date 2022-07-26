package entity

import (
	"time"
)

type ProductDesaStockHistory struct {
	Id            string       `gorm:"primaryKey;column:id;"`
	IdProductDesa string       `gorm:"column:id_product_desa;"`
	ProductsDesa  ProductsDesa `gorm:"foreignKey:IdProductDesa;"`
	TransDate     time.Time    `gorm:"column:trans_date;"`
	AddStockQty   int          `gorm:"column:add_stock_qty;"`
	MinStockQty   int          `gorm:"column:min_stock_qty;"`
	StockOpname   int          `gorm:"column:stock_opname;"`
	StockFinal    int          `gorm:"column:stock_final;"`
	HargaBeli     float64      `gorm:"column:harga_beli;"`
	HargaJual     float64      `gorm:"column:harga_jual;"`
	Description   string       `gorm:"column:description;"`
	CreatedDate   time.Time    `gorm:"column:created_at;"`
}

func (ProductDesaStockHistory) TableName() string {
	return "products_desa_stock_history"
}
