package main

import (
	"gopkg.in/telebot.v3"
	"log"
	"openai-tg-bot/includes/commands"
	"openai-tg-bot/includes/commands/cmd-config"
	"openai-tg-bot/includes/commands/cmd-help"
	"openai-tg-bot/includes/commands/cmd-hist"
	"openai-tg-bot/includes/commands/cmd-start"
	"openai-tg-bot/includes/commands/on-text"
	"openai-tg-bot/includes/config"
	"openai-tg-bot/includes/open-ai"
	user_state "openai-tg-bot/includes/user-state"
	"time"
)

func main() {
	// configure logs
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// init subpackages
	config.Init()
	user_state.Init(config.Data)
	open_ai.Init()

	// init bot
	theBot, err := telebot.NewBot(telebot.Settings{
		Token:  config.Data.BotApiKey,
		Poller: &telebot.LongPoller{Timeout: 60 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("Authorized on Telegram bot account @%s", theBot.Me.Username)

	commands.AllCommands = []telebot.Command{{
		Text:        "start",
		Description: "Start new conversation",
	}, {
		Text:        "config",
		Description: "Show config params",
	}, {
		Text:        "hist",
		Description: "Conversation history",
	}, {
		Text:        "image",
		Description: "Switch to Image Mode",
	}, {
		Text:        "text",
		Description: "Switch to Text Mode",
	}, {
		Text:        "help",
		Description: "Help and instructions",
	}}

	err = theBot.SetCommands(commands.AllCommands)

	if err != nil {
		log.Fatal(err)
		return
	}

	theBot.Handle("/start", cmd_start.Handler)
	theBot.Handle("/hist", cmd_hist.Handler)
	theBot.Handle("/config", cmd_config.Handler)
	theBot.Handle("/help", cmd_help.Handler)
	theBot.Handle("/image", func(context telebot.Context) error {
		user_state.Load(context.Sender().ID)
		state := user_state.FindById(context.Sender().ID)
		state.Mode = commands.ImageMode
		user_state.Save(context.Sender().ID, state)
		return context.Send("Switched to image generation mode")

	})
	theBot.Handle("/text", func(context telebot.Context) error {
		user_state.Load(context.Sender().ID)
		state := user_state.FindById(context.Sender().ID)
		state.Mode = commands.TextMode
		user_state.Save(context.Sender().ID, state)
		return context.Send("Switched to text completion mode")
	})

	// All the text messages that weren't captured by existing handlers.
	// The handling depends on 'mode' (text, image)
	theBot.Handle(telebot.OnText, on_text.Handler)

	theBot.Start()
}
