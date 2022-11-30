package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
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
	UserRepositoryInterface       repository.UserRepositoryInterface
}

func NewOtpManagerService(
	db *gorm.DB,
	configJwt config.Jwt,
	validate *validator.Validate,
	logger *logrus.Logger,
	otpManaOtpManagerServiceInterface repository.OtpManagerRepositoryInterface,
	userRepositoryInterface repository.UserRepositoryInterface,
) OtpManagerServiceInterface {
	return &OtpManagerServiceImplementation{
		DB:                            db,
		ConfigJwt:                     configJwt,
		Validate:                      validate,
		Logger:                        logger,
		OtpManagerRepositoryInterface: otpManaOtpManagerServiceInterface,
		UserRepositoryInterface:       userRepositoryInterface,
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

	user, _ := service.UserRepositoryInterface.FindUserByPhone(service.DB, sendOtpBySmsRequest.Phone)

	if user.StatusPaylater == -1 {
		exceptions.PanicIfUserNotHavePassword(errors.New("user dont have password"), requestId, []string{"User dont have password"}, service.Logger)
	}

	resultOtp, err := service.OtpManagerRepositoryInterface.FindOtpByPhone(service.DB, sendOtpBySmsRequest.Phone)
	exceptions.PanicIfError(err, requestId, service.Logger)

	// Find user by phone
	if sendOtpBySmsRequest.TypeOtp == 1 {
		user, _ := service.UserRepositoryInterface.FindUserByPhone(service.DB, sendOtpBySmsRequest.Phone)
		if len(user.Id) != 0 {
			exceptions.PanicIfBadRequest(errors.New("phone already user"), requestId, []string{"phone already user"}, service.Logger)
		}
	}

	// validasi data
	if len(resultOtp.Id) == 0 {
		otpManagerEntity := &entity.OtpManager{}
		otpManagerEntity.Id = utilities.RandomUUID()
		otpManagerEntity.OtpCode = GenerateRandomOtpCode()
		// otpManagerEntity.OtpCode = "123456"
		otpManagerEntity.Phone = sendOtpBySmsRequest.Phone
		otpManagerEntity.PhoneLimit = 5
		otpManagerEntity.IpAddressLimit = 5
		otpManagerEntity.OtpExperiedAt = time.Now().Add(time.Minute * 5)
		otpManagerEntity.CreatedDate = time.Now()

		// Send OTP
		go SendOTP(requestId, sendOtpBySmsRequest.Phone, otpManagerEntity.OtpCode)

		createOtpErr := service.OtpManagerRepositoryInterface.CreateOtp(service.DB, otpManagerEntity)
		exceptions.PanicIfError(createOtpErr, requestId, service.Logger)

	} else {
		// log.Println("phone limit = ", resultOtp.PhoneLimit)

		if resultOtp.PhoneLimit <= 0 {
			if resultOtp.FreezeDueDate.Time.Before(time.Now()) {
				otpManagerEntity := &entity.OtpManager{}
				otpManagerEntity.Id = resultOtp.Id
				otpManagerEntity.OtpCode = GenerateRandomOtpCode()
				// otpManagerEntity.OtpCode = "123456"
				otpManagerEntity.Phone = resultOtp.Phone
				otpManagerEntity.PhoneLimit = 5
				otpManagerEntity.IpAddressLimit = 5
				otpManagerEntity.OtpExperiedAt = time.Now().Add(time.Minute * 5)
				otpManagerEntity.CreatedDate = time.Now()

				// Send OTP
				go SendOTP(requestId, resultOtp.Phone, otpManagerEntity.OtpCode)

				updateOtpErr := service.OtpManagerRepositoryInterface.UpdateOtp(service.DB, resultOtp.Id, otpManagerEntity)
				exceptions.PanicIfError(updateOtpErr, requestId, service.Logger)

			} else {
				exceptions.PanicIfBadRequest(errors.New("phone daily limit"), requestId, []string{"phone daily limit"}, service.Logger)
			}
		} else {
			otpManagerEntity := &entity.OtpManager{}
			otpManagerEntity.Id = resultOtp.Id
			otpManagerEntity.OtpCode = GenerateRandomOtpCode()
			// otpManagerEntity.OtpCode = "123456"
			otpManagerEntity.Phone = resultOtp.Phone
			if resultOtp.PhoneLimit <= 1 {
				otpManagerEntity.FreezeDueDate = null.NewTime(time.Now().Add(time.Hour*24), true)
			}
			otpManagerEntity.PhoneLimit = resultOtp.PhoneLimit - 1
			otpManagerEntity.IpAddressLimit = 5
			otpManagerEntity.OtpExperiedAt = time.Now().Add(time.Minute * 5)
			otpManagerEntity.CreatedDate = time.Now()

			// Send OTP
			go SendOTP(requestId, resultOtp.Phone, otpManagerEntity.OtpCode)

			updateOtpErr := service.OtpManagerRepositoryInterface.UpdateOtp(service.DB, resultOtp.Id, otpManagerEntity)
			exceptions.PanicIfError(updateOtpErr, requestId, service.Logger)
		}
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
	var userModelService modelService.User
	userModelService.Phone = otp.Phone
	token, _ := service.GenerateFormToken(userModelService)
	verifyOtpResponse := response.ToVerifyOtpResponse(token)

	return &verifyOtpResponse

}

func (service *OtpManagerServiceImplementation) GenerateFormToken(user modelService.User) (token string, err error) {
	// Create the Claims
	claims := modelService.TokenClaims{
		Phone: user.Phone,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(10)).Unix(),
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

func SendOTP(requestId string, phone string, otpCode string) {
	message := fmt.Sprintf("Kode Verifikasi Akun Bupda Bali Anda adalah: %s *JANGAN BERIKAN KODE INI KEPADA SIAPAPUN, TERMASUK PIHAK BUPDA BALI* Hubungi 085960144218 untuk bantuan.", otpCode)

	postBody, _ := json.Marshal(map[string]string{
		"userkey": config.GetConfig().Sms.UserKey,
		"passkey": config.GetConfig().Sms.PassKey,
		"to":      phone,
		"message": message,
	})

	responseBody := bytes.NewBuffer(postBody)
	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post("https://console.zenziva.net/masking/api/sendOTP", "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Printf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	sb := string(body)
	log.Println(sb)
}
