package entity

type ProductsDesa struct {
	Id              string         `gorm:"primaryKey;column:id;"`
	IdProduct       string         `gorm:"column:id_product;"`
	IdPromo         string         `gorm:"column:id_promo;"`
	IdType          int            `gorm:"column:id_type;"`
	ProductsMaster  ProductsMaster `gorm:"foreignKey:IdProduct;"`
	Price           float64        `gorm:"column:price;"`
	PricePromo      float64        `gorm:"column:price_promo;"`
	PriceGrosir     float64        `gorm:"column:price_grosir;"`
	PercentagePromo float64        `gorm:"column:percentage_promo;"`
	IdDesa          string         `gorm:"column:id_desa;"`
	StockOpname     int            `gorm:"column:stock_opname;"`
	IsPromo         int            `gorm:"column:is_promo;"`
}

func (ProductsDesa) TableName() string {
	return "products_desa"
}

type ProductsPackageItems struct {
	Id                  string  `gorm:"primaryKey;column:id;"`
	IdProductPackgeDesa string  `gorm:"column:id_product_primary;"`
	IdProductItem       string  `gorm:"column:id_product_item;"`
	NoSku               string  `gorm:"column:no_sku;"`
	ProductName         string  `gorm:"column:product_name;"`
	Price               float64 `gorm:"column:price;"`
	Qty                 int     `gorm:"column:qty;"`
}

func (ProductsPackageItems) TableName() string {
	return "products_package_items"
}
