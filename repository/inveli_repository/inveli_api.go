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
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/inveli"
)

type InveliAPIRepositoryInterface interface {
	InveliResgisration(inveliRegistrationModel *inveli.InveliRegistrationModel) error
	InveliLogin(username, password string) *inveli.InveliLoginModel
	InveliUbahPassword(id, password, token string) (interface{}, error)
	InveliUbahPasswordUserExisting(id, password, token string) error
	InveliUpdateMember(user *entity.User, userProfile *entity.UserProfile, accessToken string, groupIdBupda string) error
	GetAccountInfo(IDMember, token string) ([]inveli.InveliAcountInfo, error)
	InveliCreatePaylater(token string, IDMember string, AccountID string, Amount float64, totalAmount float64, isMerchant float64, bunga float64, loanProductId string) error
	GetStatusAccount(IDMember, token string) (bool, error)
	GetBalanceAccount(Code, token string) (*inveli.InveliAcountInfo, error)
	GetKodeBIN(token string) (string, error)
	InquiryVaNasabah(phone, token string) (*inveli.InquiryVaNasabah, error)
	ApiPayment(creditAccount, debitAccount, token string, amount float64, isMerchant float64) error
	GetTunggakan(LoanID, token string) ([]inveli.TunggakanPaylater, error)
	GetLimitPayLater(IDMember, token string) (*inveli.LimitPaylater, error)
	GetTagihanPaylater(IDMember, token string) ([]inveli.RiwayatPinjaman2, error)
	PayPaylater(loanID, token string) error
	GetLoanProduct(token string) (float64, error)
	GetLoanProductId(token string) (string, error)
	GetSaldoBupda(token, groupID string) (float64, error)
}

type InveliAPIRepositoryImplementation struct {
}

func NewInveliAPIRepository() InveliAPIRepositoryInterface {
	return &InveliAPIRepositoryImplementation{}
}

func (r *InveliAPIRepositoryImplementation) GetSaldoBupda(token, groupID string) (float64, error) {
	client := graphql.NewClient(config.GetConfig().Inveli.InveliAPI + "/query")

	req := graphql.NewRequest(`
  	query ($groupID: String!) {
			accountByGroupID(groupID: $groupID){
				id
				code
				accountName
				recordStatus
				balance
				isPrimary
			}
  	}
	`)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Var("groupID", groupID)
	ctx := context.Background()
	var respData interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Println(err)
		return 0, err
	}

	bupdaSaldo := respData.(map[string]interface{})["accountByGroupID"].(map[string]interface{})["balance"].(float64)

	return bupdaSaldo, nil
}

func (r *InveliAPIRepositoryImplementation) InveliUbahPasswordUserExisting(id, password, token string) error {
	client := graphql.NewClient(config.GetConfig().Inveli.InveliAPI + "/query")
	req := graphql.NewRequest(`
		mutation ($id: String!, $password: String!) {
			changePassword(
				id: $id, 
				newPassword: $password)
		}
	`)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Var("id", id)
	req.Var("password", password)
	ctx := context.Background()

	// fmt.Println("request : ", req)

	var respData interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Println("error change password : ", err)
		return err
	}

	fmt.Println("response ubah password : ", respData)

	return nil
}

func (r *InveliAPIRepositoryImplementation) PayPaylater(IDMember, token string) error {
	client := graphql.NewClient(config.GetConfig().Inveli.InveliAPI + "/query")

	req := graphql.NewRequest(`
		mutation ($id: String!) {
			autodebetLoan(
				customerID: $id
			)
		}
	`)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Var("id", IDMember)
	ctx := context.Background()
	var respData interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *InveliAPIRepositoryImplementation) GetLoanProductId(token string) (string, error) {
	client := graphql.NewClient(config.GetConfig().Inveli.InveliAPI + "/query")

	req := graphql.NewRequest(`
	 {
			loanProducts{
        loanProductID
        loanProductNo
        productName
        interestCalcType
        interestDayBasisID
        interestPercentage
        calcBase
        schemeType
        scheduleType
        collectionType
        roundingType
        penaltyTerminationAmount
        minPenaltyTermination
        penaltyAmount
        minLoanAmount
        recordStatus
        companyBranchID
        loanProductCharges{
          loanProductChargeID
          loanProductChargeNo
          parentCOAID
          description
          loanProductID
          amount
          amountPercentage
          cOANo
          recordStatus
        }
    	}
		}
	`)

	req.Header.Set("Authorization", "Bearer "+token)
	ctx := context.Background()
	var respData interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Println(err)
		return "", err
	}

	var loandProductID string
	// var data []interface{}
	for _, v := range respData.(map[string]interface{})["loanProducts"].([]interface{}) {
		if v.(map[string]interface{})["loanProductNo"].(string) == "LPDCO160012211001" {
			loandProductID = v.(map[string]interface{})["loanProductID"].(string)
		}
	}

	fmt.Println("loan product id ", loandProductID)
	return loandProductID, nil
}

func (r *InveliAPIRepositoryImplementation) GetLoanProduct(token string) (float64, error) {
	client := graphql.NewClient(config.GetConfig().Inveli.InveliAPI + "/query")

	req := graphql.NewRequest(`
	 {
			loanProducts{
        loanProductID
        loanProductNo
        productName
        interestCalcType
        interestDayBasisID
        interestPercentage
        calcBase
        schemeType
        scheduleType
        collectionType
        roundingType
        penaltyTerminationAmount
        minPenaltyTermination
        penaltyAmount
        minLoanAmount
        recordStatus
        companyBranchID
        loanProductCharges{
          loanProductChargeID
          loanProductChargeNo
          parentCOAID
          description
          loanProductID
          amount
          amountPercentage
          cOANo
          recordStatus
        }
    	}
		}
	`)

	req.Header.Set("Authorization", "Bearer "+token)
	ctx := context.Background()
	var respData interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Println(err)
		return 0, err
	}

	var bunga float64
	// var data []interface{}
	for _, v := range respData.(map[string]interface{})["loanProducts"].([]interface{}) {
		// dev loan no LPDCO160012210001
		// prod loan no LPDCO160012211001
		if v.(map[string]interface{})["loanProductNo"].(string) == "LPDCO160012211001" {
			bunga = v.(map[string]interface{})["interestPercentage"].(float64)
		}
	}

	fmt.Println("riwayat pinjaman ", bunga)
	return bunga, nil
}

func (r *InveliAPIRepositoryImplementation) GetTagihanPaylater(IDMember, token string) ([]inveli.RiwayatPinjaman2, error) {
	client := graphql.NewClient(config.GetConfig().Inveli.InveliAPI + "/query")

	req := graphql.NewRequest(`
		query ($id: String!) {
		loans(memberID: $id){
      loanID
      code
      customerID
      customerName
      productDesc
      loanProductID
      startDate
      endDate
      tenorMonth
      loanAmount
      interestPercentage
      repaymentMethod
      accountID
      userInsert
      dateInsert
      dateAuthor
      userAuthor
      recordStatus
      isLiquidated
      outstandingAmount
      nominalWajib
      filePDFName
      loanAccountRepayments{
      	id
        loanAccountID
        repaymentType
        repaymentDate
        repaymentInterest
        repaymentPrincipal
        repaymentAmount
        repaymentInterestPaid
        repaymentPrincipalPaid
        outStandingBakiDebet
        tellerId
        isPaid
        amountPaid
        paymentTxnID
        recordStatus
        userInsert
        dateInsert
        userUpdate
        dateUpdate
      	}
    	}
		}
	`)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Var("id", IDMember)
	ctx := context.Background()
	var respData interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Println(err)
		return nil, err
	}

	riwayatPinjamans := []inveli.RiwayatPinjaman2{}
	// var data []interface{}
	for _, loan := range respData.(map[string]interface{})["loans"].([]interface{}) {
		riwayatPinjaman := inveli.RiwayatPinjaman2{}
		riwayatPinjaman.ID = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["id"].(string)
		riwayatPinjaman.LoanAccountID = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["loanAccountID"].(string)
		riwayatPinjaman.RepaymentDate = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["repaymentDate"].(string)
		riwayatPinjaman.RepaymentInterest = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["repaymentInterest"].(float64)
		riwayatPinjaman.RepaymentPrincipal = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["repaymentPrincipal"].(float64)
		riwayatPinjaman.RepaymentAmount = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["repaymentAmount"].(float64)
		riwayatPinjaman.RepaymentInterestPaid = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["repaymentInterestPaid"].(float64)
		riwayatPinjaman.RepaymentPrincipalPaid = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["repaymentPrincipalPaid"].(float64)
		riwayatPinjaman.OutStandingBakiDebet = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["outStandingBakiDebet"].(float64)
		riwayatPinjaman.TellerID = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["tellerId"].(string)
		riwayatPinjaman.IsPaid = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["isPaid"].(bool)
		riwayatPinjaman.AmountPaid = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["amountPaid"].(float64)
		riwayatPinjaman.PaymentTxnID = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["paymentTxnID"].(string)
		riwayatPinjaman.UserInsert = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["userInsert"].(string)
		riwayatPinjaman.DateInsert = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["dateInsert"].(string)
		riwayatPinjaman.UserUpdate = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["userUpdate"].(string)
		riwayatPinjaman.DateUpdate = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["dateUpdate"].(string)
		riwayatPinjamans = append(riwayatPinjamans, riwayatPinjaman)
	}

	fmt.Println("riwayat pinjaman ", riwayatPinjamans)
	return riwayatPinjamans, nil
}

func (r *InveliAPIRepositoryImplementation) GetLimitPayLater(IDMember, token string) (*inveli.LimitPaylater, error) {
	client := graphql.NewClient(config.GetConfig().Inveli.InveliAPI + "/query")

	req := graphql.NewRequest(`
		query ($id: String!) {
			limitByCustomerID(customerID: $id) {
	    	customerID
        nominal
        nominalAvailable
			}
		}
	`)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Var("id", IDMember)
	ctx := context.Background()
	var respData interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Println(err)
		return nil, err
	}

	limitPayLater := &inveli.LimitPaylater{}

	limitPayLater.MaxLimit = respData.(map[string]interface{})["limitByCustomerID"].(map[string]interface{})["nominal"].(float64)
	limitPayLater.AvailableLimit = respData.(map[string]interface{})["limitByCustomerID"].(map[string]interface{})["nominalAvailable"].(float64)

	return limitPayLater, nil
}

func (r *InveliAPIRepositoryImplementation) GetTunggakan(LoanID, token string) ([]inveli.TunggakanPaylater, error) {

	client := graphql.NewClient(config.GetConfig().Inveli.InveliAPI + "/query")

	req := graphql.NewRequest(`
		 query ($id: String!) {
			getLoanPassdue(loanID: $id) {
    		loanPassdueID               
    		loanPassdueNo               
    		loanAccountRepaymentID      
    		loanID                      
    		overduePrincipal            
    		overdueInterest             
    		overduePenalty              
    		overdueAmount               
    		isPaid                      
    		isWaivePenalty              
    		userInsert                  
    		dateInsert                  
    		userUpdate                  
    		dateUpdate                  
    		passdueCode                 
    	}
		}
	`)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Var("id", LoanID)
	ctx := context.Background()
	var respData interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Println(err)
		return nil, err
	}

	tunggakanResp := []inveli.TunggakanPaylater{}
	for _, v := range respData.(map[string]interface{})["getLoanPassdue"].([]interface{}) {
		tunggakanResp = append(tunggakanResp, inveli.TunggakanPaylater{
			LoanPassdueID:          v.(map[string]interface{})["loanPassdueID"].(string),
			LoanPassdueNo:          v.(map[string]interface{})["loanPassdueNo"].(string),
			LoanAccountRepaymentID: v.(map[string]interface{})["loanAccountRepaymentID"].(string),
			LoanID:                 v.(map[string]interface{})["loanID"].(string),
			OverduePrincipal:       v.(map[string]interface{})["overduePrincipal"].(float64),
			OverdueInterest:        v.(map[string]interface{})["overdueInterest"].(float64),
			OverduePenalty:         v.(map[string]interface{})["overduePenalty"].(float64),
			OverdueAmount:          v.(map[string]interface{})["overdueAmount"].(float64),
			IsPaid:                 v.(map[string]interface{})["isPaid"].(bool),
			IsWaivePenalty:         v.(map[string]interface{})["isWaivePenalty"].(bool),
			UserInsert:             v.(map[string]interface{})["userInsert"].(string),
			DateInsert:             v.(map[string]interface{})["dateInsert"].(string),
			UserUpdate:             v.(map[string]interface{})["userUpdate"].(string),
			DateUpdate:             v.(map[string]interface{})["dateUpdate"].(string),
			PassdueCode:            v.(map[string]interface{})["passdueCode"].(string),
		})
	}

	log.Println("tunggakan : ", tunggakanResp)

	return tunggakanResp, nil
}

func (r *InveliAPIRepositoryImplementation) ApiPayment(creditAccount, debitAccount, token string, amount float64, isMerchant float64) error {
	client := graphql.NewClient(config.GetConfig().Inveli.InveliAPI + "/query")

	log.Println("credit account api payment: ", creditAccount)

	req := graphql.NewRequest(`
		mutation ($payment: MemberPaymentInput!) {
			requestPayment(memberPayment: $payment)
			{
				id
        code
        debitAccountID
        currencyID
        amount
        debitDate
        transferType
        creditAccountID
        creditDate
        note
        postransactionRefNo
        posName
        postransactionDate
        revNo
        recordStatus
			}
		}
	`)

	log.Println("creditAccount : ", creditAccount)
	log.Println("debitAccount : ", debitAccount)
	log.Println("amount : ", amount)
	log.Println("isMerchant : ", isMerchant)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Var("payment", map[string]interface{}{
		"creditAccount":       creditAccount,
		"transferType":        "1",
		"debitAccount":        debitAccount,
		"amount":              amount,
		"note":                "",
		"paymentInterest":     isMerchant,
		"memberPaymentDetail": []string{},
	})

	fmt.Println("request api payment : ", req)

	ctx := context.Background()

	var respData interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Println("error api payment : ", err)
		return err
	}

	fmt.Println("response api payment : ", respData)

	return nil
}

func (r *InveliAPIRepositoryImplementation) GetKodeBIN(token string) (string, error) {

	client := graphql.NewClient(config.GetConfig().Inveli.InveliAPI + "/query")

	req := graphql.NewRequest(`
	{
		getApplicationConfigByKey(configKey: "BINCodeByHp"){
			configKey
        	configValue
			}
		}
	`)

	req.Header.Set("Authorization", "Bearer "+token)
	ctx := context.Background()
	var respData interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Println(err)
		return "", err
	}

	kodeBin := respData.(map[string]interface{})["getApplicationConfigByKey"].(map[string]interface{})["configValue"].(string)

	log.Println("kode bin ", kodeBin)

	return kodeBin, nil
}

func (r *InveliAPIRepositoryImplementation) InquiryVaNasabah(phone, token string) (*inveli.InquiryVaNasabah, error) {
	client := graphql.NewClient(config.GetConfig().Inveli.InveliAPI + "/query")

	req := graphql.NewRequest(` 
		query ($phone: String!) {
			accountByHandPhone(handPhone: $phone) {
				id
				code
				accountName
				accountName2
				memberID
			}
		}
	`)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Var("phone", phone)
	ctx := context.Background()
	var respData interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Println(err)
		return nil, err
	}

	inquiryVaResult := &inveli.InquiryVaNasabah{}

	inquiryVaResult.ID = respData.(map[string]interface{})["accountByHandPhone"].(map[string]interface{})["id"].(string)
	inquiryVaResult.Code = respData.(map[string]interface{})["accountByHandPhone"].(map[string]interface{})["code"].(string)
	inquiryVaResult.AccountName = respData.(map[string]interface{})["accountByHandPhone"].(map[string]interface{})["accountName"].(string)
	inquiryVaResult.AccountName2 = respData.(map[string]interface{})["accountByHandPhone"].(map[string]interface{})["accountName2"].(string)
	inquiryVaResult.MemberID = respData.(map[string]interface{})["accountByHandPhone"].(map[string]interface{})["memberID"].(string)

	return inquiryVaResult, nil
}

func (r *InveliAPIRepositoryImplementation) InveliCreatePaylater(token string, IDMember string, AccountID string, Amount float64, totalAmount float64, isMerchant float64, bunga float64, loanProductId string) error {

	client := graphql.NewClient(config.GetConfig().Inveli.InveliAPI + "/query")

	log.Println("loan id di inveli", loanProductId)

	req := graphql.NewRequest(`
		mutation ($loanInputParam: LoanInput!) {
			createLoan(loanInputParam: $loanInputParam) 
		}
	`)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Var("loanInputParam", map[string]interface{}{
		"accountID":             AccountID,
		"memberID":              IDMember,
		"loanProductID":         loanProductId,
		"tenor":                 1,
		"loanAmount":            Amount,
		"interestPercent":       bunga,
		"isAutoApprove":         true,
		"paymentInterest":       isMerchant,
		"totalAmount":           totalAmount,
		"memberLoanAttachments": []string{},
	})

	fmt.Println("request create pinjaman : ", req)

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
	client := graphql.NewClient(config.GetConfig().Inveli.InveliAPI + "/query")
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

func (r *InveliAPIRepositoryImplementation) InveliUpdateMember(user *entity.User, userProfile *entity.UserProfile, accessToken string, groupIdBupda string) error {
	client := graphql.NewClient(config.GetConfig().Inveli.InveliAPI + "/query")

	// make a request
	req := graphql.NewRequest(`
	  mutation ($member: MemberInput!, $memberDetail: MemberDetailInput!) {
			updateMember(member: $member, memberDetail: $memberDetail) {
				id
			}
	  }
	`)

	log.Println("groupId :", user.Phone)
	log.Println("groupId : ", groupIdBupda)

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
		"referralMemberID": groupIdBupda,
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

	urlString := config.GetConfig().Inveli.InveliAPI + "/login/member"
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
	client := graphql.NewClient(config.GetConfig().Inveli.InveliAPI + "/register")

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

func (r *InveliAPIRepositoryImplementation) GetAccountInfo(IDMember, token string) ([]inveli.InveliAcountInfo, error) {

	client := graphql.NewClient(config.GetConfig().Inveli.InveliAPI + "/query")

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
	// log.Println("req ", respData)

	userInfos := []inveli.InveliAcountInfo{}
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
		userInfo.Balance = float64(value.(map[string]interface{})["balance"].(float64))
		userInfo.BalanceMerchant = int(value.(map[string]interface{})["balanceMerchant"].(float64))
		userInfo.ClosingBalance = int(value.(map[string]interface{})["closingBalance"].(float64))
		userInfo.BlockingBalance = int(value.(map[string]interface{})["blockingBalance"].(float64))
		userInfo.ProductType = int(value.(map[string]interface{})["productType"].(float64))
		userInfo.IsPrimary = value.(map[string]interface{})["isPrimary"].(bool)
		userInfos = append(userInfos, userInfo)
	}

	// fmt.Println(userInfos)

	return userInfos, nil
}

func (r *InveliAPIRepositoryImplementation) GetStatusAccount(IDMember, token string) (bool, error) {
	client := graphql.NewClient(config.GetConfig().Inveli.InveliAPI + "/query")

	req := graphql.NewRequest(` 
		query ($id: String!) {	
			member (id: $id) {
        recordStatus
			}
		}
	`)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Var("id", IDMember)
	ctx := context.Background()
	var respData interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		return false, err
	}
	// log.Println("res ", respData)

	resStatus := respData.(map[string]interface{})["member"].(map[string]interface{})["recordStatus"].(float64)
	if resStatus == 2 {
		return true, nil
	} else {
		return false, nil
	}
}

func (r *InveliAPIRepositoryImplementation) GetBalanceAccount(Code, token string) (*inveli.InveliAcountInfo, error) {
	client := graphql.NewClient(config.GetConfig().Inveli.InveliAPI + "/query")

	req := graphql.NewRequest(` 
		query ($keyword: String!) {	
			accounts (search: CODE, keyword: $keyword) {
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
	req.Var("keyword", Code)
	ctx := context.Background()
	var respData interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("req ", respData)

	accountBalance := &inveli.InveliAcountInfo{
		ID:              respData.(map[string]interface{})["accounts"].([]interface{})[0].(map[string]interface{})["id"].(string),
		Code:            respData.(map[string]interface{})["accounts"].([]interface{})[0].(map[string]interface{})["code"].(string),
		AccountName2:    respData.(map[string]interface{})["accounts"].([]interface{})[0].(map[string]interface{})["accountName"].(string),
		ProductName:     respData.(map[string]interface{})["accounts"].([]interface{})[0].(map[string]interface{})["productName"].(string),
		Balance:         respData.(map[string]interface{})["accounts"].([]interface{})[0].(map[string]interface{})["balance"].(float64),
		BalanceMerchant: int(respData.(map[string]interface{})["accounts"].([]interface{})[0].(map[string]interface{})["balanceMerchant"].(float64)),
		ClosingBalance:  int(respData.(map[string]interface{})["accounts"].([]interface{})[0].(map[string]interface{})["closingBalance"].(float64)),
		BlockingBalance: int(respData.(map[string]interface{})["accounts"].([]interface{})[0].(map[string]interface{})["blockingBalance"].(float64)),
	}

	return accountBalance, nil
}
