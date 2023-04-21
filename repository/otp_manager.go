package repository

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type OtpManagerRepositoryInterface interface {
	CreateOtp(db *gorm.DB, otpManager *entity.OtpManager) error
	UpdateOtp(db *gorm.DB, idOtpManager string, otpManager *entity.OtpManager) error
	FindOtpByPhone(db *gorm.DB, phone string) (*entity.OtpManager, error)
}

type OtpManagerRepositoryImplementation struct {
	DB *config.Database
}

func NewOtpManagerRepository(
	db *config.Database,
) OtpManagerRepositoryInterface {
	return &OtpManagerRepositoryImplementation{
		DB: db,
	}
}

func (repository *OtpManagerRepositoryImplementation) CreateOtp(db *gorm.DB, otpManager *entity.OtpManager) error {
	result := db.Create(&otpManager)
	return result.Error
}

func (repository *OtpManagerRepositoryImplementation) UpdateOtp(db *gorm.DB, idOtpManager string, otpManager *entity.OtpManager) error {
	updateOtp := make(map[string]interface{})
	updateOtp["otp_code"] = otpManager.OtpCode
	updateOtp["otp_experied_at"] = otpManager.OtpExperiedAt
	updateOtp["freeze_due_date"] = otpManager.FreezeDueDate
	updateOtp["phone_limit"] = otpManager.PhoneLimit
	result := db.
		Model(entity.OtpManager{}).
		Where("id = ?", idOtpManager).
		Updates(&updateOtp)
	return result.Error
}

func (repository *OtpManagerRepositoryImplementation) FindOtpByPhone(db *gorm.DB, phone string) (*entity.OtpManager, error) {
	otpManager := &entity.OtpManager{}
	result := db.
		Where("phone = ?", phone).
		Find(otpManager)
	return otpManager, result.Error
}
