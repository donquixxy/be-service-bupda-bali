package invelirepository

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type LoginRepositoryInterface interface {
	Login(username string, password string) (string, error)
}

type LoginRepositoryImplementation struct {
}

func NewLoginRepository() LoginRepositoryInterface {
	return &LoginRepositoryImplementation{}
}

func (service *LoginRepositoryImplementation) Login(username string, password string) (string, error) {
	body, _ := json.Marshal(map[string]interface{}{
		"username": username,
		"password": password,
		"imei":     "",
	})

	reqBody := io.NopCloser(strings.NewReader(string(body)))

	// URL
	url, _ := url.Parse("http://api-dev.cardlez.com:8089/login/member")

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
		log.Fatalf("An Error Occured %v", err)
	}

	defer resp.Body.Close()

	// Read response body
	data, _ := io.ReadAll(resp.Body)
	fmt.Printf("body: %s\n", data)

	return "", nil
}
