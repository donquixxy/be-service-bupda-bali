package repository

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type PointRepositoryInterface interface {
	CreatePoint(db *gorm.DB, point *entity.Point) error
	FindPointByUser(db *gorm.DB, idUser string) (*entity.Point, error)
}

type PointRepositoryImplementation struct {
	DB *config.Database
}

func NewPointRepository(
	db *config.Database,
) PointRepositoryInterface {
	return &PointRepositoryImplementation{
		DB: db,
	}
}

func (repository *PointRepositoryImplementation) CreatePoint(db *gorm.DB, point *entity.Point) error {
	result := db.Create(&point)
	return result.Error
}

func (repository *PointRepositoryImplementation) FindPointByUser(db *gorm.DB, idUser string) (*entity.Point, error) {
	point := &entity.Point{}
	result := db.
		Find(point, "points.id_user = ?", idUser)
	return point, result.Error
}
