package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/middleware"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/service"
)

type ProductDesaControllerInterface interface {
	FindProductsDesa(c echo.Context) error
	FindProductsDesaByCategory(c echo.Context) error
	FindProductsDesaBySubCategory(c echo.Context) error
	FindProductsDesaNotoken(c echo.Context) error
	FindProductsDesaByCategoryNotoken(c echo.Context) error
	FindProductsDesaBySubCategoryNotoken(c echo.Context) error
	FindProductDesaById(c echo.Context) error
	FindProductsDesaByPromo(c echo.Context) error
}

type ProductDesaControllerImplementation struct {
	Logger                      *logrus.Logger
	ProductDesaServiceInterface service.ProductDesaServiceInterface
}

func NewProductDesaController(
	logger *logrus.Logger,
	productDesaServiceInterface service.ProductDesaServiceInterface,
) ProductDesaControllerInterface {
	return &ProductDesaControllerImplementation{
		Logger:                      logger,
		ProductDesaServiceInterface: productDesaServiceInterface,
	}
}

func (controller *ProductDesaControllerImplementation) FindProductsDesa(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idDesa := middleware.TokenClaimsIdDesa(c)
	accountType := middleware.TokenClaimsAccountType(c)
	productsDesaResponse := controller.ProductDesaServiceInterface.FindProductsDesa(requestId, idDesa, accountType)
	responses := response.Response{Code: 200, Mssg: "success", Data: productsDesaResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *ProductDesaControllerImplementation) FindProductsDesaNotoken(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idDesa := "44dcd5ce-aa07-43da-ba04-2bff14a4d7ac"
	accountType := 1
	productsDesaResponse := controller.ProductDesaServiceInterface.FindProductsDesa(requestId, idDesa, accountType)
	responses := response.Response{Code: 200, Mssg: "success", Data: productsDesaResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *ProductDesaControllerImplementation) FindProductDesaById(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idProductDesa := c.QueryParam("id_product")
	accountType := middleware.TokenClaimsAccountType(c)
	productsDesaResponse := controller.ProductDesaServiceInterface.FindProductDesaById(requestId, idProductDesa, accountType)
	responses := response.Response{Code: 200, Mssg: "success", Data: productsDesaResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *ProductDesaControllerImplementation) FindProductsDesaByCategory(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idDesa := middleware.TokenClaimsIdDesa(c)
	idCategory, _ := strconv.Atoi(c.QueryParam("id_category"))
	accountType := middleware.TokenClaimsAccountType(c)
	productsDesaResponse := controller.ProductDesaServiceInterface.FindProductsDesaByCategory(requestId, idDesa, idCategory, accountType)
	responses := response.Response{Code: 200, Mssg: "success", Data: productsDesaResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *ProductDesaControllerImplementation) FindProductsDesaByCategoryNotoken(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idDesa := "44dcd5ce-aa07-43da-ba04-2bff14a4d7ac"
	idCategory, _ := strconv.Atoi(c.QueryParam("id_category"))
	accountType := 1
	productsDesaResponse := controller.ProductDesaServiceInterface.FindProductsDesaByCategory(requestId, idDesa, idCategory, accountType)
	responses := response.Response{Code: 200, Mssg: "success", Data: productsDesaResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *ProductDesaControllerImplementation) FindProductsDesaBySubCategory(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idDesa := middleware.TokenClaimsIdDesa(c)
	idSubCategory, _ := strconv.Atoi(c.QueryParam("id_sub_category"))
	accountType := middleware.TokenClaimsAccountType(c)
	productsDesaResponse := controller.ProductDesaServiceInterface.FindProductsDesaBySubCategory(requestId, idDesa, idSubCategory, accountType)
	responses := response.Response{Code: 200, Mssg: "success", Data: productsDesaResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *ProductDesaControllerImplementation) FindProductsDesaBySubCategoryNotoken(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idDesa := "44dcd5ce-aa07-43da-ba04-2bff14a4d7ac"
	idSubCategory, _ := strconv.Atoi(c.QueryParam("id_sub_category"))
	accountType := 1
	productsDesaResponse := controller.ProductDesaServiceInterface.FindProductsDesaBySubCategory(requestId, idDesa, idSubCategory, accountType)
	responses := response.Response{Code: 200, Mssg: "success", Data: productsDesaResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *ProductDesaControllerImplementation) FindProductsDesaByPromo(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idDesa := middleware.TokenClaimsIdDesa(c)
	IdPromo := c.QueryParam("id_promo")
	accountType := middleware.TokenClaimsAccountType(c)
	productsDesaResponse := controller.ProductDesaServiceInterface.FindProductsDesaByPromo(requestId, idDesa, IdPromo, accountType)
	responses := response.Response{Code: 200, Mssg: "success", Data: productsDesaResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}
