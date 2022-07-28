package repository

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type OperatorPrefixRepositoryInterface interface {
	FindOperatorPrefixByPhone(db *gorm.DB, phonePrefix string) (*entity.OperatorPrefix, error)
}

type OperatorPrefixRepositoryImplementation struct {
	DB *config.Database
}

func NewOperatorPrefixRepository(
	db *config.Database,
) OperatorPrefixRepositoryInterface {
	return &OperatorPrefixRepositoryImplementation{
		DB: db,
	}
}

func (repository *OperatorPrefixRepositoryImplementation) FindOperatorPrefixByPhone(db *gorm.DB, phonePrefix string) (*entity.OperatorPrefix, error) {
	operatorPrefix := &entity.OperatorPrefix{}
	result := db.Find(operatorPrefix, "prefix_number = ?", phonePrefix)
	return operatorPrefix, result.Error
}
