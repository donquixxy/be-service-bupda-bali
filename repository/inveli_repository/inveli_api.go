package invelirepository

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/machinebox/graphql"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/inveli"
)

type InveliAPIRepositoryInterface interface {
	InveliResgisration(inveliRegistrationModel *inveli.InveliRegistrationModel) error
	InveliLogin(username, password string) *inveli.InveliLoginModel
	InveliUbahPin(id, pin string) error
	InveliUpdateMember() error
}

type InveliAPIRepositoryImplementation struct {
}

func NewInveliAPIRepository() InveliAPIRepositoryInterface {
	return &InveliAPIRepositoryImplementation{}
}

func (r *InveliAPIRepositoryImplementation) InveliUpdateMember() error {
	graphql.NewClient("http://api-dev.cardlez.com:8089/query")

	return nil
}

func (r *InveliAPIRepositoryImplementation) InveliUbahPin(id, pin string) error {
	client := graphql.NewClient("http://api-dev.cardlez.com:8089/query")

	// make a request
	req := graphql.NewRequest(`
		mutation ($object: MemberInput!) {
			updateMember(member: $object)
			{
				id
			}
		}
	`)

	req.Var("memberObject", map[string]interface{}{
		"id":       id,
		"loginPin": pin,
	})

	// run it and capture the response
	ctx := context.Background()
	var respData interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Println(err)
		return err
	}

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

	reqDump, _ := httputil.DumpRequestOut(req, true)
	fmt.Printf("REQUEST:\n%s", string(reqDump))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("An Error Occured %v", err)
	}

	// Read response body
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("body: %s\n", data)

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
		log.Println(err)
		return err
	}

	return nil
}
