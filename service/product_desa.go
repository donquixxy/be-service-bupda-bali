package service

import (
	"errors"
	"time"

	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/repository"
	"gorm.io/gorm"
)

type ProductDesaServiceInterface interface {
	FindProductsDesa(requestId string, IdDesa string, AccountType int) (productsDesaResponses []response.FindProductsDesaResponse)
	FindProductsDesaByCategory(requestId string, IdDesa string, IdCaetgory, AccountType int) (productsDesaResponses []response.FindProductsDesaResponse)
	FindProductsDesaBySubCategory(requestId string, IdDesa string, IdSubCaetgory, AccountType int) (productsDesaResponses []response.FindProductsDesaResponse)
	FindProductsDesaByPromo(requestId string, IdDesa string, IdPromo string, AccountType int) (productsDesaResponses []response.FindProductsDesaResponse)
	FindProductDesaById(requestId string, IdProductDesa string, AccountType int) (productDesaReponse response.FindProductDesaByIdResponse)
	UpdateProductStock(requestId string, IdOrder string, db *gorm.DB)
}

type ProductDesaServiceImplementation struct {
	DB                               *gorm.DB
	Validate                         *validator.Validate
	Logger                           *logrus.Logger
	ProductDesaRepositoryInterface   repository.ProductDesaRepositoryInterface
	OrderItemRepositoryInterface     repository.OrderItemRepositoryInterface
	ProductDesaStockHistoryInterface repository.ProductDesaStockHistoryRepositoryInterface
}

func NewProductDesaService(
	db *gorm.DB,
	validate *validator.Validate,
	logger *logrus.Logger,
	productDesaRepositoryInterface repository.ProductDesaRepositoryInterface,
	orderItemRepositoryInterface repository.OrderItemRepositoryInterface,
	productDesaStockHistoryInterface repository.ProductDesaStockHistoryRepositoryInterface,
) ProductDesaServiceInterface {
	return &ProductDesaServiceImplementation{
		DB:                               db,
		Validate:                         validate,
		Logger:                           logger,
		ProductDesaRepositoryInterface:   productDesaRepositoryInterface,
		OrderItemRepositoryInterface:     orderItemRepositoryInterface,
		ProductDesaStockHistoryInterface: productDesaStockHistoryInterface,
	}
}

func (service *ProductDesaServiceImplementation) FindProductsDesa(requestid string, IdDesa string, AccountType int) (productsDesaResponses []response.FindProductsDesaResponse) {
	productsDesa, err := service.ProductDesaRepositoryInterface.FindProductsDesa(service.DB, IdDesa)
	exceptions.PanicIfError(err, requestid, service.Logger)
	if len(productsDesa) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("not found"), requestid, []string{"data not found"}, service.Logger)
	}
	productsDesaResponses = response.ToFindProductsDesaResponse(productsDesa, AccountType)
	return productsDesaResponses
}

func (service *ProductDesaServiceImplementation) FindProductsDesaByCategory(requestid string, IdDesa string, IdCaetgory, AccountType int) (productsDesaResponses []response.FindProductsDesaResponse) {
	productsDesa, err := service.ProductDesaRepositoryInterface.FindProductsDesaByCategory(service.DB, IdDesa, IdCaetgory)
	exceptions.PanicIfError(err, requestid, service.Logger)
	if len(productsDesa) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("not found"), requestid, []string{"data not found"}, service.Logger)
	}
	productsDesaResponses = response.ToFindProductsDesaResponse(productsDesa, AccountType)
	return productsDesaResponses
}

func (service *ProductDesaServiceImplementation) FindProductsDesaBySubCategory(requestid string, IdDesa string, IdSubCaetgory, AccountType int) (productsDesaResponses []response.FindProductsDesaResponse) {
	productsDesa, err := service.ProductDesaRepositoryInterface.FindProductsDesaBySubCategory(service.DB, IdDesa, IdSubCaetgory)
	exceptions.PanicIfError(err, requestid, service.Logger)
	if len(productsDesa) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("not found"), requestid, []string{"data not found"}, service.Logger)
	}
	productsDesaResponses = response.ToFindProductsDesaResponse(productsDesa, AccountType)
	return productsDesaResponses
}

func (service *ProductDesaServiceImplementation) FindProductsDesaByPromo(requestid string, IdDesa string, IdPromo string, AccountType int) (productsDesaResponses []response.FindProductsDesaResponse) {
	productsDesa, err := service.ProductDesaRepositoryInterface.FindProductsDesaByPromo(service.DB, IdDesa, IdPromo)
	// fmt.Println("product = ", productsDesa)
	exceptions.PanicIfError(err, requestid, service.Logger)
	if len(productsDesa) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("not found"), requestid, []string{"data not found"}, service.Logger)
	}
	productsDesaResponses = response.ToFindProductsDesaResponse(productsDesa, AccountType)
	return productsDesaResponses
}

func (service *ProductDesaServiceImplementation) FindProductDesaById(requestid string, IdProductDesaDesa string, AccountType int) (productDesaResponse response.FindProductDesaByIdResponse) {
	var errr error
	productDesa, err := service.ProductDesaRepositoryInterface.FindProductDesaById(service.DB, IdProductDesaDesa)
	exceptions.PanicIfError(err, requestid, service.Logger)
	if len(productDesa.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("not found"), requestid, []string{"data not found"}, service.Logger)
	}

	var ListItemsPackage []entity.ProductsPackageItems

	if productDesa.IdType == 2 {
		ListItemsPackage, errr = service.ProductDesaRepositoryInterface.FindListPackageByIdProductDesa(service.DB, IdProductDesaDesa)
		exceptions.PanicIfError(errr, requestid, service.Logger)
	}

	productDesaResponse = response.ToFindProductDesaByIdResponse(productDesa, AccountType, ListItemsPackage)
	return productDesaResponse
}

func (service *ProductDesaServiceImplementation) UpdateProductStock(requestId string, IdOrder string, db *gorm.DB) {
	var err error
	orderItems, err := service.OrderItemRepositoryInterface.FindOrderItemsByIdOrder(db, IdOrder)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(orderItems) == 0 {
		exceptions.PanicIfErrorWithRollback(errors.New("product not found"), requestId, []string{"product not found"}, service.Logger, db)
	}

	for _, orderItem := range orderItems {
		productDesa, errFindProduct := service.ProductDesaRepositoryInterface.FindProductDesaById(db, orderItem.IdProductDesa)
		exceptions.PanicIfErrorWithRollback(errFindProduct, requestId, []string{"product not found"}, service.Logger, db)

		productDesaEntityStockHistory := &entity.ProductDesaStockHistory{}
		productDesaEntityStockHistory.IdProductDesa = orderItem.IdProductDesa
		productDesaEntityStockHistory.TransDate = time.Now()
		productDesaEntityStockHistory.MinStockQty = orderItem.Qty
		productDesaEntityStockHistory.StockFinal = productDesa.StockOpname - orderItem.Qty
		productDesaEntityStockHistory.Description = "PEMBELIAN " + orderItem.NumberOrder
		productDesaEntityStockHistory.CreatedDate = time.Now()
		productDesaEntityStockHistory.HargaBeli = 0
		productDesaEntityStockHistory.HargaJual = orderItem.Price
		err = service.ProductDesaStockHistoryInterface.CreateProductDesaStockHistory(db, productDesaEntityStockHistory)
		exceptions.PanicIfErrorWithRollback(err, requestId, []string{"add stock history error"}, service.Logger, db)

		productDesaEntity := &entity.ProductsDesa{}
		productDesaEntity.StockOpname = productDesa.StockOpname - orderItem.Qty
		err = service.ProductDesaRepositoryInterface.UpdateProductStock(db, productDesa.Id, productDesaEntity)
		exceptions.PanicIfErrorWithRollback(err, requestId, []string{"update stock error"}, service.Logger, db)
	}
}
