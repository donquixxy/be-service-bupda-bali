package service

import (
	"errors"
	"log"
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
	invelirepository "github.com/tensuqiuwulu/be-service-bupda-bali/repository/inveli_repository"
	"github.com/tensuqiuwulu/be-service-bupda-bali/utilities"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthServiceInterface interface {
	Login(requestId string, loginRequest *request.LoginRequest) (loginResponse interface{})
	FirstTimeUbahPasswordInveli(requestId string, ubahPasswordInveliRequest *request.UbahPasswordInveliRequest) error
	NewToken(requestId string, refreshToken string) (token string)
	GenerateToken(user modelService.User) (token string, err error)
	GenerateRefreshToken(user modelService.User) (token string, err error)
	GetUserAccountInveli(inveliIDMember string, inveliAccessToken string, userID string)
	FirstTimeLoginInveli(phone string, passwordFromInveli string) string
}

type AuthServiceImplementation struct {
	DB                             *gorm.DB
	ConfigJwt                      config.Jwt
	Validate                       *validator.Validate
	Logger                         *logrus.Logger
	UserRepositoryInterface        repository.UserRepositoryInterface
	InveliAPIRespositoryInterface  invelirepository.InveliAPIRepositoryInterface
	UserProfileRepositoryInterface repository.UserProfileRepositoryInterface
	DesaRepositoryInterface        repository.DesaRepositoryInterface
}

func NewAuthService(
	db *gorm.DB,
	configJwt config.Jwt,
	validate *validator.Validate,
	logger *logrus.Logger,
	userRepositoryInterface repository.UserRepositoryInterface,
	inveliAPIRespositoryInterface invelirepository.InveliAPIRepositoryInterface,
	userProfileRepositoryInterface repository.UserProfileRepositoryInterface,
	desaRepositoryInterface repository.DesaRepositoryInterface,
) AuthServiceInterface {
	return &AuthServiceImplementation{
		DB:                             db,
		ConfigJwt:                      configJwt,
		Validate:                       validate,
		Logger:                         logger,
		UserRepositoryInterface:        userRepositoryInterface,
		InveliAPIRespositoryInterface:  inveliAPIRespositoryInterface,
		UserProfileRepositoryInterface: userProfileRepositoryInterface,
		DesaRepositoryInterface:        desaRepositoryInterface,
	}
}

func (service *AuthServiceImplementation) Login(requestId string, loginRequest *request.LoginRequest) (loginResponse interface{}) {
	var userModelService modelService.User

	request.ValidateRequest(service.Validate, loginRequest, requestId, service.Logger)

	// jika username tidak ditemukan
	user, _ := service.UserRepositoryInterface.FindUserByPhone(service.DB, loginRequest.Phone)
	if len(user.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("user not found"), requestId, []string{"not found"}, service.Logger)
	}

	if user.IsDelete == 1 {
		exceptions.PanicIfRecordNotFound(errors.New("user not found"), requestId, []string{"not found"}, service.Logger)
	}

	if user.IsActive == 1 {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
		exceptions.PanicIfBadRequest(err, requestId, []string{"Invalid Credentials"}, service.Logger)

		userModelService.Id = user.Id
		userModelService.IdDesa = user.IdDesa
		userModelService.AccountType = user.AccountType

		token, err := service.GenerateToken(userModelService)
		exceptions.PanicIfError(err, requestId, service.Logger)

		refreshToken, err := service.GenerateRefreshToken(userModelService)
		exceptions.PanicIfError(err, requestId, service.Logger)

		_, err = service.UserRepositoryInterface.SaveUserRefreshToken(service.DB, userModelService.Id, refreshToken)
		exceptions.PanicIfError(err, requestId, service.Logger)

		loginResponse = response.ToLoginResponse(token, refreshToken)

		if user.StatusPaylater == 2 {
			go service.FirstTimeLoginInveli(user.Phone, user.InveliPassword)
		}
		// service.FirstTimeLoginInveli(user.Phone, user.InveliPassword)

		// Get User Paylater List
		if user.IsPaylater == 0 {
			go func() {
				userResult, _ := service.UserRepositoryInterface.FindUserById(service.DB, user.Id)

				userPaylaterList, _ := service.UserRepositoryInterface.GetUserPaylaterList(service.DB, userResult.NoIdentitas)

				if len(userPaylaterList.Id) != 0 {
					userEntity := &entity.User{
						IsPaylater: 1,
					}
					service.UserRepositoryInterface.UpdateUser(service.DB, user.Id, userEntity)
				}
			}()
		}

		return loginResponse
	} else {
		exceptions.PanicIfUnauthorized(errors.New("account is not active"), requestId, []string{"not active"}, service.Logger)
		return nil
	}
}

func (service *AuthServiceImplementation) FirstTimeLoginInveli(phone string, passwordFromInveli string) string {
	loginResult := service.InveliAPIRespositoryInterface.InveliLogin(phone, passwordFromInveli)

	if len(loginResult.AccessToken) == 0 {
		log.Println("login inveli gagal")
		return ""
		// exceptions.PanicIfBadRequest(errors.New("gagal login to inveli"), "requestId", []string{"Invalid Credentials Inveli Login"}, service.Logger)
	}

	userResult, _ := service.UserRepositoryInterface.FindUserByPhone(service.DB, phone)
	if userResult.StatusPaylater == 2 {
		user := &entity.User{
			InveliAccessToken: loginResult.AccessToken,
			InveliIDMember:    loginResult.UserID,
		}

		if len(userResult.Id) == 0 {
			exceptions.PanicIfBadRequest(errors.New("user tidak ditemukan 1"), "requestId", []string{"User Not Found"}, service.Logger)
		}

		err := service.UserRepositoryInterface.SaveUserInveliToken(service.DB, userResult.Id, user)

		if err != nil {
			exceptions.PanicIfBadRequest(errors.New("gagal update token inveli"), "requestId", []string{"Failed Update Token Inveli"}, service.Logger)
		}

		return loginResult.AccessToken
	} else {
		user := &entity.User{
			InveliAccessToken: loginResult.AccessToken,
			InveliIDMember:    loginResult.UserID,
			StatusPaylater:    1,
		}

		if len(userResult.Id) == 0 {
			exceptions.PanicIfBadRequest(errors.New("user tidak ditemukan 1"), "requestId", []string{"User Not Found"}, service.Logger)
		}

		err := service.UserRepositoryInterface.SaveUserInveliToken(service.DB, userResult.Id, user)

		if err != nil {
			exceptions.PanicIfBadRequest(errors.New("gagal update token inveli"), "requestId", []string{"Failed Update Token Inveli"}, service.Logger)
		}

		return loginResult.AccessToken
	}

}

func (service *AuthServiceImplementation) GetUserAccountInveli(IDMember, AccessToken, IdUser string) {
	accountInfo, _ := service.InveliAPIRespositoryInterface.GetAccountInfo(IDMember, AccessToken)
	// fmt.Println("accountInfo : ", accountInfo)
	if accountInfo == nil {
		log.Println("akun belum aktif")
	} else {
		go func() {
			codeBIN, err := service.InveliAPIRespositoryInterface.GetKodeBIN(AccessToken)
			if err != nil {
				log.Println("Error Get code bin", err.Error())
			}
			var userAccounts []*entity.UserAccount
			for _, account := range accountInfo {
				userAccount := &entity.UserAccount{}
				userAccount.Id = utilities.RandomUUID()
				userAccount.IdUser = IdUser
				userAccount.IdAccount = account.ID
				userAccount.AccountName = account.AccountName2
				userAccount.IdProduct = account.ProductID
				userAccount.Code = account.Code
				userAccount.BIN = codeBIN
				userAccounts = append(userAccounts, userAccount)
			}

			user := &entity.User{
				StatusPaylater: 2,
			}

			service.UserRepositoryInterface.SaveUserInveliToken(service.DB, IdUser, user)
			if err != nil {
				exceptions.PanicIfBadRequest(errors.New("gagal update token inveli"), "requestId", []string{"Failed Update Token Inveli : ", err.Error()}, service.Logger)
			}

			err = service.UserRepositoryInterface.UpdateUser(service.DB, IdUser, user)
			if err != nil {
				log.Println("error update user : ", err.Error())
			}

			err = service.UserRepositoryInterface.SaveUserAccount(service.DB, userAccounts)
			if err != nil {
				log.Println("error save user account : ", err.Error())
			}
		}()

	}
}

func (service *AuthServiceImplementation) FirstTimeUbahPasswordInveli(requestId string, ubahPasswordInveliRequest *request.UbahPasswordInveliRequest) error {

	request.ValidateRequest(service.Validate, ubahPasswordInveliRequest, requestId, service.Logger)

	accessToken := service.FirstTimeLoginInveli(ubahPasswordInveliRequest.Phone, ubahPasswordInveliRequest.PasswordFromInveli)

	userResult, _ := service.UserRepositoryInterface.FindUserByPhone(service.DB, ubahPasswordInveliRequest.Phone)
	if len(userResult.Id) == 0 {
		exceptions.PanicIfBadRequest(errors.New("user not found"), requestId, []string{"User Not Found"}, service.Logger)
	}

	resp, err := service.InveliAPIRespositoryInterface.InveliUbahPassword(userResult.InveliIDMember, ubahPasswordInveliRequest.NewPassword, accessToken)
	if err != nil {
		exceptions.PanicIfBadRequest(errors.New("error inveli ubah password"), requestId, []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger)
	}

	if resp == nil {
		exceptions.PanicIfBadRequest(errors.New("error inveli ubah password"), requestId, []string{"error change password inveli"}, service.Logger)
	}

	// Hash password
	password := strings.ReplaceAll(ubahPasswordInveliRequest.NewPassword, " ", "")
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	exceptions.PanicIfBadRequest(err, requestId, []string{"Error Generate Password"}, service.Logger)

	userUpdateEntity := &entity.User{
		Password:       string(bcryptPassword),
		InveliPassword: ubahPasswordInveliRequest.NewPassword,
	}

	errr := service.UserRepositoryInterface.UpdateUser(service.DB, userResult.Id, userUpdateEntity)

	if errr != nil {
		exceptions.PanicIfBadRequest(errors.New("failed change password to db"), requestId, []string{"failed to update password user db"}, service.Logger)
	}

	userProfile, _ := service.UserProfileRepositoryInterface.FindUserProfileByIdUser(service.DB, userResult.Id)

	if len(userProfile.Id) == 0 {
		exceptions.PanicIfBadRequest(errors.New("user profile not found"), requestId, []string{"User Profile Not Found"}, service.Logger)
	}

	user, _ := service.UserRepositoryInterface.FindUserById(service.DB, userResult.Id)

	desa, _ := service.DesaRepositoryInterface.FindDesaById(service.DB, user.User.IdDesa)

	if len(desa.GroupIdBupda) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("groupd id not found"), requestId, []string{"groupd id not found"}, service.Logger)
	}

	errrr := service.InveliAPIRespositoryInterface.InveliUpdateMember(userResult, userProfile, accessToken, desa.GroupIdBupda)

	if errrr != nil {
		exceptions.PanicIfBadRequest(errors.New("failed activate account"), requestId, []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger)
	}

	return nil

}

func (service *AuthServiceImplementation) NewToken(requestId string, refreshToken string) (token string) {
	tokenParse, err := jwt.ParseWithClaims(refreshToken, &modelService.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(service.ConfigJwt.Key), nil
	})

	if !tokenParse.Valid {
		exceptions.PanicIfUnauthorized(err, requestId, []string{"invalid token"}, service.Logger)
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			exceptions.PanicIfUnauthorized(err, requestId, []string{"invalid token"}, service.Logger)
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			exceptions.PanicIfUnauthorized(err, requestId, []string{"expired token"}, service.Logger)
		} else {
			exceptions.PanicIfError(err, requestId, service.Logger)
		}
	}

	if claims, ok := tokenParse.Claims.(*modelService.TokenClaims); ok && tokenParse.Valid {
		user, err := service.UserRepositoryInterface.FindUserByIdAndRefreshToken(service.DB, claims.Id, refreshToken)
		exceptions.PanicIfRecordNotFound(err, requestId, []string{"User tidak ada"}, service.Logger)

		if user.StatusPaylater == 2 {
			service.FirstTimeLoginInveli(user.Phone, user.InveliPassword)
		}

		var userModelService modelService.User
		userModelService.Id = user.Id
		userModelService.IdDesa = user.IdDesa
		userModelService.AccountType = user.AccountType
		token, err := service.GenerateRefreshToken(userModelService)
		exceptions.PanicIfError(err, requestId, service.Logger)
		return token
	} else {
		err := errors.New("no claims")
		exceptions.PanicIfBadRequest(err, requestId, []string{"no claims"}, service.Logger)
		return ""
	}
}

func (service *AuthServiceImplementation) GenerateToken(user modelService.User) (token string, err error) {
	// Create the Claims
	claims := modelService.TokenClaims{
		Id:          user.Id,
		IdDesa:      user.IdDesa,
		AccountType: user.AccountType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(service.ConfigJwt.Tokenexpiredtime)).Unix(),
			Issuer:    "cyrilia",
		},
	}

	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenWithClaims.SignedString([]byte(service.ConfigJwt.Key))
	if err != nil {
		return "", err
	}
	return token, err
}

func (service *AuthServiceImplementation) GenerateRefreshToken(user modelService.User) (token string, err error) {
	// Create the Claims
	claims := modelService.TokenClaims{
		Id:          user.Id,
		IdDesa:      user.IdDesa,
		AccountType: user.AccountType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(0, 0, int(service.ConfigJwt.Refreshtokenexpiredtime)).Unix(),
			Issuer:    "cyrilia",
		},
	}

	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenWithClaims.SignedString([]byte(service.ConfigJwt.Key))
	if err != nil {
		return "", err
	}
	return token, err
}
