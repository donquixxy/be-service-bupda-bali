package invelirepository

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/tensuqiuwulu/be-service-bupda-bali/model/inveli"
)

type InveliRegistarationRepositoryInterface interface {
	InveliResgisration(inveliRegistrationModel *inveli.InveliRegistrationModel) error
}

type InveliRegistarationRepositoryImplementation struct {
}

func NewInveliRegistarationRepository() InveliRegistarationRepositoryInterface {
	return &InveliRegistarationRepositoryImplementation{}
}

func (service *InveliRegistarationRepositoryImplementation) InveliResgisration(inveliRegistrationModel *inveli.InveliRegistrationModel) error {

	str := fmt.Sprintf(`{"query":"mutation {createMember(member: {handphone: %s,memName: %s,birthPlace: "",
						birthDate: \"\",
						emailAddress: \"%s\",
						identityNumber: \"%s\",
						identityType: \"1\",
						noNPWP: \"\",
						gender: \"\",
						address: \"\",
						alamatRumah2: \"\",
						alamatRumahKerja1: \"\",
						alamatRumahKerja2: \"\",
						bankCode: \"014\",
						noInduk: \"\",
						fileNameFotoKtp: \"\",
						fileNameFotoSelfie: \"\",
						bankAccountName: \"\",
						bankAccountNo: \"\",
						bankAccountType: \"\",
						grade: \"bronze\",
						maritalStatus: \"\",
						recordStatus: \"\",
						nationality: \"\",
						memType: \"\",
						isLocked: \"\",
						virtualAccountNo: \"\",
						goldStatus: \"\",
						bankID: \"\",
					}
				)"}
	`, inveliRegistrationModel.Phone, inveliRegistrationModel.MemberName, inveliRegistrationModel.Email, inveliRegistrationModel.NIK)

	fmt.Println(str)

	body := strings.NewReader(str)

	fmt.Println(str)
	req, err := http.NewRequest("POST", "http://api-dev.cardlez.com:8089/register", body)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	log.Println(resp)

	return nil
}
