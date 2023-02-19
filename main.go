package main

import (
	"context"
	"github.com/PullRequestInc/go-gpt3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	gCfg Config = Config{
		// default config values here
		Temperature: 0.7,
		MaxTokens:   512,
		HistorySize: 1024,
	}
	endSequence = "@END"
)

// Config is an application configuration structure
type Config struct {
	OpenAiKey   string  `yaml:"open-ai-key"`
	BotApiKey   string  `yaml:"bot-api-key"`
	Temperature float32 `yaml:"temperature"`
	MaxTokens   int     `yaml:"max-tokens"`
	HistorySize int     `yaml:"history-size"`
}

func InitConfig() {
	// load file
	err := cleanenv.ReadConfig("config.yaml", &gCfg)
	if err != nil {
		log.Println(err)
	}
	log.Printf("CONFIG: %+v", gCfg)
}

func saveToHistory(userName string, userMessage string, aiResponse string) {
	hist := getHistory(userName)

	file, err := os.OpenFile(userName+".txt", os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	newHist := hist + userMessage + aiResponse
	start := utf8.RuneCountInString(newHist) - gCfg.HistorySize
	if start < 0 {
		start = 0
	}
	newHist = string([]rune(newHist)[start:]) // convert to rune is important, because we don't want to break utf-8
	_, _ = file.WriteString(newHist)
}

func getHistory(userName string) string {
	b, err := os.ReadFile(userName + ".txt")
	if err != nil {
		log.Println(err)
		return ""
	}

	str := string(b) // convert content to a 'string'
	return str
}

func clearHistory(userName string) {
	file, err := os.OpenFile(userName+".txt", os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0600)
	defer file.Close()
	if err != nil {
		log.Println(err)
	}
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

		if update.Message != nil && update.Message.IsCommand() { // commands
			// Create a new MessageConfig. We don't have text yet,
			// so we leave it empty.
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			userName := strconv.Itoa(int(update.Message.From.ID)) + "-" + update.Message.From.UserName

			// Extract the command from the Message.
			switch update.Message.Command() {
			case "help":
				msg.Text = "I understand commands /start and /hist."
			case "start":
				clearHistory(userName)
				msg.Text = "History cleared"
			case "hist":
				hist := getHistory(userName)
				if hist == "" {
					msg.Text = "_(empty)_"
					msg.ParseMode = "markdown"
				} else {
					msg.Text = hist
					msg.Text += "\n======\n"
					msg.Text += "Length in runes: " + strconv.Itoa(utf8.RuneCountInString(hist)) + "\n"
					msg.Text += "History size limit: " + strconv.Itoa(gCfg.HistorySize) + "\n"
				}
			default:
				msg.Text = "I don't know that command"
			}

			if _, err := bot.Send(msg); err != nil {
				log.Println("error:", err)
			}
			continue
		}

		if update.Message != nil { // If we got normal a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			typingMsg := tgbotapi.NewChatAction(update.Message.Chat.ID, "typing")
			bot.Send(typingMsg)

			userName := strconv.Itoa(int(update.Message.From.ID)) + "-" + update.Message.From.UserName

			history := getHistory(userName)
			humanPart := "Human: " + update.Message.Text + endSequence + "\n"

			log.Println("SENDING TO API:", history+humanPart)
			resp, err := client.Completion(ctx, gpt3.CompletionRequest{
				Prompt:      []string{history + humanPart + "AI: "},
				MaxTokens:   gpt3.IntPtr(gCfg.MaxTokens),
				Temperature: gpt3.Float32Ptr(gCfg.Temperature),
				Stop:        []string{endSequence},
			})

			response := ""
			if err != nil {
				log.Println("Error:", err)
				response = err.Error()
			} else {
				log.Println("GPT RESPONSE:", strings.Trim(resp.Choices[0].Text, "\n"), err)
				response = strings.Trim(resp.Choices[0].Text, "\n")
			}

			aiPart := "AI: " + response + endSequence + "\n"
			saveToHistory(userName, humanPart, aiPart)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}
