package cmd_hist

import (
	tele "gopkg.in/telebot.v3"
	"openai-tg-bot/includes/config"
	"openai-tg-bot/includes/history"
	"strconv"
	"unicode/utf8"
)

func Handler(c tele.Context) error {
	var (
		user = c.Sender()
	)
	userName := strconv.Itoa(int(user.ID)) + "-" + user.Username
	hist := history.GetHistory(userName)
	if hist == "" {
		return c.Send("_(empty)_", &tele.SendOptions{
			ParseMode: "markdown",
		})
	} else {
		resp := hist
		resp += "\n======\n"
		resp += "Length in runes: " + strconv.Itoa(utf8.RuneCountInString(hist)) + "\n"
		resp += "History size limit: " + strconv.Itoa(config.Data.HistorySize) + "\n"
		return c.Send(resp)
	}
}
