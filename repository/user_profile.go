package repository

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type UserProfileRepositoryInterface interface {
	CreateUserProfile(db *gorm.DB, userProfile *entity.UserProfile) error
	FindUserByEmail(db *gorm.DB, email string) (*entity.UserProfile, error)
	FindUserByNoIdentitas(db *gorm.DB, NoIdentitas string) (*entity.UserProfile, error)
}

type UserProfileRepositoryImplementation struct {
	DB *config.Database
}

func NewUserProfileRepository(
	db *config.Database,
) UserProfileRepositoryInterface {
	return &UserProfileRepositoryImplementation{
		DB: db,
	}
}

func (repository *UserProfileRepositoryImplementation) CreateUserProfile(db *gorm.DB, userProfile *entity.UserProfile) error {
	result := db.Create(userProfile)
	return result.Error
}

func (repository *UserProfileRepositoryImplementation) FindUserByEmail(db *gorm.DB, email string) (*entity.UserProfile, error) {
	userProfile := &entity.UserProfile{}
	result := db.
		Joins("User").
		Find(userProfile, "users_profile.email = ?", email)
	return userProfile, result.Error
}

func (repository *UserProfileRepositoryImplementation) FindUserByNoIdentitas(db *gorm.DB, NoIdentitas string) (*entity.UserProfile, error) {
	userProfile := &entity.UserProfile{}
	result := db.
		Joins("User").
		Find(userProfile, "users_profile.no_identitas = ?", NoIdentitas)
	return userProfile, result.Error
}
