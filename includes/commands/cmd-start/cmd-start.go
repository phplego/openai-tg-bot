package cmd_start

import (
	tele "gopkg.in/telebot.v3"
	"openai-tg-bot/includes/history"
	"strconv"
)

func Handler(c tele.Context) error {
	var (
		user = c.Sender()
	)
	userName := strconv.Itoa(int(user.ID)) + "-" + user.Username
	history.ClearHistory(userName)
	return c.Send("History cleared")
}
