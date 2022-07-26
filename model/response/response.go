package response

type Response struct {
	Code  int         `json:"code"`
	Mssg  string      `json:"message"`
	Data  interface{} `json:"data"`
	Error []string    `json:"error"`
}
