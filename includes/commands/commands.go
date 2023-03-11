package commands

import "gopkg.in/telebot.v3"

type Mode int

const (
	TextMode Mode = iota
	ImageMode
)

var CurrentMode = TextMode

var AllCommands []telebot.Command
