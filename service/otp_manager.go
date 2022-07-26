package service

import (
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/request"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	modelService "github.com/tensuqiuwulu/be-service-bupda-bali/model/service"
	"github.com/tensuqiuwulu/be-service-bupda-bali/repository"
	"github.com/tensuqiuwulu/be-service-bupda-bali/utilities"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type OtpManagerServiceInterface interface {
	SendOtpBySms(requestId string, sendOtpBySmsRequest *request.SendOtpBySmsRequest)
	VerifyOtp(requestId string, verifyOtpRequest *request.VerifyOtpRequest) *response.VerifyOtpResponse
}

type OtpManagerServiceImplementation struct {
	DB                            *gorm.DB
	ConfigJwt                     config.Jwt
	Validate                      *validator.Validate
	Logger                        *logrus.Logger
	OtpManagerRepositoryInterface repository.OtpManagerRepositoryInterface
	AuthServiceInterface          AuthServiceInterface
	SmsServiceInterface           SmsServiceInterface
}

func NewOtpManagerService(
	db *gorm.DB,
	configJwt config.Jwt,
	validate *validator.Validate,
	logger *logrus.Logger,
	otpManaOtpManagerServiceInterface repository.OtpManagerRepositoryInterface,
) OtpManagerServiceInterface {
	return &OtpManagerServiceImplementation{
		DB:                            db,
		ConfigJwt:                     configJwt,
		Validate:                      validate,
		Logger:                        logger,
		OtpManagerRepositoryInterface: otpManaOtpManagerServiceInterface,
	}
}

func GenerateRandomOtpCode() string {
	rand.Seed(time.Now().Unix())
	charSet := "1234567890"
	var output strings.Builder
	length := 6

	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}

	return output.String()
}

func (service *OtpManagerServiceImplementation) SendOtpBySms(requestId string, sendOtpBySmsRequest *request.SendOtpBySmsRequest) {
	var err error

	request.ValidateRequest(service.Validate, sendOtpBySmsRequest, requestId, service.Logger)

	resultOtp, err := service.OtpManagerRepositoryInterface.FindOtpByPhone(service.DB, sendOtpBySmsRequest.Phone)
	exceptions.PanicIfError(err, requestId, service.Logger)

	// validasi data
	if len(resultOtp.Id) == 0 {
		otpManagerEntity := &entity.OtpManager{}
		otpManagerEntity.Id = utilities.RandomUUID()
		// otpManagerEntity.OtpCode = GenerateRandomOtpCode()
		otpManagerEntity.OtpCode = "123456"
		otpManagerEntity.Phone = sendOtpBySmsRequest.Phone
		otpManagerEntity.PhoneLimit = 5
		otpManagerEntity.IpAddressLimit = 5
		otpManagerEntity.OtpExperiedAt = time.Now().Add(time.Minute * 5)
		otpManagerEntity.CreatedDate = time.Now()

		createOtpErr := service.OtpManagerRepositoryInterface.CreateOtp(service.DB, otpManagerEntity)
		exceptions.PanicIfError(createOtpErr, requestId, service.Logger)

		// Send OTP
		// go service.SmsServiceInterface.SendOTP(requestId, otpManagerEntity.Phone, otpManagerEntity.OtpCode)
	} else {
		if resultOtp.PhoneLimit <= 0 {
			exceptions.PanicIfBadRequest(errors.New("phone daily limit"), requestId, []string{"phone daily limit"}, service.Logger)
		}
		otpManagerEntity := &entity.OtpManager{}
		// otpManagerEntity.OtpCode = GenerateRandomOtpCode()
		otpManagerEntity.OtpCode = "123456"
		otpManagerEntity.PhoneLimit = resultOtp.PhoneLimit - 1
		otpManagerEntity.OtpExperiedAt = time.Now().Add(time.Minute * 5)
		otpManagerEntity.UpdatedDate = null.NewTime(time.Now(), true)

		// Update OTP
		updateOtpErr := service.OtpManagerRepositoryInterface.UpdateOtp(service.DB, resultOtp.Id, otpManagerEntity)
		exceptions.PanicIfError(updateOtpErr, requestId, service.Logger)
		// Send OTP

		// go service.SmsServiceInterface.SendOTP(requestId, otpManagerEntity.Phone, otpManagerEntity.OtpCode)
	}
}

func (service *OtpManagerServiceImplementation) VerifyOtp(requestId string, verifyOtpRequest *request.VerifyOtpRequest) *response.VerifyOtpResponse {
	request.ValidateRequest(service.Validate, verifyOtpRequest, requestId, service.Logger)

	otp, _ := service.OtpManagerRepositoryInterface.FindOtpByPhone(service.DB, verifyOtpRequest.Phone)
	if len(otp.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("otp not found"), requestId, []string{"otp not found"}, service.Logger)
	}

	if verifyOtpRequest.OtpCode != otp.OtpCode {
		exceptions.PanicIfBadRequest(errors.New("otp not match"), requestId, []string{"otp not match"}, service.Logger)
	}

	if time.Now().After(otp.OtpExperiedAt) {
		exceptions.PanicIfBadRequest(errors.New("otp code has expired"), requestId, []string{"otp code has expired"}, service.Logger)
	}

	token, _ := service.GenerateFormToken()

	verifyOtpResponse := response.ToVerifyOtpResponse(token)

	return &verifyOtpResponse

}

func (service *OtpManagerServiceImplementation) GenerateFormToken() (token string, err error) {
	// Create the Claims
	claims := modelService.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(service.ConfigJwt.Tokenexpiredtime)).Unix(),
			Issuer:    "cyrilia",
		},
	}

	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenWithClaims.SignedString([]byte(service.ConfigJwt.FormKey))
	if err != nil {
		return "", err
	}
	return token, err
}
