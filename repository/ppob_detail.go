package repository

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type PpobDetailRepositoryInterface interface {
	CreateOrderPpobDetailPostpaidPdam(db *gorm.DB, ppobDetailPostpaidPdam *entity.PpobDetailPostpaidPdam) error
	CreateOrderPpobDetailPostpaidPln(db *gorm.DB, ppobDetailPostpaidPln *entity.PpobDetailPostpaidPln) error
	CreateOrderPpobDetailPrepaidPulsa(db *gorm.DB, ppobDetailPrepaidPulsa *entity.PpobDetailPrepaidPulsa) error
	CreateOrderPpobDetailPrepaidPln(db *gorm.DB, ppobDetailPrepaidPln *entity.PpobDetailPrepaidPln) error
	FindPpobDetailPrepaidPulsaById(db *gorm.DB, idOrderItemsPpob string) (*entity.PpobDetailPrepaidPulsa, error)
	FindPpobDetailPrepaidPlnById(db *gorm.DB, idOrderItemsPpob string) (*entity.PpobDetailPrepaidPln, error)
	FindPpobDetailPostpaidPlnById(db *gorm.DB, idOrderItemsPpob string) (*entity.PpobDetailPostpaidPln, error)
	UpdatePpobPrepaidPulsaById(db *gorm.DB, idOrderItemPpob string, ppobDetailUpdatePrepaidPulsa *entity.PpobDetailPrepaidPulsa) error
	UpdatePpobPrepaidPlnById(db *gorm.DB, idOrderItemPpob string, ppobDetailUpdatePrepaidPln *entity.PpobDetailPrepaidPln) error
}

type PpobDetailRepositoryImplementation struct {
	DB *config.Database
}

func NewPpobDetailRepository(
	db *config.Database,
) PpobDetailRepositoryInterface {
	return &PpobDetailRepositoryImplementation{
		DB: db,
	}
}

func (repository *PpobDetailRepositoryImplementation) CreateOrderPpobDetailPostpaidPdam(db *gorm.DB, ppobDetailPostpaidPdam *entity.PpobDetailPostpaidPdam) error {
	result := db.Create(&ppobDetailPostpaidPdam)
	return result.Error
}

func (repository *PpobDetailRepositoryImplementation) CreateOrderPpobDetailPostpaidPln(db *gorm.DB, ppobDetailPostpaidPln *entity.PpobDetailPostpaidPln) error {
	result := db.Create(&ppobDetailPostpaidPln)
	return result.Error
}

func (repository *PpobDetailRepositoryImplementation) CreateOrderPpobDetailPrepaidPulsa(db *gorm.DB, ppobDetailPrepaidPulsa *entity.PpobDetailPrepaidPulsa) error {
	result := db.Create(&ppobDetailPrepaidPulsa)
	return result.Error
}

func (repository *PpobDetailRepositoryImplementation) CreateOrderPpobDetailPrepaidPln(db *gorm.DB, ppobDetailPrepaidPln *entity.PpobDetailPrepaidPln) error {
	result := db.Create(&ppobDetailPrepaidPln)
	return result.Error
}

func (repository *PpobDetailRepositoryImplementation) FindPpobDetailPrepaidPulsaById(db *gorm.DB, idOrderItemsPpob string) (*entity.PpobDetailPrepaidPulsa, error) {
	ppobDetailPrepaidPulsa := &entity.PpobDetailPrepaidPulsa{}
	result := db.
		Find(ppobDetailPrepaidPulsa, "id_order_item_ppob = ?", idOrderItemsPpob)
	return ppobDetailPrepaidPulsa, result.Error
}

func (repository *PpobDetailRepositoryImplementation) FindPpobDetailPrepaidPlnById(db *gorm.DB, idOrderItemsPpob string) (*entity.PpobDetailPrepaidPln, error) {
	ppobDetailPrepaidPln := &entity.PpobDetailPrepaidPln{}
	result := db.
		Find(ppobDetailPrepaidPln, "id_order_item_ppob = ?", idOrderItemsPpob)
	return ppobDetailPrepaidPln, result.Error
}

func (repository *PpobDetailRepositoryImplementation) FindPpobDetailPostpaidPlnById(db *gorm.DB, idOrderItemsPpob string) (*entity.PpobDetailPostpaidPln, error) {
	ppobDetailPostpaidPln := &entity.PpobDetailPostpaidPln{}
	result := db.
		Find(ppobDetailPostpaidPln, "id_order_item_ppob = ?", idOrderItemsPpob)
	return ppobDetailPostpaidPln, result.Error
}

func (repository *PpobDetailRepositoryImplementation) UpdatePpobPrepaidPulsaById(db *gorm.DB, idOrderItemPpob string, ppobDetailUpdatePrepaidPulsa *entity.PpobDetailPrepaidPulsa) error {
	ppobDetailPrepaidPulsa := &entity.PpobDetailPrepaidPulsa{}
	result := db.
		Model(ppobDetailPrepaidPulsa).
		Where("id = ?", idOrderItemPpob).
		Updates(ppobDetailUpdatePrepaidPulsa)
	return result.Error
}

func (repository *PpobDetailRepositoryImplementation) UpdatePpobPrepaidPlnById(db *gorm.DB, idOrderItemPpob string, ppobDetailUpdatePrepaidPln *entity.PpobDetailPrepaidPln) error {
	ppobDetailPrepaidPln := &entity.PpobDetailPrepaidPln{}
	result := db.
		Model(ppobDetailPrepaidPln).
		Where("id = ?", idOrderItemPpob).
		Updates(ppobDetailUpdatePrepaidPln)
	return result.Error
}
