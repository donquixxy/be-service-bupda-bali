package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/controller"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/repository"
	"github.com/tensuqiuwulu/be-service-bupda-bali/routes"
	"github.com/tensuqiuwulu/be-service-bupda-bali/service"
	"github.com/tensuqiuwulu/be-service-bupda-bali/utilities"
)

func main() {
	appConfig := config.GetConfig()

	// Database
	DBConn := repository.NewDatabaseConnection(&appConfig.Database)

	// Timezone
	location, err := time.LoadLocation(appConfig.Timezone.Timezone)
	time.Local = location
	fmt.Println("Location:", location, err)

	// Server App
	fmt.Println("Server App : ", string(appConfig.Application.Server))

	// Logger
	logrusLogger := utilities.NewLogger(appConfig.Log)

	// Validator
	validate := validator.New()

	e := echo.New()

	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      nil,
		ErrorMessage: "Request Timeout",
		Timeout:      15 * time.Second,
	}))
	e.Use(middleware.Recover())
	e.HTTPErrorHandler = exceptions.ErrorHandler
	e.Use(middleware.RequestID())

	// Repository
	userRepository := repository.NewUserRepository(&appConfig.Database)
	userProfileRepository := repository.NewUserProfileRepository(&appConfig.Database)
	otpManagerRepository := repository.NewOtpManagerRepository(&appConfig.Database)
	kecamatanRepository := repository.NewKecamatanRepository(&appConfig.Database)
	kelurahanRepository := repository.NewKelurahanRepository(&appConfig.Database)
	desaRepository := repository.NewDesaRepository(&appConfig.Database)
	productDesaRepository := repository.NewProductDesaRepository(&appConfig.Database)
	cartRepository := repository.NewCartRepository(&appConfig.Database)
	promoRepository := repository.NewPromoRepository(&appConfig.Database)
	pointRepository := repository.NewPointRepository(&appConfig.Database)
	orderRepository := repository.NewOrderRepository(&appConfig.Database)
	orderItemRepository := repository.NewOrderItemRepository(&appConfig.Database)
	orderItemPpobRepository := repository.NewOrderItemPpobRepository(&appConfig.Database)
	paymentChannelRepository := repository.NewPaymentChannelRepository(&appConfig.Database)
	productDesaStockRepository := repository.NewProductDesaStockHistoryRepository(&appConfig.Database)
	settingRepository := repository.NewSettingRepository(&appConfig.Database)
	userShippingAddressRepository := repository.NewUserShippingAddressRepository(&appConfig.Database)
	operatorPrefixRepository := repository.NewOperatorPrefixRepository(&appConfig.Database)
	ppobDetailRepository := repository.NewPpobDetailRepository(&appConfig.Database)
	infoDesaRepository := repository.NewInfoDesaRepository(&appConfig.Database)

	// Service
	infoDesaService := service.NewInfoDesaService(
		DBConn,
		validate,
		logrusLogger,
		infoDesaRepository,
	)
	authService := service.NewAuthService(
		DBConn,
		appConfig.Jwt,
		validate,
		logrusLogger,
		userRepository,
	)
	otpManagerService := service.NewOtpManagerService(
		DBConn,
		appConfig.Jwt,
		validate,
		logrusLogger,
		otpManagerRepository,
		userRepository,
	)
	kecamatanService := service.NewKecamatanService(
		DBConn,
		validate,
		logrusLogger,
		kecamatanRepository,
	)
	kelurahanService := service.NewKelurahanService(
		DBConn,
		validate,
		logrusLogger,
		kelurahanRepository,
	)
	desaService := service.NewDesaService(
		DBConn,
		validate,
		logrusLogger,
		desaRepository,
	)
	userService := service.NewUserService(
		DBConn,
		validate,
		appConfig.Jwt,
		logrusLogger,
		userRepository,
		userProfileRepository,
		pointRepository,
	)
	productDesaService := service.NewProductDesaService(
		DBConn,
		validate,
		logrusLogger,
		productDesaRepository,
		orderItemRepository,
		productDesaStockRepository,
	)
	cartService := service.NewCartService(
		DBConn,
		validate,
		logrusLogger,
		cartRepository,
		productDesaRepository,
	)
	promoService := service.NewPromoService(
		DBConn,
		validate,
		logrusLogger,
		promoRepository,
	)
	pointService := service.NewPointService(
		DBConn,
		validate,
		logrusLogger,
		pointRepository,
	)
	paymentService := service.NewPaymentService(
		logrusLogger,
	)
	orderService := service.NewOrderService(
		DBConn,
		validate,
		logrusLogger,
		orderRepository,
		userRepository,
		paymentService,
		cartRepository,
		orderItemRepository,
		paymentChannelRepository,
		productDesaRepository,
		productDesaService,
		operatorPrefixRepository,
		orderItemPpobRepository,
		ppobDetailRepository,
		desaRepository,
	)
	paymentChannelService := service.NewPaymentChannelService(
		DBConn,
		validate,
		logrusLogger,
		paymentChannelRepository,
	)
	settingService := service.NewSettingService(
		DBConn,
		validate,
		logrusLogger,
		settingRepository,
	)
	userShippingAddressService := service.NewUserShippingAddressService(
		DBConn,
		validate,
		logrusLogger,
		userShippingAddressRepository,
	)
	ppobService := service.NewPpobService(
		DBConn,
		validate,
		logrusLogger,
		operatorPrefixRepository,
		orderService,
	)

	infoDesaController := controller.NewInfoDesaController(
		infoDesaService,
	)
	// Controller
	authController := controller.NewAuthController(
		logrusLogger,
		authService,
	)
	otpManagerController := controller.NewOtpManagerController(
		logrusLogger,
		otpManagerService,
	)
	kecamatanController := controller.NewKecamatanController(
		kecamatanService,
	)
	kelurahanController := controller.NewKelurahanController(
		kelurahanService,
	)
	desaController := controller.NewDesaController(
		desaService,
	)
	userController := controller.NewUserController(
		logrusLogger,
		userService,
	)
	productDesaController := controller.NewProductDesaController(
		logrusLogger,
		productDesaService,
	)
	cartController := controller.NewCartController(
		logrusLogger,
		cartService,
	)
	promoController := controller.NewPromoController(
		promoService,
	)
	pointController := controller.NewPointController(
		pointService,
	)
	orderController := controller.NewOrderController(
		logrusLogger,
		orderService,
	)
	paymentChannelController := controller.NewPaymentChannelController(
		paymentChannelService,
	)
	settingController := controller.NewSettingController(
		settingService,
	)
	userShippingAddressController := controller.NewUserShippingAddressController(
		logrusLogger,
		userShippingAddressService,
	)
	ppobController := controller.NewPpobController(
		logrusLogger,
		ppobService,
	)

	// Route
	routes.InfoDesaRoute(e, appConfig.Jwt, infoDesaController)
	routes.AuthRoute(e, authController)
	routes.OtpManagerRoute(e, otpManagerController)
	routes.KecamatanRoute(e, kecamatanController)
	routes.KelurahanRoute(e, kelurahanController)
	routes.DesaRoute(e, desaController)
	routes.UserRoute(e, appConfig.Jwt, userController)
	routes.ProductDesaRoute(e, appConfig.Jwt, productDesaController)
	routes.CartRoute(e, appConfig.Jwt, cartController)
	routes.PromoRoute(e, appConfig.Jwt, promoController)
	routes.PointRoute(e, appConfig.Jwt, pointController)
	routes.OrderRoute(e, appConfig.Jwt, orderController)
	routes.PaymentChannelRoute(e, appConfig.Jwt, paymentChannelController)
	routes.SettingRoute(e, appConfig.Jwt, settingController)
	routes.UserShippingAddressRoute(e, appConfig.Jwt, userShippingAddressController)
	routes.PpobRoute(e, appConfig.Jwt, ppobController)

	// Careful shutdown
	go func() {
		if err := e.Start(":" + strconv.Itoa(int(appConfig.Webserver.Port))); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	fmt.Println("Running cleanup tasks...")

	// Your cleanup tasks go here
	// mysql database
	repository.Close(DBConn)
	fmt.Println("Echo was successful shutdown.")

}
