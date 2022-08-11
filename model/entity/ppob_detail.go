package entity

import "gopkg.in/guregu/null.v4"

type PpobDetailPostpaidPdam struct {
	Id                  string        `gorm:"primaryKey;column:id;"`
	IdOrderItemPpob     string        `gorm:"column:id_order_item_ppob;"`
	OrderItemPpob       OrderItemPpob `gorm:"foreignKey:IdOrderItemPpob;"`
	TrId                int           `gorm:"column:tr_id;"`
	CustomerId          string        `gorm:"column:customer_id;"`
	CustomerName        string        `gorm:"column:customer_name;"`
	RefId               string        `gorm:"column:ref_id;"`
	Period              string        `gorm:"column:period;"`
	BillQty             int           `gorm:"column:bill_qty;"`
	DueDate             string        `gorm:"column:due_date;"`
	PdamName            string        `gorm:"column:pdam_name;"`
	PdamAddress         string        `gorm:"column:pdam_address;"`
	StampDuty           string        `gorm:"column:stamp_duty;"`
	Address             string        `gorm:"column:address;"`
	StatusTopUp         int           `gorm:"column:status_topup;"`
	LastBalance         float64       `gorm:"column:last_balance;"`
	TopupProccesingDate null.Time     `gorm:"column:topup_proccesing_date;"`
	TopupSuccessDate    null.Time     `gorm:"column:topup_success_date;"`
	TopupFailedDate     null.Time     `gorm:"column:topup_failed_date;"`
	JsonDetailTagihan   string        `gorm:"column:json_detail_tagihan;"`
}

func (PpobDetailPostpaidPdam) TableName() string {
	return "ppob_detail_postpaid_pdam"
}

type PpobDetailPostpaidPln struct {
	Id                  string        `gorm:"primaryKey;column:id;"`
	IdOrderItemPpob     string        `gorm:"column:id_order_item_ppob;"`
	OrderItemPpob       OrderItemPpob `gorm:"foreignKey:IdOrderItemPpob;"`
	RefId               string        `gorm:"column:ref_id;"`
	CustomerId          string        `gorm:"column:customer_id;"`
	CustomerName        string        `gorm:"column:customer_name;"`
	Tarif               string        `gorm:"column:tarif;"`
	Daya                int           `gorm:"column:daya;"`
	LembarTagihan       string        `gorm:"column:lembar_tagihan;"`
	Period              string        `gorm:"column:period;"`
	StatusTopUp         int           `gorm:"column:status_topup;"`
	LastBalance         float64       `gorm:"column:last_balance;"`
	TopupProccesingDate null.Time     `gorm:"column:topup_proccesing_date;"`
	TopupSuccessDate    null.Time     `gorm:"column:topup_success_date;"`
	TopupFailedDate     null.Time     `gorm:"column:topup_failed_date;"`
	JsonDetailTagihan   string        `gorm:"column:json_detail_tagihan;"`
}

func (PpobDetailPostpaidPln) TableName() string {
	return "ppob_detail_postpaid_pln"
}

type PpobDetailPrepaidPulsa struct {
	Id                  string        `gorm:"primaryKey;column:id;"`
	IdOrderItemPpob     string        `gorm:"column:id_order_item_ppob;"`
	OrderItemPpob       OrderItemPpob `gorm:"foreignKey:IdOrderItemPpob;"`
	ProductCode         string        `gorm:"column:product_code;"`
	ProductName         string        `gorm:"column:product_name;"`
	ProductDescription  string        `gorm:"column:product_description;"`
	CustomerId          string        `gorm:"column:customer_id;"`
	Operator            string        `gorm:"column:operator;"`
	ActivePeriod        string        `gorm:"column:active_period;"`
	StatusTopUp         int           `gorm:"column:status_topup;"`
	LastBalance         float64       `gorm:"column:last_balance;"`
	IconUrl             string        `gorm:"column:icon_url;"`
	TopupProccesingDate null.Time     `gorm:"column:topup_proccesing_date;"`
	TopupSuccessDate    null.Time     `gorm:"column:topup_success_date;"`
	TopupFailedDate     null.Time     `gorm:"column:topup_failed_date;"`
}

func (PpobDetailPrepaidPulsa) TableName() string {
	return "ppob_detail_prepaid_pulsa"
}

type PpobDetailPrepaidPln struct {
	Id                  string        `gorm:"primaryKey;column:id;"`
	IdOrderItemPpob     string        `gorm:"column:id_order_item_ppob;"`
	OrderItemPpob       OrderItemPpob `gorm:"foreignKey:IdOrderItemPpob;"`
	ProductCode         string        `gorm:"column:product_code;"`
	ProductName         string        `gorm:"column:product_name;"`
	ProductDescription  string        `gorm:"column:product_description;"`
	CustomerId          string        `gorm:"column:customer_id;"`
	MeterNo             string        `gorm:"column:meter_no;"`
	SubscriberId        string        `gorm:"column:subscriber_id;"`
	CustomerName        string        `gorm:"column:customer_name;"`
	SegmentPower        string        `gorm:"column:segment_power;"`
	StatusTopUp         int           `gorm:"column:status_topup;"`
	LastBalance         float64       `gorm:"column:last_balance;"`
	NoToken             string        `gorm:"column:no_token;"`
	TopupProccesingDate null.Time     `gorm:"column:topup_proccesing_date;"`
	TopupSuccessDate    null.Time     `gorm:"column:topup_success_date;"`
	TopupFailedDate     null.Time     `gorm:"column:topup_failed_date;"`
}

func (PpobDetailPrepaidPln) TableName() string {
	return "ppob_detail_prepaid_pln"
}
