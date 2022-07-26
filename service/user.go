package service

import (
	"errors"
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
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServiceInterface interface {
	CreateUserNonSuveyed(requestId string, createUserRequest *request.CreateUserRequest)
	FindUserById(requestId string, idUser string) (userResponse response.FindUserIdResponse)
}

type UserServiceImplementation struct {
	DB                             *gorm.DB
	Validate                       *validator.Validate
	ConfigJwt                      config.Jwt
	Logger                         *logrus.Logger
	UserRepositoryInterface        repository.UserRepositoryInterface
	UserProfileRepositoryInterface repository.UserProfileRepositoryInterface
	PointRepositoryInterface       repository.PointRepositoryInterface
}

func NewUserService(
	db *gorm.DB,
	validate *validator.Validate,
	configJwt config.Jwt,
	logger *logrus.Logger,
	userRepositoryInterface repository.UserRepositoryInterface,
	userProfileRepositoryInterface repository.UserProfileRepositoryInterface,
	pointRepositoryInterface repository.PointRepositoryInterface,
) UserServiceInterface {
	return &UserServiceImplementation{
		DB:                             db,
		Validate:                       validate,
		ConfigJwt:                      configJwt,
		Logger:                         logger,
		UserRepositoryInterface:        userRepositoryInterface,
		UserProfileRepositoryInterface: userProfileRepositoryInterface,
		PointRepositoryInterface:       pointRepositoryInterface,
	}
}

func (service *UserServiceImplementation) VerifyFormToken(requestId, token string) {
	tokenParse, err := jwt.ParseWithClaims(token, &modelService.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(service.ConfigJwt.FormKey), nil
	})

	if !tokenParse.Valid {
		exceptions.PanicIfUnauthorized(err, requestId, []string{"invalid token"}, service.Logger)
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			exceptions.PanicIfUnauthorized(err, requestId, []string{"invalid token"}, service.Logger)
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			exceptions.PanicIfUnauthorized(err, requestId, []string{"invalid token"}, service.Logger)
		} else {
			exceptions.PanicIfUnauthorized(err, requestId, []string{"invalid token"}, service.Logger)
		}
	}
}

func (service *UserServiceImplementation) CreateUserNonSuveyed(requestId string, createUserRequest *request.CreateUserRequest) {
	var err error

	request.ValidateRequest(service.Validate, createUserRequest, requestId, service.Logger)

	// validate token
	service.VerifyFormToken(requestId, createUserRequest.FormToken)

	// Check No Identitas
	NoIdentitasCheck, err := service.UserProfileRepositoryInterface.FindUserByNoIdentitas(service.DB, createUserRequest.NoIdentitas)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(NoIdentitasCheck.Id) != 0 {
		exceptions.PanicIfRecordAlreadyExists(errors.New("no identitas already exist"), requestId, []string{"no identitas sudah digunakan"}, service.Logger)
	}

	// Check email if exsict
	emailLowerCase := strings.ToLower(createUserRequest.Email)
	emailChek, err := service.UserProfileRepositoryInterface.FindUserByEmail(service.DB, emailLowerCase)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(emailChek.Id) != 0 {
		exceptions.PanicIfRecordAlreadyExists(errors.New("email already exist"), requestId, []string{"Email sudah digunakan"}, service.Logger)
	}

	// Check No Hp
	phoneCheck, err := service.UserRepositoryInterface.FindUserByPhone(service.DB, createUserRequest.Phone)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(phoneCheck.Id) != 0 {
		exceptions.PanicIfRecordAlreadyExists(errors.New("phone already exist"), requestId, []string{"phone sudah digunakan"}, service.Logger)
	}

	// Hash password
	password := strings.ReplaceAll(createUserRequest.Password, " ", "")
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	exceptions.PanicIfBadRequest(err, requestId, []string{"Error Generate Password"}, service.Logger)

	// Begin Transcation
	tx := service.DB.Begin()
	exceptions.PanicIfError(tx.Error, requestId, service.Logger)

	userEntity := &entity.User{
		Id:              utilities.RandomUUID(),
		Phone:           createUserRequest.Phone,
		Password:        string(bcryptPassword),
		IdDesa:          createUserRequest.IdDesa,
		IsActive:        1,
		IdLimitPayLater: "1006588e-da08-4e1b-8cd4-c14fff9059e1", //default limit 0
		AccountType:     1,                                      // 1 Normal 2 Merchant
		StatusSurvey:    0,                                      // 0 Blum survey 1 sudah survey
		CreatedDate:     time.Now(),
	}

	// Save user to database
	err = service.UserRepositoryInterface.CreateUser(tx, userEntity)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error create user"}, service.Logger, service.DB)

	userProfileEntity := &entity.UserProfile{
		Id:          utilities.RandomUUID(),
		IdUser:      userEntity.Id,
		NoIdentitas: createUserRequest.NoIdentitas,
		NamaLengkap: createUserRequest.NamaLengkap,
		Email:       emailLowerCase,
		CreatedDate: time.Now(),
	}

	// Save user profile to database
	err = service.UserProfileRepositoryInterface.CreateUserProfile(tx, userProfileEntity)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error create user profile"}, service.Logger, service.DB)

	pointEntity := &entity.Point{
		Id:          utilities.RandomUUID(),
		IdUser:      userEntity.Id,
		JmlPoint:    0,
		StatusPoint: 1,
		CreatedDate: time.Now(),
	}

	// Save point to database
	err = service.PointRepositoryInterface.CreatePoint(tx, pointEntity)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error create point"}, service.Logger, service.DB)

	commit := tx.Commit()
	exceptions.PanicIfError(commit.Error, requestId, service.Logger)
}

func (service *UserServiceImplementation) FindUserById(requestid string, idUser string) (userResponse response.FindUserIdResponse) {
	user, err := service.UserRepositoryInterface.FindUserById(service.DB, idUser)
	exceptions.PanicIfError(err, requestid, service.Logger)
	userResponse = response.ToFindUserIdResponse(user)
	return userResponse
}
