package repository

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type ListPinjamanRepositoryInterface interface {
	CreateListPinjaman(db *gorm.DB, listPinjaman *entity.ListPinjaman) error
	FindListPinjamanByIdUser(db *gorm.DB, idUser string) ([]entity.ListPinjaman, error)
	FindListPinjamanById(db *gorm.DB, IdListPinjaman string) (*entity.ListPinjaman, error)
}

type ListPinjamanRepositoryImplementation struct {
	DB *config.Database
}

func NewListPinjamanRepository(
	db *config.Database,
) ListPinjamanRepositoryInterface {
	return &ListPinjamanRepositoryImplementation{
		DB: db,
	}
}

func (repository *ListPinjamanRepositoryImplementation) CreateListPinjaman(db *gorm.DB, listPinjaman *entity.ListPinjaman) error {
	err := db.Create(listPinjaman).Error
	if err != nil {
		return err
	}

	return nil
}

func (repository *ListPinjamanRepositoryImplementation) FindListPinjamanByIdUser(db *gorm.DB, idUser string) ([]entity.ListPinjaman, error) {
	listPinjaman := []entity.ListPinjaman{}

	err := db.Where("id_user = ?", idUser).Find(&listPinjaman).Error
	if err != nil {
		return nil, err
	}

	return listPinjaman, nil
}

func (repository *ListPinjamanRepositoryImplementation) FindListPinjamanById(db *gorm.DB, IdListPinjaman string) (*entity.ListPinjaman, error) {
	listPinjaman := &entity.ListPinjaman{}

	err := db.Where("id = ?", IdListPinjaman).First(listPinjaman).Error
	if err != nil {
		return nil, err
	}

	return listPinjaman, nil
}
