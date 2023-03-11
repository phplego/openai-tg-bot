package cmd_start

import (
	tele "gopkg.in/telebot.v3"
	user_state "openai-tg-bot/includes/user-state"
)

func Handler(c tele.Context) error {
	var (
		user = c.Sender()
	)
	user_state.ClearHistory(user.ID)
	return c.Send("History cleared")
}
