package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sync"
	"unicode/utf8"

	"xelbot.com/reprogl/api"
	"xelbot.com/reprogl/api/backend"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
)

type message struct {
	Chat           int    `json:"chat_id"`
	Text           string `json:"text"`
	ParseMode      string `json:"parse_mode"`
	DisablePreview bool   `json:"disable_web_page_preview"`
}

var telegramLocker sync.Mutex

func SendNotification(
	app *container.Application,
	comment *backend.CreatedCommentDTO,
	article *models.ArticleForComment,
) {
	text := generateText(comment, article)
	app.InfoLog.Printf("Telegram notification text:\n%s", text)

	jsonBody, err := json.Marshal(createMessage(text))
	if err != nil {
		app.LogError(err)
		return
	}

	request, err := http.NewRequest(
		http.MethodPost,
		"https://api.telegram.org/bot"+container.GetConfig().TelegramToken+"/sendMessage",
		bytes.NewReader(jsonBody))
	if err != nil {
		app.LogError(err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	resp, err := send(request)
	if err != nil {
		app.LogError(err)
		return
	}

	defer resp.Body.Close()
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		app.LogError(err)
		return
	}

	app.InfoLog.Printf("Telegram answer:\nStatus: %s\n\n%s", resp.Status, string(buf))
}

func createMessage(text string) message {
	return message{
		Chat:           container.GetConfig().TelegramAdminID,
		Text:           text,
		ParseMode:      "MarkdownV2",
		DisablePreview: true,
	}
}

func generateText(
	comment *backend.CreatedCommentDTO,
	article *models.ArticleForComment,
) (msg string) {
	msg = fmt.Sprintf(
		"Кто\\-то оставил [комментарий](%s)\n\n*ID*: %d\n",
		container.GenerateAbsoluteURL("article", "slug", article.Slug),
		comment.ID,
	)
	if len(comment.Name) > 0 {
		msg += "*Name*: " + escapeMarkdownCharacters(comment.Name)
	}
	if len(comment.Country) > 0 {
		msg += " " + comment.Country + "\n"
	} else {
		msg += "\n"
	}
	if len(comment.Email) > 0 {
		msg += "*Email*: " + escapeMarkdownCharacters(comment.Email) + "\n"
	}
	if len(comment.Website) > 0 {
		msg += "*Website*: " + escapeMarkdownCharacters(comment.Website) + "\n"
	}

	msg += "\n" + escapeMarkdownCharacters(stripTags(comment.Text))

	return
}

func stripTags(content string) string {
	re := regexp.MustCompile(`<(.|\n)*?>`)

	return re.ReplaceAllString(content, "")
}

func escapeMarkdownCharacters(content string) string {
	buffer := make([]rune, 0, 2*utf8.RuneCountInString(content))
	for _, e := range []rune(content) {
		switch e {
		case '_', '*', '[', ']', '(', ')', '~', '`', '>', '#', '+', '-', '=', '|', '{', '}', '.', '!':
			buffer = append(buffer, '\\', e)
		default:
			buffer = append(buffer, e)
		}
	}

	return string(buffer)
}

func send(req *http.Request) (*http.Response, error) {
	telegramLocker.Lock()
	defer telegramLocker.Unlock()

	return api.Send(req)
}
