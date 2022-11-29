package repository

import (
	"log"
	"net/http"
	"net/url"
)

type TelegramRepositoryInterface interface {
	SendMessageToTelegram(message, chatId, token string)
}

type TelegramRepositoryImplementation struct{}

func NewTelegramRepository() TelegramRepositoryInterface {
	return nil
}

func (repository *TelegramRepositoryImplementation) SendMessageToTelegram(message, numberOrder, chatId, token string) {
	url, _ := url.Parse("https://api.telegram.org/bot" + token + "/sendMessage?chat_id=" + chatId + "&text=" + message + " " + numberOrder + "")

	req := &http.Request{
		Method: "POST",
		URL:    url,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println("an error occured %v", err)
	}
	defer resp.Body.Close()
}
