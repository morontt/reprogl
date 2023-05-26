package telegram

import (
	"bytes"
	"encoding/json"
	"net/http"
	"xelbot.com/reprogl/api"
	"xelbot.com/reprogl/container"
)

type message struct {
	Chat           int    `json:"chat_id"`
	Text           string `json:"text"`
	ParseMode      string `json:"parse_mode"`
	DisablePreview bool   `json:"disable_web_page_preview"`
}

var telegramAdminId int
var telegramToken string

func init() {
	cnf := container.GetConfig()
	telegramAdminId = cnf.TelegramAdminID
	telegramToken = cnf.TelegramToken
}

func SendNotification(app *container.Application, text string) {
	jsonBody, err := json.Marshal(createMessage(text))
	if err != nil {
		app.ErrorLog.Printf("telegram notification: %s\n", err.Error())
		return
	}

	request, err := http.NewRequest(
		http.MethodPost,
		"https://api.telegram.org/bot"+telegramToken+"/sendMessage",
		bytes.NewReader(jsonBody))
	if err != nil {
		app.ErrorLog.Printf("telegram notification: %s\n", err.Error())
		return
	}

	request.Header.Set("Content-Type", "application/json")

	_, err = api.Send(request)
	if err != nil {
		app.ErrorLog.Printf("telegram notification: %s\n", err.Error())
		return
	}
}

func createMessage(text string) message {
	return message{
		Chat:           telegramAdminId,
		Text:           text,
		ParseMode:      "MarkdownV2",
		DisablePreview: true,
	}
}
