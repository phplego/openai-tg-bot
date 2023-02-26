package cmd_help

import (
	tele "gopkg.in/telebot.v3"
)

func Handler(c tele.Context) error {
	return c.Send("I understand commands /start and /hist.")
}
