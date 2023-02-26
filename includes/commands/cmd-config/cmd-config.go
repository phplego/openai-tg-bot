package cmd_config

import (
	tele "gopkg.in/telebot.v3"
	"gopkg.in/yaml.v3"
	"openai-tg-bot/includes/config"
)

func Handler(c tele.Context) error {
	resp := ""
	resp += "*CONFIG*\n"
	var cfgCopy = config.Data
	cfgCopy.OpenAiKey = ""
	cfgCopy.BotApiKey = ""
	bytes, _ := yaml.Marshal(cfgCopy)
	resp += "```\n"
	resp += string(bytes)
	resp += "```\n"
	return c.Send(resp, &tele.SendOptions{
		ParseMode: "markdown",
	})
}
