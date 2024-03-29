package goutil

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type teleService struct {
	chatId int64
	bot    *tgbotapi.BotAPI
}

// TeleService is a interface all of function third-party Telegram
type TeleService interface {
	SendError(ctx context.Context, path string, line int, msg string) (err error)
}

// NewTele is a function open connection
func NewTele(
	token string,
	chatId int64) (tele TeleService, err error) {
	bot, err := tgbotapi.NewBotAPI(token)
	return &teleService{bot: bot, chatId: chatId}, err
}

func (t *teleService) SendError(ctx context.Context, path string, line int, msg string) (err error) {

	// validate bot tele before send request
	if t.bot == nil {
		return nil
	}

	// template chat to telegram
	var template = "<b>-===ERROR NOTIFICATION===-</b>\n\n" +
		"<b>RequestId:</b> --REQUESTID--\n" +
		"<b>Method:</b> --METHOD--\n" +
		"<b>Endpoint:</b> --ENDPOINT--\n" +
		"<b>Error Message:</b> --MESSAGE--\n" +
		"<b>Path:</b> --PATH--\n" +
		"<b>Line:</b> --LINE--"

	// get context value from context
	res := GetContext(ctx)

	// replate template with message value
	template = strings.Replace(template, "--MESSAGE--", msg, 1)
	template = strings.Replace(template, "--PATH--", path, 1)
	template = strings.Replace(template, "--LINE--", fmt.Sprint(line), 1)

	// if request id exist
	if res.RequestId != "" {
		// replace template with value request id
		template = strings.Replace(template, "--REQUESTID--", res.RequestId, 1)
	}

	// if method exist
	if res.Method != "" {
		// replace template with value method
		template = strings.Replace(template, "--METHOD--", res.Method, 1)
	}

	// if endpoint exist
	if res.Endpoint != "" {
		// replace template with value endpoint
		template = strings.Replace(template, "--ENDPOINT--", res.Endpoint, 1)
	}

	// config bot
	message := tgbotapi.NewMessage(t.chatId, template)
	message.ParseMode = "HTML"

	// send message
	_, err = t.bot.Send(message)
	return

}
