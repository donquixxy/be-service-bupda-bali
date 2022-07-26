package response

import "github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"

type FindUserShippingAddress struct {
	Id               string  `json:"id"`
	AlamatPengiriman string  `json:"alamat_pengiriman"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	Radius           float64 `json:"radius"`
	StatusPrimary    int     `json:"status_primary"`
	Catatan          string  `json:"catatan"`
}

func ToFindUserShippingAddressResponse(userShippingAddresss []entity.UserShippingAddress) (userShippingAddressResponses []FindUserShippingAddress) {
	for _, userShippingAddress := range userShippingAddresss {
		var userShippingAddressResponse FindUserShippingAddress
		userShippingAddressResponse.Id = userShippingAddress.Id
		userShippingAddressResponse.AlamatPengiriman = userShippingAddress.AlamatPengiriman
		userShippingAddressResponse.Latitude = userShippingAddress.Latitude
		userShippingAddressResponse.Longitude = userShippingAddress.Longitude
		userShippingAddressResponse.Radius = userShippingAddress.Radius
		userShippingAddressResponse.StatusPrimary = userShippingAddress.StatusPrimary
		userShippingAddressResponse.Catatan = userShippingAddress.Catatan
		userShippingAddressResponses = append(userShippingAddressResponses, userShippingAddressResponse)
	}
	return userShippingAddressResponses
}
