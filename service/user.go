package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/bigis"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/inveli"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/request"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	modelService "github.com/tensuqiuwulu/be-service-bupda-bali/model/service"
	"github.com/tensuqiuwulu/be-service-bupda-bali/repository"
	invelirepository "github.com/tensuqiuwulu/be-service-bupda-bali/repository/inveli_repository"
	"github.com/tensuqiuwulu/be-service-bupda-bali/utilities"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type UserServiceInterface interface {
	CreateUserNonSuveyed(requestId string, createUserRequest *request.CreateUserRequest)
	CreateUserSuveyed(requestId string, createUserSurveyedRequest *request.CreateUserSurveyedRequest)
	FindUserById(requestId string, idUser string) (userResponse response.FindUserIdResponse)
	DeleteUserById(requestId string, idUser string)
	UpdateUserPassword(reqeustId string, idUser string, updateUserPasswordRequest *request.UpdateUserPasswordRequest)
	UpdateUserForgotPassword(reqeustId string, updateUserForgotPasswordRequest *request.UpdateUserForgotPasswordRequest)
	UpdateUserProfile(requestId string, idUser string, updateUserProfileRequest *request.UpdateUserProfileRequest)
	UpdateUserPhone(requestId string, idUser string, updateUserPhoneRequest *request.UpdateUserPhoneRequest)
	FindUserFromBigis(requestId string, requestUser *request.FindBigisResponsesRequest) (userBigisResponse response.FindUserFromBigisResponse)
}

type UserServiceImplementation struct {
	DB                             *gorm.DB
	Validate                       *validator.Validate
	ConfigJwt                      config.Jwt
	Logger                         *logrus.Logger
	UserRepositoryInterface        repository.UserRepositoryInterface
	UserProfileRepositoryInterface repository.UserProfileRepositoryInterface
	PointRepositoryInterface       repository.PointRepositoryInterface
	DesaRepositoryInterface        repository.DesaRepositoryInterface
	InveliRepositoryInterface      invelirepository.InveliAPIRepositoryInterface
}

func NewUserService(
	db *gorm.DB,
	validate *validator.Validate,
	configJwt config.Jwt,
	logger *logrus.Logger,
	userRepositoryInterface repository.UserRepositoryInterface,
	userProfileRepositoryInterface repository.UserProfileRepositoryInterface,
	pointRepositoryInterface repository.PointRepositoryInterface,
	desaRepositoryInterface repository.DesaRepositoryInterface,
	inveliRepositoryInterface invelirepository.InveliAPIRepositoryInterface,
) UserServiceInterface {
	return &UserServiceImplementation{
		DB:                             db,
		Validate:                       validate,
		ConfigJwt:                      configJwt,
		Logger:                         logger,
		UserRepositoryInterface:        userRepositoryInterface,
		UserProfileRepositoryInterface: userProfileRepositoryInterface,
		PointRepositoryInterface:       pointRepositoryInterface,
		DesaRepositoryInterface:        desaRepositoryInterface,
		InveliRepositoryInterface:      inveliRepositoryInterface,
	}
}

func (service *UserServiceImplementation) FindUserFromBigis(requestId string, requestUser *request.FindBigisResponsesRequest) (userBigisResponse response.FindUserFromBigisResponse) {
	// Create Request
	body, _ := json.Marshal(map[string]interface{}{
		"nik": requestUser.Nik,
	})

	reqBody := ioutil.NopCloser(strings.NewReader(string(body)))

	urlString := "http://117.53.44.216:9070/api/v1/response"
	// URL
	url, _ := url.Parse(urlString)

	req := &http.Request{
		Method: "POST",
		URL:    url,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Body: reqBody,
	}

	reqDump, _ := httputil.DumpRequestOut(req, true)
	fmt.Printf("REQUEST:\n%s", string(reqDump))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	// Read response body
	data, _ := ioutil.ReadAll(resp.Body)
	// fmt.Printf("body: %s\n", data)

	defer resp.Body.Close()

	userFromBigis := &bigis.Response{}
	// fmt.Printf("body: %s\n", prepaidPriceList)

	if err = json.Unmarshal([]byte(data), &userFromBigis); err != nil {
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	if len(userFromBigis.DataResponse.Nik) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("user tidak ditemukan"), requestId, []string{"user tidak ditemukan"}, service.Logger)
	}

	desa, err := service.DesaRepositoryInterface.FindOneDesaByIdKelu(service.DB, userFromBigis.DataResponse.IdKelu)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(desa.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("desa tidak ditemukan"), requestId, []string{"desa tidak ditemukan"}, service.Logger)
	}

	userBigisResponse = response.ToFindUserFromBigisResponse(userFromBigis, desa)

	return userBigisResponse
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

func (service *UserServiceImplementation) DeleteUserById(requestId string, idUser string) {
	err := service.UserRepositoryInterface.UpdateUser(service.DB, idUser, &entity.User{
		IsDelete:     1,
		IsDeleteDate: null.NewTime(time.Now(), true),
	})
	exceptions.PanicIfError(err, requestId, service.Logger)
}

func (service *UserServiceImplementation) CreateUserSuveyed(requestId string, createUserRequest *request.CreateUserSurveyedRequest) {
	var err error

	request.ValidateRequest(service.Validate, createUserRequest, requestId, service.Logger)

	// Check No Identitas
	NoIdentitasCheck, err := service.UserProfileRepositoryInterface.FindUserByNoIdentitas(service.DB, createUserRequest.NoIdentitas)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(NoIdentitasCheck.Id) != 0 {
		exceptions.PanicIfRecordAlreadyExists(errors.New("no identitas already exist"), requestId, []string{"no identitas sudah digunakan"}, service.Logger)
	}

	// Check email if exsict
	var emailLowerCase string
	emailLowerCase = strings.ToLower(createUserRequest.Email)
	if len(emailLowerCase) == 0 {
		emailLowerCase = "bupdabali@gmail.com"
	} else {
		emailChek, err := service.UserProfileRepositoryInterface.FindUserByEmail(service.DB, emailLowerCase)
		exceptions.PanicIfError(err, requestId, service.Logger)
		if len(emailChek.Id) != 0 {
			exceptions.PanicIfRecordAlreadyExists(errors.New("email already exist"), requestId, []string{"Email sudah digunakan"}, service.Logger)
		}
	}

	// Check No Hp
	phoneCheck, err := service.UserRepositoryInterface.FindUserByPhone(service.DB, createUserRequest.Phone)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(phoneCheck.Id) != 0 {
		exceptions.PanicIfRecordAlreadyExists(errors.New("phone already exist"), requestId, []string{"phone sudah digunakan"}, service.Logger)
	}

	// Hash password
	// password := strings.ReplaceAll(createUserRequest.Password, " ", "")
	// bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// exceptions.PanicIfBadRequest(err, requestId, []string{"Error Generate Password"}, service.Logger)

	// Begin Transcation
	tx := service.DB.Begin()
	exceptions.PanicIfError(tx.Error, requestId, service.Logger)

	userEntity := &entity.User{
		Id:    utilities.RandomUUID(),
		Phone: createUserRequest.Phone,
		// Password:        string(bcryptPassword),
		IdDesa:          createUserRequest.IdDesa,
		IsActive:        1,
		IdLimitPayLater: "1006588e-da08-4e1b-8cd4-c14fff9059e1", //default limit 0
		AccountType:     1,                                      // 1 Normal 2 Merchant
		StatusSurvey:    1,                                      // 0 Blum survey 1 sudah survey
		CreatedDate:     time.Now(),
	}

	// Save user to database
	err = service.UserRepositoryInterface.CreateUser(tx, userEntity)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error create user"}, service.Logger, service.DB)

	userProfileEntity := &entity.UserProfile{
		Id:                    utilities.RandomUUID(),
		IdUser:                userEntity.Id,
		NoIdentitas:           createUserRequest.NoIdentitas,
		NamaLengkap:           createUserRequest.NamaLengkap,
		AlamatSesuaiIdentitas: createUserRequest.Alamat,
		Email:                 emailLowerCase,
		CreatedDate:           time.Now(),
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

	inveliRegistrationModel := &inveli.InveliRegistrationModel{
		Email:      emailLowerCase,
		Phone:      createUserRequest.Phone,
		NIK:        createUserRequest.NoIdentitas,
		Address:    createUserRequest.Alamat,
		MemberName: createUserRequest.NamaLengkap,
	}

	// Register to inveli
	err = service.InveliRepositoryInterface.InveliResgisration(inveliRegistrationModel)
	if err != nil {
		exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error create data inveli"}, service.Logger, service.DB)
	}

	commit := tx.Commit()
	exceptions.PanicIfError(commit.Error, requestId, service.Logger)
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
	var emailLowerCase string
	emailLowerCase = strings.ToLower(createUserRequest.Email)
	if len(emailLowerCase) == 0 {
		emailLowerCase = "bupdabali@gmail.com"
	} else {
		emailChek, err := service.UserProfileRepositoryInterface.FindUserByEmail(service.DB, emailLowerCase)
		exceptions.PanicIfError(err, requestId, service.Logger)
		if len(emailChek.Id) != 0 {
			exceptions.PanicIfRecordAlreadyExists(errors.New("email already exist"), requestId, []string{"Email sudah digunakan"}, service.Logger)
		}
	}

	// Check No Hp
	phoneCheck, err := service.UserRepositoryInterface.FindUserByPhone(service.DB, createUserRequest.Phone)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(phoneCheck.Id) != 0 {
		exceptions.PanicIfRecordAlreadyExists(errors.New("phone already exist"), requestId, []string{"phone sudah digunakan"}, service.Logger)
	}

	// Hash password
	// password := strings.ReplaceAll(createUserRequest.Password, " ", "")
	// bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// exceptions.PanicIfBadRequest(err, requestId, []string{"Error Generate Password"}, service.Logger)

	// Begin Transcation
	tx := service.DB.Begin()
	exceptions.PanicIfError(tx.Error, requestId, service.Logger)

	userEntity := &entity.User{
		Id:    utilities.RandomUUID(),
		Phone: createUserRequest.Phone,
		// Password:        string(bcryptPassword),
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

	inveliRegistrationModel := &inveli.InveliRegistrationModel{
		Email:      emailLowerCase,
		Phone:      createUserRequest.Phone,
		NIK:        createUserRequest.NoIdentitas,
		Address:    "nil",
		MemberName: createUserRequest.NamaLengkap,
	}

	// Register to inveli
	err = service.InveliRepositoryInterface.InveliResgisration(inveliRegistrationModel)
	if err != nil {
		exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error create data inveli"}, service.Logger, service.DB)
	}

	commit := tx.Commit()
	exceptions.PanicIfError(commit.Error, requestId, service.Logger)
}

func (service *UserServiceImplementation) FindUserById(requestId string, idUser string) (userResponse response.FindUserIdResponse) {
	user, err := service.UserRepositoryInterface.FindUserById(service.DB, idUser)
	exceptions.PanicIfError(err, requestId, service.Logger)
	userResponse = response.ToFindUserIdResponse(user)
	return userResponse
}

func (service *UserServiceImplementation) UpdateUserPassword(requestId string, idUser string, updateUserPasswordRequest *request.UpdateUserPasswordRequest) {
	var err error

	user, err := service.UserRepositoryInterface.FindUserById(service.DB, idUser)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(user.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("data user not found"), requestId, []string{"data user not found"}, service.Logger)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.User.Password), []byte(updateUserPasswordRequest.PasswordLama))
	exceptions.PanicIfBadRequest(err, requestId, []string{"Invalid Credentials"}, service.Logger)

	password := strings.ReplaceAll(updateUserPasswordRequest.PasswordBaru, " ", "")
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	exceptions.PanicIfBadRequest(err, requestId, []string{"Error Generate Password"}, service.Logger)

	userEntity := &entity.User{
		Password: string(bcryptPassword),
	}

	err = service.UserRepositoryInterface.UpdateUser(service.DB, idUser, userEntity)
	exceptions.PanicIfError(err, requestId, service.Logger)
}

func (service *UserServiceImplementation) UpdateUserForgotPassword(requestId string, updateUserForgotPasswordRequest *request.UpdateUserForgotPasswordRequest) {
	var err error

	// validate request
	request.ValidateRequest(service.Validate, updateUserForgotPasswordRequest, requestId, service.Logger)

	// validate form token
	service.VerifyFormToken(requestId, updateUserForgotPasswordRequest.FormToken)

	user, err := service.UserRepositoryInterface.FindUserByPhone(service.DB, updateUserForgotPasswordRequest.Phone)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(user.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("data user not found"), requestId, []string{"data user not found"}, service.Logger)
	}

	password := strings.ReplaceAll(updateUserForgotPasswordRequest.PasswordBaru, " ", "")
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	exceptions.PanicIfBadRequest(err, requestId, []string{"Error Generate Password"}, service.Logger)

	userEntity := &entity.User{
		Password: string(bcryptPassword),
	}

	err = service.UserRepositoryInterface.UpdateUser(service.DB, user.Id, userEntity)
	exceptions.PanicIfError(err, requestId, service.Logger)
}

func (service *UserServiceImplementation) UpdateUserProfile(requestId string, idUser string, updateUserProfileRequest *request.UpdateUserProfileRequest) {
	var err error

	// validate reqeust
	request.ValidateRequest(service.Validate, updateUserProfileRequest, requestId, service.Logger)

	userProfileEntity := &entity.UserProfile{
		NamaLengkap: updateUserProfileRequest.NamaLengkap,
		Email:       updateUserProfileRequest.Email,
	}

	err = service.UserProfileRepositoryInterface.UpdateUserProfile(service.DB, idUser, userProfileEntity)
	exceptions.PanicIfError(err, requestId, service.Logger)
}

func (service *UserServiceImplementation) UpdateUserPhone(requestId string, idUser string, updateUserPhoneRequest *request.UpdateUserPhoneRequest) {
	var err error

	// validate reqeust
	request.ValidateRequest(service.Validate, updateUserPhoneRequest, requestId, service.Logger)

	// validate form token
	service.VerifyFormToken(requestId, updateUserPhoneRequest.FormToken)

	userEntity := &entity.User{
		Phone: updateUserPhoneRequest.Phone,
	}

	err = service.UserRepositoryInterface.UpdateUser(service.DB, idUser, userEntity)
	exceptions.PanicIfError(err, requestId, service.Logger)
}
