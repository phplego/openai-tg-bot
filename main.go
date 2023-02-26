package main

import (
	tele "gopkg.in/telebot.v3"
	"log"
	"openai-tg-bot/includes/commands/cmd-config"
	"openai-tg-bot/includes/commands/cmd-help"
	"openai-tg-bot/includes/commands/cmd-hist"
	"openai-tg-bot/includes/commands/cmd-start"
	"openai-tg-bot/includes/commands/on-text"
	"openai-tg-bot/includes/config"
	"openai-tg-bot/includes/history"
	"openai-tg-bot/includes/open-ai"
	"time"
)

func main() {
	config.Init()
	history.Init(config.Data)
	open_ai.Init()

	pref := tele.Settings{
		Token:  config.Data.BotApiKey,
		Poller: &tele.LongPoller{Timeout: 60 * time.Second},
	}

	theBot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("Authorized on Telegram bot account @%s", theBot.Me.Username)

	err = theBot.SetCommands([]tele.Command{{
		Text:        "start",
		Description: "Start new conversation",
	}, {
		Text:        "config",
		Description: "Show config params",
	}, {
		Text:        "hist",
		Description: "Conversation history",
	}, {
		Text:        "help",
		Description: "Help and instructions",
	}})

	if err != nil {
		log.Fatal(err)
		return
	}

	theBot.Handle("/start", cmd_start.Handler)
	theBot.Handle("/hist", cmd_hist.Handler)
	theBot.Handle("/config", cmd_config.Handler)
	theBot.Handle("/help", cmd_help.Handler)

	// All the text messages that weren't captured by existing handlers.
	theBot.Handle(tele.OnText, on_text.Handler)

	theBot.Start()
}
