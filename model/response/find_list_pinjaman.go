package response

import (
	"log"
	"time"

	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
)

type ListPinjamanResponse struct {
	Id               string    `json:"id"`
	IdOrder          string    `json:"id_order"`
	IdUser           string    `json:"id_user"`
	IdDesa           string    `json:"id_desa"`
	JmlTagihan       float64   `json:"jml_tagihan"`
	BungaPinjaman    float64   `json:"bunga_pinjaman"`
	BiayaAdmin       float64   `json:"biaya_admin"`
	Total            float64   `json:"total"`
	TglPeminjaman    time.Time `json:"tanggal_pinjaman"`
	TglJatuhTempo    time.Time `json:"tanggal_jatuh_tempo"`
	StatusPembayaran int       `json:"status_pembayaran"`
}

func ToListPinjamanResponses(listPinjamans []entity.ListPinjaman) (listPinjamanResponses []ListPinjamanResponse) {
	log.Println("list pinjaman response", listPinjamans)
	for _, listPinjaman := range listPinjamans {
		listPinjamanResponse := &ListPinjamanResponse{}
		listPinjamanResponse.Id = listPinjaman.Id
		listPinjamanResponse.IdOrder = listPinjaman.IdOrder
		listPinjamanResponse.IdUser = listPinjaman.IdUser
		listPinjamanResponse.IdDesa = listPinjaman.IdDesa
		listPinjamanResponse.JmlTagihan = listPinjaman.JmlTagihan
		listPinjamanResponse.BungaPinjaman = listPinjaman.BungaPinjaman
		listPinjamanResponse.BiayaAdmin = listPinjaman.BiayaAdmin
		listPinjamanResponse.Total = listPinjaman.Total
		listPinjamanResponse.TglPeminjaman = listPinjaman.TglPeminjaman
		listPinjamanResponse.TglJatuhTempo = listPinjaman.TglJatuhTempo.Time
		listPinjamanResponse.StatusPembayaran = listPinjaman.StatusPembayaran
		listPinjamanResponses = append(listPinjamanResponses, *listPinjamanResponse)
	}
	return listPinjamanResponses
}

func ToListPinjamanResponse(listPinjaman *entity.ListPinjaman) (listPinjamanResponse ListPinjamanResponse) {
	listPinjamanResponse.Id = listPinjaman.Id
	listPinjamanResponse.IdOrder = listPinjaman.IdOrder
	listPinjamanResponse.IdUser = listPinjaman.IdUser
	listPinjamanResponse.IdDesa = listPinjaman.IdDesa
	listPinjamanResponse.JmlTagihan = listPinjaman.JmlTagihan
	listPinjamanResponse.BungaPinjaman = listPinjaman.BungaPinjaman
	listPinjamanResponse.BiayaAdmin = listPinjaman.BiayaAdmin
	listPinjamanResponse.Total = listPinjaman.Total
	listPinjamanResponse.TglPeminjaman = listPinjaman.TglPeminjaman
	listPinjamanResponse.TglJatuhTempo = listPinjaman.TglJatuhTempo.Time
	listPinjamanResponse.StatusPembayaran = listPinjaman.StatusPembayaran
	return listPinjamanResponse
}

