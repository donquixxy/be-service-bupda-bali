package bigis

type Response struct {
	Message      string       `json:"message"`
	DataResponse DataResponse `json:"data"`
}

type DataResponse struct {
	Nik    string `json:"nik"`
	Name   string `json:"name"`
	Alamat string `json:"alamat"`
	Phone  string `json:"phone"`
}
