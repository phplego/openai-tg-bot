package cmd_help

import (
	tele "gopkg.in/telebot.v3"
	"openai-tg-bot/includes/commands"
	user_state "openai-tg-bot/includes/user-state"
	"strconv"
)

func Handler(context tele.Context) error {

	user_state.Load(context.Sender().ID)
	state := user_state.FindById(context.Sender().ID)

	result := "Supported commands are:\n\n"

	for i, command := range commands.AllCommands {
		result += strconv.Itoa(i+1) + ". /" + command.Text
		result += " - " + command.Description + "\n"
	}

	result += "\n"
	result += "Current mode: "
	result += strconv.Itoa(int(state.Mode))
	return context.Send(result, &tele.SendOptions{
		ParseMode: "markdown",
	})
}
