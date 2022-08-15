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
	FindPpobDetailPostpaidPdamById(db *gorm.DB, idOrderItemsPpob string) (*entity.PpobDetailPostpaidPdam, error)
	UpdatePpobPrepaidPulsaById(db *gorm.DB, idOrderItemPpob string, ppobDetailUpdatePrepaidPulsa *entity.PpobDetailPrepaidPulsa) error
	UpdatePpobPrepaidPlnById(db *gorm.DB, idOrderItemPpob string, ppobDetailUpdatePrepaidPln *entity.PpobDetailPrepaidPln) error
	UpdatePpobPostpaidPlnById(db *gorm.DB, idOrderItemPpob string, ppobDetailUpdatePostpaidPln *entity.PpobDetailPostpaidPln) error
	UpdatePpobPostpaidPdamById(db *gorm.DB, idOrderItemPpob string, ppobDetailUpdatePostpaidPdam *entity.PpobDetailPostpaidPdam) error
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

func (repository *PpobDetailRepositoryImplementation) FindPpobDetailPostpaidPdamById(db *gorm.DB, idOrderItemsPpob string) (*entity.PpobDetailPostpaidPdam, error) {
	ppobDetailPostpaidPdam := &entity.PpobDetailPostpaidPdam{}
	result := db.
		Find(ppobDetailPostpaidPdam, "id_order_item_ppob = ?", idOrderItemsPpob)
	return ppobDetailPostpaidPdam, result.Error
}

func (repository *PpobDetailRepositoryImplementation) UpdatePpobPrepaidPulsaById(db *gorm.DB, idOrderItemPpob string, ppobDetailUpdatePrepaidPulsa *entity.PpobDetailPrepaidPulsa) error {
	updatePrepaidPulsa := make(map[string]interface{})
	updatePrepaidPulsa["status_topup"] = ppobDetailUpdatePrepaidPulsa.StatusTopUp
	updatePrepaidPulsa["topup_proccesing_date"] = ppobDetailUpdatePrepaidPulsa.TopupProccesingDate
	updatePrepaidPulsa["last_balance"] = ppobDetailUpdatePrepaidPulsa.LastBalance
	result := db.
		Model(entity.PpobDetailPrepaidPulsa{}).
		Where("id = ?", idOrderItemPpob).
		Updates(updatePrepaidPulsa)
	return result.Error
}

func (repository *PpobDetailRepositoryImplementation) UpdatePpobPrepaidPlnById(db *gorm.DB, idOrderItemPpob string, ppobDetailUpdatePrepaidPln *entity.PpobDetailPrepaidPln) error {
	updatePrepaidPln := make(map[string]interface{})
	updatePrepaidPln["status_topup"] = ppobDetailUpdatePrepaidPln.StatusTopUp
	updatePrepaidPln["topup_proccesing_date"] = ppobDetailUpdatePrepaidPln.TopupProccesingDate
	updatePrepaidPln["last_balance"] = ppobDetailUpdatePrepaidPln.LastBalance
	updatePrepaidPln["no_token"] = ppobDetailUpdatePrepaidPln.NoToken
	result := db.
		Model(entity.PpobDetailPrepaidPln{}).
		Where("id = ?", idOrderItemPpob).
		Updates(updatePrepaidPln)
	return result.Error
}

func (repository *PpobDetailRepositoryImplementation) UpdatePpobPostpaidPlnById(db *gorm.DB, idOrderItemPpob string, ppobDetailUpdatePostpaidPln *entity.PpobDetailPostpaidPln) error {
	updatePostpaidPln := make(map[string]interface{})
	updatePostpaidPln["status_topup"] = ppobDetailUpdatePostpaidPln.StatusTopUp
	updatePostpaidPln["topup_proccesing_date"] = ppobDetailUpdatePostpaidPln.TopupProccesingDate
	updatePostpaidPln["last_balance"] = ppobDetailUpdatePostpaidPln.LastBalance
	result := db.
		Model(entity.PpobDetailPostpaidPln{}).
		Where("id = ?", idOrderItemPpob).
		Updates(updatePostpaidPln)
	return result.Error
}

func (repository *PpobDetailRepositoryImplementation) UpdatePpobPostpaidPdamById(db *gorm.DB, idOrderItemPpob string, ppobDetailUpdatePostpaidPdam *entity.PpobDetailPostpaidPdam) error {
	updatePostpaidPdam := make(map[string]interface{})
	updatePostpaidPdam["status_topup"] = ppobDetailUpdatePostpaidPdam.StatusTopUp
	updatePostpaidPdam["topup_proccesing_date"] = ppobDetailUpdatePostpaidPdam.TopupProccesingDate
	updatePostpaidPdam["last_balance"] = ppobDetailUpdatePostpaidPdam.LastBalance
	result := db.
		Model(entity.PpobDetailPostpaidPdam{}).
		Where("id = ?", idOrderItemPpob).
		Updates(ppobDetailUpdatePostpaidPdam)
	return result.Error
}
