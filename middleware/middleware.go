package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	modelService "github.com/tensuqiuwulu/be-service-bupda-bali/model/service"
)

func Authentication(configurationJWT config.Jwt) echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:       &modelService.TokenClaims{},
		SigningKey:   []byte(configurationJWT.Key),
		ErrorHandler: ErrorHandler,
	})
}

func RateLimit() echo.MiddlewareFunc {
	return middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: 10, Burst: 0, ExpiresIn: 6 * time.Hour},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, "Error")
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, "Error")
		},
	})
}

func Timeout() echo.MiddlewareFunc {
	return middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		ErrorMessage: "Request Timeout",
		Timeout:      15 * time.Second,
	})
}

func ErrorHandler(err error) error {
	fmt.Println(err)
	return err
}

func TokenClaimsIdUser(c echo.Context) (id string) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*modelService.TokenClaims)
	idUser := claims.Id
	return idUser
}

func TokenClaimsIdDesa(c echo.Context) (id string) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*modelService.TokenClaims)
	IdDesa := claims.IdDesa
	return IdDesa
}

func TokenClaimsAccountType(c echo.Context) (id int) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*modelService.TokenClaims)
	AccountType := claims.AccountType
	return AccountType
}
