package response

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/ppob"
)

type InquiryPrepaidPlnResponse struct {
	Status       string `json:"status"`
	CustomerId   string `json:"customer_id"`
	MeterNo      string `json:"meter_no"`
	SubscriberId string `json:"subscriber_id"`
	Name         string `json:"name"`
	SegmentPower string `json:"segment_power"`
	Message      string `json:"message"`
	Rc           string `json:"rc"`
}

func ToInquiryPrepaidPlnResponse(inquiryPrepaidPln *ppob.InquiryPrepaidPln) (inquiryPrepaidPlnResponse InquiryPrepaidPlnResponse) {
	inquiryPrepaidPlnResponse.Status = inquiryPrepaidPln.Data.Status
	inquiryPrepaidPlnResponse.CustomerId = inquiryPrepaidPln.Data.CustomerId
	inquiryPrepaidPlnResponse.MeterNo = inquiryPrepaidPln.Data.MeterNo
	inquiryPrepaidPlnResponse.SubscriberId = inquiryPrepaidPln.Data.SubscriberId
	inquiryPrepaidPlnResponse.Name = inquiryPrepaidPln.Data.Name
	inquiryPrepaidPlnResponse.SegmentPower = inquiryPrepaidPln.Data.SegmentPower
	inquiryPrepaidPlnResponse.Message = inquiryPrepaidPln.Data.Message
	inquiryPrepaidPlnResponse.Rc = inquiryPrepaidPln.Data.Rc
	return inquiryPrepaidPlnResponse
}
