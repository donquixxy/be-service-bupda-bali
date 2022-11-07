package invelirepository

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/machinebox/graphql"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/inveli"
)

type InveliAPIRepositoryInterface interface {
	InveliResgisration(inveliRegistrationModel *inveli.InveliRegistrationModel) error
	InveliLogin(username, password string) *inveli.InveliLoginModel
	InveliUbahPassword(id, password, token string) (interface{}, error)
	InveliUpdateMember(user *entity.User, userProfile *entity.UserProfile, accessToken string) error
	GetAccountInfo(IDMember, token string) ([]*inveli.InveliAcountInfo, error)
	InveliCreatePaylater(token string, IDMember string, AccountID string, Amount float64) error
}

type InveliAPIRepositoryImplementation struct {
}

func NewInveliAPIRepository() InveliAPIRepositoryInterface {
	return &InveliAPIRepositoryImplementation{}
}

func (r *InveliAPIRepositoryImplementation) InveliCreatePaylater(token string, IDMember string, AccountID string, Amount float64) error {

	client := graphql.NewClient("http://api-dev.cardlez.com:8089/query")
	req := graphql.NewRequest(`
		mutation ($inputLoanObject: LoanObjectInput!) {
			createLoan(loanInputParam: $inputLoanObject) {
		}
	`)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Var("inputLoanObject", map[string]interface{}{
		"accountID":             AccountID,
		"memberID":              IDMember,
		"loanProductID":         "7CC2EEEC-9DC0-4456-A755-53A50F28990C",
		"tenor":                 1,
		"loanAmount":            Amount,
		"interestPercent":       0,
		"isAutoApprove":         true,
		"memberLoanAttachments": "",
	})

	ctx := context.Background()

	var respData interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Println("error create pinjaman : ", err)
		return err
	}

	fmt.Println("response create pinjaman : ", respData)

	return nil
}

func (r *InveliAPIRepositoryImplementation) InveliUbahPassword(id, password, token string) (interface{}, error) {
	client := graphql.NewClient("http://api-dev.cardlez.com:8089/query")
	req := graphql.NewRequest(`
		mutation ($id: String!, $password: String!, $isactivation: Boolean!, $isnewmember: Boolean!) {
			changePassword(
				id: $id, 
				newPassword: $password, 
				isactivation: $isactivation, 
				isnewmember: $isnewmember)
		}
	`)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Var("id", id)
	req.Var("password", password)
	req.Var("isactivation", false)
	req.Var("isnewmember", true)
	ctx := context.Background()

	// fmt.Println("request : ", req)

	var respData interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Println("error change password : ", err)
		return respData, err
	}

	// fmt.Println("response ubah password : ", respData)

	return respData, nil
}

func (r *InveliAPIRepositoryImplementation) InveliUpdateMember(user *entity.User, userProfile *entity.UserProfile, accessToken string) error {
	client := graphql.NewClient("http://api-dev.cardlez.com:8089/query")

	// make a request
	req := graphql.NewRequest(`
	  mutation ($member: MemberInput!, $memberDetail: MemberDetailInput!) {
			updateMember(member: $member, memberDetail: $memberDetail) {
				id
			}
	  }
	`)

	// set any variables
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Var("member", map[string]interface{}{
		"id":               user.InveliIDMember,
		"address":          userProfile.AlamatSesuaiIdentitas,
		"alamatRumah2":     "",
		"phone":            user.Phone,
		"province":         "D1FE2F93-9A39-4710-8C03-E266C9F1BEE1",
		"city":             "CEABD366-D270-4286-93A9-418B274C9A23",
		"kecamatan":        "Setia Budi",
		"kelurahan":        "",
		"kodePos":          "Kuningan Timur",
		"isSendToCore":     true,
		"referralMemberID": "",
		"bankCode":         "",
		"bankID":           "9B2FC1C5-9F3A-44C5-915C-60E8653F32D6",
		"bankAccountName":  "",
		"bankAccountNo":    "",
		"profileImage":     "Ranti Puspita_Selfie_ND2.jpg",
		"birthDate":        "2021-08-20 00:00:00.0000000",
		"identityNumber":   userProfile.NoIdentitas,
	})

	// set any variables
	req.Var("memberDetail", map[string]interface{}{
		"namaLengkapTanpaSingkatan": userProfile.NamaLengkap,
		"statusPendidikan":          "D73D02FD-4F8A-4AB9-B4AD-4AC01963B8B0",
		"namaIbuKandung":            "ibu",
		"namaKontakDarurat":         "ibu",
		"nomorKontakDarurat":        user.Phone,
		"hubunganKontakDarurat":     "anak",
		"provinsi":                  "D1FE2F93-9A39-4710-8C03-E266C9F1BEE1",
		"kabupaten":                 "58FD0E0F-12F0-45CF-B960-9F0F5951F759",
		"kecamatan":                 "Jatinegara",
		"kelurahan":                 "Kampung Melayu",
		"kodePos":                   "13320",
		"phone":                     user.Phone,
		"provinsiTempatKerja":       "D1FE2F93-9A39-4710-8C03-E266C9F1BEE1",
		"kabupatenTempatKerja":      "CEABD366-D270-4286-93A9-418B274C9A23",
		"kecamatanTempatKerja":      "Setia Budi",
		"kelurahanTempatKerja":      "Kuningan Timur",
		"kodePosTempatKerja":        "12950",
		"kodePekerjaan":             "241BEE62-DC7C-4CAE-82AA-F44B75266B94",
		"namaTempatKerja":           "default temp",
		"kodeBidangUsaha":           "7342E607-ED2D-4E32-837E-9A41FEB3EC7F",
		"penghasilanKotorPerTahun":  "0",
		"kodeSumberPenghasilan":     "4C4C1204-514B-42AF-B47C-3E92B4EA05E5",
		"maritalStatus":             "5",
		"jumlahTanggungan":          "0",
		"noIdentitasPasangan":       "",
		"namaPasangan":              "",
		"tanggalLahirPasangan":      "",
		"pisahHarta":                "0",
		"fileNameFotoKTP":           "Ranti Puspita_KTP_ND2.jpg",
		"fileNameFotoSelfie":        "Ranti Puspita_Selfie_ND2.jpg",
		"fileNameFotoKTP64":         "",
		"fileNameFotoSelfie64":      "",
		"isNewFotoKTP":              true,
		"isNewFotoSelfie":           true,
		"fileVideo64":               "",
	})

	// fmt.Println("req : ", req)

	ctx := context.Background()
	var respData interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Println("error update member inveli : ", err)
		return err
	}
	// fmt.Println("Response update member", respData)

	return nil
}

func (r *InveliAPIRepositoryImplementation) InveliLogin(username, password string) *inveli.InveliLoginModel {

	// Create Request
	body, _ := json.Marshal(map[string]interface{}{
		"username": username,
		"password": password,
		"imei":     "",
	})

	reqBody := ioutil.NopCloser(strings.NewReader(string(body)))

	urlString := "http://api-dev.cardlez.com:8089/login/member"
	// URL
	url, _ := url.Parse(urlString)

	req := &http.Request{
		Method: "POST",
		URL:    url,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Body: reqBody,
	}

	// reqDump, _ := httputil.DumpRequestOut(req, true)
	// fmt.Printf("REQUEST:\n%s", string(reqDump))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("An Error Occured %v", err)
	}

	// Read response body
	data, _ := ioutil.ReadAll(resp.Body)
	// fmt.Printf("body: %s\n", data)

	defer resp.Body.Close()

	inveliLogin := &inveli.InveliLoginModel{}
	// fmt.Printf("body: %s\n", prepaidPriceList)

	if err = json.Unmarshal([]byte(data), &inveliLogin); err != nil {
		log.Printf("An Error Occured %v", err)
	}

	return inveliLogin
}

func (r *InveliAPIRepositoryImplementation) InveliResgisration(inveliRegistrationModel *inveli.InveliRegistrationModel) error {
	client := graphql.NewClient("http://api-dev.cardlez.com:8089/register")

	// make a request
	req := graphql.NewRequest(`
    mutation ($memberObject: MemberInput!) {
			registerMember(memberObject: $memberObject) 
    }
	`)

	// set any variables
	req.Var("memberObject", map[string]interface{}{
		"handphone":          inveliRegistrationModel.Phone,
		"memberName":         inveliRegistrationModel.MemberName,
		"birthPlace":         "",
		"birthDate":          "",
		"emailAddress":       inveliRegistrationModel.Email,
		"identityNumber":     inveliRegistrationModel.NIK,
		"identityType":       1,
		"noNPWP":             "",
		"gender":             1,
		"address":            inveliRegistrationModel.Address,
		"alamatRumah2":       "",
		"alamatTempatKerja1": "",
		"alamatTempatKerja2": "",
		"bankCode":           "014",
		"noInduk":            "",
		"fileNameFotoKTP":    "",
		"fileNameFotoSelfie": "",
		"bankAccountName":    "",
		"bankAccountNo":      "",
		"grade":              "bronze",
		"maritalStatus":      0,
		"recordStatus":       1,
		"nationality":        0,
		"memberType":         1,
		"isLocked":           false,
		"virtualAccountNo":   "",
		"goldStatus":         0,
		"bankID":             "9B2FC1C5-9F3A-44C5-915C-60E8653F32D6",
	})

	// run it and capture the response
	ctx := context.Background()
	var respData interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Println("error inveli registration : ", err)
		return err
	}

	return nil
}

func (r *InveliAPIRepositoryImplementation) GetAccountInfo(IDMember, token string) ([]*inveli.InveliAcountInfo, error) {
	client := graphql.NewClient("http://api-dev.cardlez.com:8089/query")

	req := graphql.NewRequest(` 
		query ($keyword: String!) {	
			accounts (search: MEMBERID, keyword: $keyword) {
				id
      	code
      	accountName
      	accountName2
      	recordStatus
      	productName
      	productID
      	memberName
      	memberID
     		memberBranchID
    		memberType
    		email
    		phone
    		balance
    		balanceMerchant
    		closingBalance
    		blockingBalance
    		productType
    		isPrimary
			}
		}
	`)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Var("keyword", IDMember)
	ctx := context.Background()
	var respData interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Println(err)
		return nil, err
	}

	fmt.Println("req", req)
	fmt.Println("response", respData)

	userInfos := []*inveli.InveliAcountInfo{}
	for _, value := range respData.(map[string]interface{})["accounts"].([]interface{}) {
		var userInfo inveli.InveliAcountInfo
		userInfo.ID = value.(map[string]interface{})["id"].(string)
		userInfo.Code = value.(map[string]interface{})["code"].(string)
		userInfo.AccountName = value.(map[string]interface{})["accountName"].(string)
		userInfo.AccountName2 = value.(map[string]interface{})["accountName2"].(string)
		userInfo.RecordStatus = int(value.(map[string]interface{})["recordStatus"].(float64))
		userInfo.ProductName = value.(map[string]interface{})["productName"].(string)
		userInfo.ProductID = value.(map[string]interface{})["productID"].(string)
		userInfo.MemberName = value.(map[string]interface{})["memberName"].(string)
		userInfo.MemberID = value.(map[string]interface{})["memberID"].(string)
		userInfo.MemberBranchID = value.(map[string]interface{})["memberBranchID"].(string)
		userInfo.MemberType = int(value.(map[string]interface{})["memberType"].(float64))
		userInfo.Email = value.(map[string]interface{})["email"].(string)
		userInfo.Phone = value.(map[string]interface{})["phone"].(string)
		userInfo.Balance = int(value.(map[string]interface{})["balance"].(float64))
		userInfo.BalanceMerchant = int(value.(map[string]interface{})["balanceMerchant"].(float64))
		userInfo.ClosingBalance = int(value.(map[string]interface{})["closingBalance"].(float64))
		userInfo.BlockingBalance = int(value.(map[string]interface{})["blockingBalance"].(float64))
		userInfo.ProductType = int(value.(map[string]interface{})["productType"].(float64))
		userInfo.IsPrimary = value.(map[string]interface{})["isPrimary"].(bool)
		userInfos = append(userInfos, &userInfo)
	}

	return userInfos, nil
}
