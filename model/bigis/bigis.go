package bigis

type Response struct {
	Message      string       `json:"message"`
	DataResponse DataResponse `json:"data"`
}

type DataResponse struct {
	Nik       string `json:"nik"`
	Name      string `json:"name"`
	Alamat    string `json:"alamat"`
	Kecamatan string `json:"kecamatan"`
	Kelurahan string `json:"kelurahan"`
	IdKelu    int    `json:"idkelu"`
}
