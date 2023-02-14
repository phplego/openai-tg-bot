package main

import (
	"context"
	"github.com/PullRequestInc/go-gpt3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"strings"
)

var (
	gCfg Config = Config{
		// default config values here
		Temperature: 0.7,
		MaxTokens:   512,
	}
)

// Config is an application configuration structure
type Config struct {
	OpenAiKey   string  `yaml:"open-ai-key"`
	BotApiKey   string  `yaml:"bot-api-key"`
	Temperature float32 `yaml:"temperature"`
	MaxTokens   int     `yaml:"max-tokens"`
}

func InitConfig() {
	// load file
	err := cleanenv.ReadConfig("config.yaml", &gCfg)
	if err != nil {
		log.Println(err)
	}
	log.Printf("CONFIG: %+v", gCfg)
}

func main() {
	InitConfig()
	client := gpt3.NewClient(gCfg.OpenAiKey, gpt3.WithDefaultEngine(gpt3.TextDavinci003Engine))
	bot, err := tgbotapi.NewBotAPI(gCfg.BotApiKey)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10
	updates := bot.GetUpdatesChan(u)

	ctx := context.Background()
	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			typingMsg := tgbotapi.NewChatAction(update.Message.Chat.ID, "typing")
			bot.Send(typingMsg)

			resp, err := client.Completion(ctx, gpt3.CompletionRequest{
				Prompt:      []string{update.Message.Text},
				MaxTokens:   gpt3.IntPtr(gCfg.MaxTokens),
				Temperature: gpt3.Float32Ptr(gCfg.Temperature),
			})

			response := ""
			if err != nil {
				log.Println("Error:", err)
				response = err.Error()
			} else {
				log.Println("GPT RESPONSE:", strings.Trim(resp.Choices[0].Text, "\n"), err)
				response = strings.Trim(resp.Choices[0].Text, "\n")
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}
