package service

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"

// 	"github.com/sirupsen/logrus"
// 	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
// )

// type SmsServiceInterface interface {
// 	SendOTP(requestId string, phone string, otpCode string)
// }

// type SmsServiceImplementation struct {
// 	Logger *logrus.Logger
// }

// func NewSmsService(
// 	logger *logrus.Logger,
// ) SmsServiceInterface {
// 	return &SmsServiceImplementation{
// 		Logger: logger,
// 	}
// }

// func (service *SmsServiceImplementation) SendOTP(requestId string, phone string, otpCode string) {
// 	message := fmt.Sprintf("Kode Verifikasi Akun Bupda Bali Anda adalah: %s *JANGAN BERIKAN KODE INI KEPADA SIAPAPUN, TERMASUK PIHAK BUPDA BALI* Hubungi xxxxxxxx untuk bantuan.", otpCode)

// 	postBody, _ := json.Marshal(map[string]string{
// 		"userkey": config.GetConfig().Sms.UserKey,
// 		"passkey": config.GetConfig().Sms.PassKey,
// 		"to":      phone,
// 		"message": message,
// 	})

// 	responseBody := bytes.NewBuffer(postBody)
// 	//Leverage Go's HTTP Post function to make request
// 	resp, err := http.Post("https://console.zenziva.net/masking/api/sendOTP", "application/json", responseBody)
// 	//Handle Error
// 	if err != nil {
// 		log.Fatalf("An Error Occured %v", err)
// 	}
// 	defer resp.Body.Close()
// 	//Read the response body
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	sb := string(body)
// 	fmt.Printf(sb)
// }
