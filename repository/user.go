package repository

import (
	"log"
	"time"

	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	CreateUser(db *gorm.DB, user *entity.User) error
	UpdateUser(db *gorm.DB, idUser string, user *entity.User) error
	SaveUserInveliToken(db *gorm.DB, idUser string, user *entity.User) error
	FindUserByPhone(db *gorm.DB, phone string) (*entity.User, error)
	FindUserById(db *gorm.DB, idUser string) (*entity.UserProfile, error)
	FindUserById2(db *gorm.DB, idUser string) (*entity.User, error)
	SaveUserRefreshToken(DB *gorm.DB, idUser string, refreshToken string) (int64, error)
	FindUserByIdAndRefreshToken(DB *gorm.DB, idUser string, refresh_token string) (*entity.User, error)
	SaveUserAccount(db *gorm.DB, userAccounts []*entity.UserAccount) error
	GetUserAccountByID(db *gorm.DB, idUser string) (*entity.UserAccount, error)
	GetUserAccountPaylaterByID(db *gorm.DB, idUser string) (*entity.UserAccount, error)
	GetUserAccountBimaByID(db *gorm.DB, idUser string) (*entity.UserAccount, error)
	GetUserPayLaterFlagThisMonth(db *gorm.DB, idUser string) (*entity.UsersPaylaterFlag, error)
	CreateUserPayLaterFlag(db *gorm.DB, userPayLaterFlag *entity.UsersPaylaterFlag) error
	UpdateUserPayLaterFlag(db *gorm.DB, idUser string, userPayLaterFlag *entity.UsersPaylaterFlag) error
	GetUserPaylaterList(db *gorm.DB, nik string) (*entity.UserGetPaylater, error)
}

type UserRepositoryImplementation struct {
	DB *config.Database
}

func NewUserRepository(
	db *config.Database,
) UserRepositoryInterface {
	return &UserRepositoryImplementation{
		DB: db,
	}
}

func (repository *UserRepositoryImplementation) GetUserPaylaterList(db *gorm.DB, nik string) (*entity.UserGetPaylater, error) {
	var userPayLaterFlag *entity.UserGetPaylater
	result := db.
		Where("nik = ?", nik).
		Find(&userPayLaterFlag)
	return userPayLaterFlag, result.Error
}

func (repository *UserRepositoryImplementation) UpdateUserPayLaterFlag(db *gorm.DB, idUser string, userPayLaterFlag *entity.UsersPaylaterFlag) error {
	result := db.
		Model(&entity.UsersPaylaterFlag{}).
		Where("id_user = ?", idUser).
		Updates(userPayLaterFlag)
	return result.Error
}

func (repository *UserRepositoryImplementation) CreateUserPayLaterFlag(db *gorm.DB, userPayLaterFlag *entity.UsersPaylaterFlag) error {
	result := db.Create(userPayLaterFlag)
	return result.Error
}

func (repository *UserRepositoryImplementation) GetUserPayLaterFlagThisMonth(db *gorm.DB, idUser string) (*entity.UsersPaylaterFlag, error) {
	userPayLaterFlag := &entity.UsersPaylaterFlag{}
	var month time.Month
	now := time.Now()
	day := now.Day()
	if day < 25 {
		month = now.Month()
	} else if day >= 25 {
		month = now.AddDate(0, 1, 0).Month()
	}

	result := db.
		Where("id_user = ?", idUser).
		Where("MONTH(paylater_date) = ?", int(month)).
		First(userPayLaterFlag)
	return userPayLaterFlag, result.Error
}

func (repository *UserRepositoryImplementation) FindUserById2(db *gorm.DB, idUser string) (*entity.User, error) {
	user := &entity.User{}
	result := db.
		Where("id = ?", idUser).
		Find(user)
	return user, result.Error
}

func (repository *UserRepositoryImplementation) GetUserAccountPaylaterByID(db *gorm.DB, idUser string) (*entity.UserAccount, error) {
	userAccount := &entity.UserAccount{}
	result := db.
		Where("id_user = ?", idUser).
		Where("account_name = ?", "Simpanan Khusus").
		Find(userAccount)
	return userAccount, result.Error
}

func (repository *UserRepositoryImplementation) GetUserAccountBimaByID(db *gorm.DB, idUser string) (*entity.UserAccount, error) {
	userAccount := &entity.UserAccount{}
	result := db.
		Where("id_user = ?", idUser).
		Where("account_name LIKE ?", "%Tabungan Bima%").
		Find(userAccount)
	return userAccount, result.Error
}

func (repository *UserRepositoryImplementation) GetUserAccountByID(db *gorm.DB, idUser string) (*entity.UserAccount, error) {
	userAccounts := &entity.UserAccount{}
	result := db.
		Where("id_user = ?", idUser).
		Find(&userAccounts)
	return userAccounts, result.Error
}

func (repository *UserRepositoryImplementation) CreateUser(db *gorm.DB, user *entity.User) error {
	result := db.Create(user)
	return result.Error
}

func (repository *UserRepositoryImplementation) SaveUserInveliToken(db *gorm.DB, idUser string, user *entity.User) error {
	result := db.
		Model(&entity.User{}).
		Where("id = ?", idUser).
		Updates(user)
	return result.Error
}

func (repository *UserRepositoryImplementation) SaveUserAccount(db *gorm.DB, userAccounts []*entity.UserAccount) error {
	result := db.Create(&userAccounts)
	return result.Error
}

func (repository *UserRepositoryImplementation) UpdateUser(db *gorm.DB, idUser string, userUpdate *entity.User) error {
	user := &entity.User{}
	log.Println("userUpdate", userUpdate)
	result := db.
		Model(user).
		Where("id = ?", idUser).
		Updates(userUpdate)
	return result.Error
}

func (repository *UserRepositoryImplementation) FindUserById(db *gorm.DB, idUser string) (*entity.UserProfile, error) {
	userProfile := &entity.UserProfile{}
	result := db.
		Joins("User").
		Preload("User.Desa").
		Find(userProfile, "User.Id = ?", idUser)
	return userProfile, result.Error
}

func (repository *UserRepositoryImplementation) FindUserByPhone(db *gorm.DB, phone string) (*entity.User, error) {
	user := &entity.User{}
	result := db.
		Where("phone = ?", phone).
		Where("is_delete = ?", 0).
		Find(user)
	return user, result.Error
}

func (repository *UserRepositoryImplementation) FindUserByIdAndRefreshToken(db *gorm.DB, idUser string, refreshToken string) (*entity.User, error) {
	user := &entity.User{}
	result := db.
		Where("id = ?", idUser).
		Where("refresh_token = ?", refreshToken).
		First(user)
	return user, result.Error
}

func (repository *UserRepositoryImplementation) SaveUserRefreshToken(db *gorm.DB, idUser string, refreshToken string) (int64, error) {
	result := db.
		Model(entity.User{}).
		Where("id = ?", idUser).
		Updates(entity.User{
			RefreshToken: refreshToken,
		})
	return result.RowsAffected, result.Error
}
