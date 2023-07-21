package repository

import (
	"fmt"
	"time"

	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type OrderRepositoryInterface interface {
	CreateOrder(db *gorm.DB, order *entity.Order) error
	FindOrderByNumberOrder(db *gorm.DB, numberOrder string) (*entity.Order, error)
	FindOrderByUser(db *gorm.DB, idUser string, orderStatus int) ([]entity.Order, error)
	FindOrderById(db *gorm.DB, idOrder string) (*entity.Order, error)
	FindOrderByRefId(db *gorm.DB, refId string) (*entity.Order, error)
	UpdateOrderByIdOrder(db *gorm.DB, idOrder string, orderUpdate *entity.Order) error
	FindOrderPrepaidPulsaById(db *gorm.DB, idUser string, productType string) (*entity.Order, error)
	FindOrderPrepaidPlnById(db *gorm.DB, idUser string) (*entity.Order, error)
	FindOrderPayLaterById(db *gorm.DB, idUser string) ([]entity.Order, error)
	FindOrderPaylaterUnpaidById(db *gorm.DB, idUser string) ([]entity.Order, error)
	FindOrderPaylaterAllPaidById(db *gorm.DB, idUser string) ([]entity.Order, error)
	UpdateOrderPaylaterPaidStatus(db *gorm.DB, idUser string, orderUpdate *entity.Order) error
	UpdateOrderPaylaterPaidStatusByIdOrder(db *gorm.DB, idOrder string, orderUpdate *entity.Order) error
	GetOrderPaylaterPerBulan(db *gorm.DB, idUser string, month int) ([]entity.Order, string, string, error)
	GetOrderPaylaterPaidPerBulan(db *gorm.DB, idUser string, month int) ([]entity.Order, string, string, error)
	FindUnPaidPaylater(db *gorm.DB, idUser string) ([]entity.Order, error)
	FindOldestUnPaidPaylater(db *gorm.DB, idUser string, limit int) ([]entity.Order, error)
	FindOrderTotalPaylaterByMonth(db *gorm.DB, idUser string, month int) ([]entity.Order, error)
}

type OrderRepositoryImplementation struct {
	DB *config.Database
}

func NewOrderRepository(
	db *config.Database,
) OrderRepositoryInterface {
	return &OrderRepositoryImplementation{
		DB: db,
	}
}

func (repository *OrderRepositoryImplementation) FindUnPaidPaylater(db *gorm.DB, idUser string) ([]entity.Order, error) {
	orders := []entity.Order{}

	// Order type 1-8 Adalah sebuah TRANSAKSI PPOB / SEMBAKO
	// 9 adalah pelunasan
	result := db.
		Where("id_user = ?", idUser).
		Where("order_type < ?", 9).
		Where("payment_method = ?", "paylater").
		Where("paylater_paid_status = ?", 0).
		Order("order_date ASC").
		Find(&orders)

	return orders, result.Error
}

func (repository *OrderRepositoryImplementation) FindOldestUnPaidPaylater(db *gorm.DB, idUser string, limit int) ([]entity.Order, error) {
	orders := []entity.Order{}

	result := db.
		Where("id_user = ?", idUser).
		Where("order_type < ?", 9).
		Where("payment_method = ?", "paylater").
		Where("paylater_paid_status = ?", 0).
		Order("order_date ASC").
		Limit(limit).
		Find(&orders)

	return orders, result.Error
}

func (repository *OrderRepositoryImplementation) GetOrderPaylaterPaidPerBulan(db *gorm.DB, idOrder string, month int) ([]entity.Order, string, string, error) {
	orders := []entity.Order{}

	var startDate string
	var endDate string

	year := time.Now().Year()

	if month == 1 {
		startDate = fmt.Sprint(year-1) + "-" + fmt.Sprint(12) + "-26 " + "00:00:00"
		endDate = fmt.Sprint(year) + "-" + "0" + fmt.Sprint(month) + "-25 " + "23:59:59"
	} else {
		if month >= 10 {
			if (month - 1) == 9 {
				startDate = fmt.Sprint(year) + "-" + "0" + fmt.Sprint(month-1) + "-26 " + "00:00:00"
				endDate = fmt.Sprint(year) + "-" + fmt.Sprint(month) + "-25 " + "23:59:59"
			} else {
				startDate = fmt.Sprint(year) + "-" + fmt.Sprint(month-1) + "-26 " + "00:00:00"
				endDate = fmt.Sprint(year) + "-" + fmt.Sprint(month) + "-25 " + "23:59:59"
			}
		} else {
			startDate = fmt.Sprint(year) + "-" + "0" + fmt.Sprint(month-1) + "-26 " + "00:00:00"
			endDate = fmt.Sprint(year) + "-" + "0" + fmt.Sprint(month) + "-25 " + "23:59:59"
		}
	}

	result := db.
		Where("id_user = ?", idOrder).
		Where("payment_method = ?", "paylater").
		Where("order_date >= ?", startDate).
		Where("order_date <= ?", endDate).
		Where("paylater_paid_status ?", 1).
		Order("created_at desc").
		Find(&orders)

	return orders, startDate, endDate, result.Error
}

func (repository *OrderRepositoryImplementation) GetOrderPaylaterPerBulan(db *gorm.DB, idOrder string, month int) ([]entity.Order, string, string, error) {
	orders := []entity.Order{}

	var startDate string
	var endDate string

	// year := time.Now().Year()

	startDate = time.Date(time.Now().Year(), time.Month(month), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02 15:04:05")
	endDate = time.Date(time.Now().Year(), time.Month(month+1), 1, 0, 59, 59, 59, time.UTC).Format("2006-01-02 15:04:05")

	result := db.
		Where("id_user = ?", idOrder).
		Where("payment_method = ?", "paylater").
		Where("order_date >= ?", startDate).
		Where("order_date <= ?", endDate).
		Order("created_at desc").
		Find(&orders)

	return orders, startDate, endDate, result.Error
}

func (repository *OrderRepositoryImplementation) UpdateOrderPaylaterPaidStatus(db *gorm.DB, idUser string, orderUpdate *entity.Order) error {
	order := &entity.Order{}
	result := db.
		Model(order).
		Where("id_user = ?", idUser).
		Where("paylater_paid_status = ?", 0).
		Updates(orderUpdate)
	return result.Error
}

func (repository *OrderRepositoryImplementation) UpdateOrderPaylaterPaidStatusByIdOrder(db *gorm.DB, idOrder string, orderUpdate *entity.Order) error {
	order := &entity.Order{}
	result := db.
		Model(order).
		Where("id = ?", idOrder).
		Where("paylater_paid_status = ?", 0).
		Updates(orderUpdate)
	return result.Error
}

func (repository *OrderRepositoryImplementation) FindOrderPaylaterUnpaidById(db *gorm.DB, idUser string) ([]entity.Order, error) {
	orders := []entity.Order{}

	result := db.
		Where("id_user = ?", idUser).
		Where("payment_method = ?", "paylater").
		Where("paylater_paid_status = ?", 0).
		Find(&orders)

	return orders, result.Error
}

func (repository *OrderRepositoryImplementation) FindOrderPaylaterAllPaidById(db *gorm.DB, idUser string) ([]entity.Order, error) {
	orders := []entity.Order{}

	result := db.
		Where("id_user = ?", idUser).
		Where("payment_method = ?", "paylater").
		Find(&orders)

	return orders, result.Error
}

func (repository *OrderRepositoryImplementation) FindOrderTotalPaylaterByMonth(db *gorm.DB, idUser string, month int) ([]entity.Order, error) {
	orders := []entity.Order{}

	result := db.
		Where("id_user = ?", idUser).
		Where("payment_method = ?", "paylater").
		Where("MONTH(order_date) = ?", month).
		Find(&orders)
	return orders, result.Error
}

func (repository *OrderRepositoryImplementation) CreateOrder(db *gorm.DB, order *entity.Order) error {
	result := db.Create(&order)
	return result.Error
}

func (repository *OrderRepositoryImplementation) FindOrderByNumberOrder(db *gorm.DB, numberOrder string) (*entity.Order, error) {
	order := &entity.Order{}
	result := db.Find(order, "number_order = ?", numberOrder)
	return order, result.Error
}

func (repository *OrderRepositoryImplementation) FindOrderByUser(db *gorm.DB, idUser string, orderStatus int) ([]entity.Order, error) {
	var result *gorm.DB
	order := []entity.Order{}
	if orderStatus >= 0 {
		result = db.Order("order_date desc").Find(&order, "id_user = ? AND order_status = ?", idUser, orderStatus)
	} else {
		result = db.Order("order_date desc").Find(&order, "id_user = ?", idUser)
	}

	return order, result.Error
}

func (repository *OrderRepositoryImplementation) FindOrderById(db *gorm.DB, idUser string) (*entity.Order, error) {
	orders := &entity.Order{}
	result := db.Find(orders, "id = ?", idUser)
	return orders, result.Error
}

// Get data Order dengan payment paylater
// Berdasarkan periode saat ini
func (repository *OrderRepositoryImplementation) FindOrderPayLaterById(db *gorm.DB, idUser string) ([]entity.Order, error) {
	orders := []entity.Order{}
	// var month time.Month
	now := time.Now()
	day := now.Day()
	var date1 time.Time
	var date2 time.Time

	// Jika tanggal hari ini adalah Kurang dari 25 (1 S/D 24)
	if day < 25 {
		// Set value date 1 dengan
		// Tahun sekarang, Bulan sekarang -1 (Contoh skrg maret, -1. maka bulan skrg adalah Feb)
		// dan tanggal 24.
		// Kenapa tanggal 24 ?. Karena setiap periode akan reset setiap tanggal 25
		date1 = time.Date(now.Year(), now.Month()-1, 24, 23, 59, 59, 0, time.UTC)
		fmt.Println("Value date 1 :", date1)
		// Set value date 2 dengan
		// Tahun Sekarang, Bulan sekarang dan Tanggal 26
		// Intinya, untuk If yang ini adalah search between
		// Bulan sekarang -1 tanggal 24 hingga Bulan sekarang Tanggal 26
		date2 = time.Date(now.Year(), now.Month(), 26, 0, 0, 0, 0, time.UTC)
		fmt.Println("Value date 2 :", date2)
	} else if day >= 25 {
		// Jika tidak.
		// Set Value date 1 dengan
		// Tahun Sekarang. Bulan Sekarang. Tanggal 24.
		date1 = time.Date(now.Year(), now.Month(), 24, 23, 59, 59, 0, time.UTC)
		fmt.Println("Value date 1 bawah :", date1)
		// Set Value date 2 dengan
		// Tahun Sekarang. Bulan Sekarang +1 (next periode). Tanggal 26
		date2 = time.Date(now.Year(), now.Month()+1, 26, 0, 0, 0, 0, time.UTC)
		fmt.Println("Value date 2 bawah :", date2)
	}

	result := db.
		Where("id_user = ?", idUser).
		Where("payment_method = ?", "paylater").
		Where("order_date BETWEEN  ? AND ?", date1, date2).
		Order("created_at desc").
		Find(&orders)

	return orders, result.Error
}

func (repository *OrderRepositoryImplementation) FindOrderByRefId(db *gorm.DB, refId string) (*entity.Order, error) {
	orders := &entity.Order{}
	result := db.Find(orders, "ref_id = ?", refId)
	return orders, result.Error
}

func (repository *OrderRepositoryImplementation) FindOrderPrepaidPulsaById(db *gorm.DB, idUser string, productType string) (*entity.Order, error) {
	orders := &entity.Order{}
	result := db.Where("product_type = ?", productType).Find(orders, "id = ?", idUser)
	return orders, result.Error
}

func (repository *OrderRepositoryImplementation) FindOrderPrepaidPlnById(db *gorm.DB, idUser string) (*entity.Order, error) {
	orders := &entity.Order{}
	result := db.Find(orders, "id = ?", idUser)
	return orders, result.Error
}

func (repository *OrderRepositoryImplementation) UpdateOrderByIdOrder(db *gorm.DB, idOrder string, orderUpdate *entity.Order) error {
	order := &entity.Order{}
	result := db.
		Model(order).
		Where("id = ?", idOrder).
		Updates(orderUpdate)
	return result.Error
}
