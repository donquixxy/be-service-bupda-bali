package repository

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type UserShippingAddressRepositoryInterface interface {
	CreateUserShippingAddress(DB *gorm.DB, userShippingAddress *entity.UserShippingAddress) (*entity.UserShippingAddress, error)
	FindUserShippingAddressByIdUser(DB *gorm.DB, idUser string) ([]entity.UserShippingAddress, error)
	FindUserShippingAddressById(DB *gorm.DB, idUserShippingAddress string) (*entity.UserShippingAddress, error)
	DeleteUserShippingAddress(DB *gorm.DB, idUserShippingAddress string) error
	FindUserShippingAddressByAddress(DB *gorm.DB, address string) (*entity.UserShippingAddress, error)
}

type UserShippingAddressRepositoryImplementation struct {
	DB *config.Database
}

func NewUserShippingAddressRepository(db *config.Database) UserShippingAddressRepositoryInterface {
	return &UserShippingAddressRepositoryImplementation{
		DB: db,
	}
}

func (repository *UserShippingAddressRepositoryImplementation) FindUserShippingAddressByAddress(DB *gorm.DB, address string) (*entity.UserShippingAddress, error) {
	userShippingAddresss := &entity.UserShippingAddress{}
	results := DB.Where("alamat_pengiriman = ?", address).Find(userShippingAddresss)
	return userShippingAddresss, results.Error
}

func (repository *UserShippingAddressRepositoryImplementation) CreateUserShippingAddress(DB *gorm.DB, userAddress *entity.UserShippingAddress) (*entity.UserShippingAddress, error) {
	results := DB.Create(userAddress)
	return userAddress, results.Error
}

func (repository *UserShippingAddressRepositoryImplementation) DeleteUserShippingAddress(DB *gorm.DB, idUserShippingAddress string) error {
	result := DB.Where("id = ?", idUserShippingAddress).Delete(&entity.UserShippingAddress{})
	return result.Error
}

func (repository *UserShippingAddressRepositoryImplementation) FindUserShippingAddressByIdUser(DB *gorm.DB, idUser string) ([]entity.UserShippingAddress, error) {
	userShippingAddresss := []entity.UserShippingAddress{}
	results := DB.Where("id_user = ?", idUser).Find(&userShippingAddresss)
	return userShippingAddresss, results.Error
}

func (repository *UserShippingAddressRepositoryImplementation) FindUserShippingAddressById(DB *gorm.DB, idUserShippingAddress string) (*entity.UserShippingAddress, error) {
	userShippingAddresss := &entity.UserShippingAddress{}
	results := DB.Where("id = ?", idUserShippingAddress).Find(userShippingAddresss)
	return userShippingAddresss, results.Error
}
