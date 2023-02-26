package main

import (
	"context"
	"github.com/ilyakaznacheev/cleanenv"
	gogpt "github.com/sashabaranov/go-gpt3"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	gCfg Config = Config{
		// default config values here
		Temperature: 0.7,
		MaxTokens:   512,
		HistorySize: 1024,
	}
	endSequence = " ␃"
)

// Config is an application configuration structure
type Config struct {
	OpenAiKey   string  `yaml:"open-ai-key,omitempty"`
	BotApiKey   string  `yaml:"bot-api-key,omitempty"`
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

	openAiClient := gogpt.NewClient(gCfg.OpenAiKey)

	pref := tele.Settings{
		Token:  gCfg.BotApiKey,
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
		Text:        "help",
		Description: "Help and instructions",
	}, {
		Text:        "config",
		Description: "Show config params",
	}, {
		Text:        "hist",
		Description: "Conversation history",
	}})
	if err != nil {
		log.Fatal(err)
		return
	}

	theBot.Handle("/help", func(c tele.Context) error {
		return c.Send("I understand commands /start and /hist.")
	})

	theBot.Handle("/hist", func(c tele.Context) error {
		var (
			user = c.Sender()
		)
		userName := strconv.Itoa(int(user.ID)) + "-" + user.Username
		hist := getHistory(userName)
		if hist == "" {
			return c.Send("_(empty)_", &tele.SendOptions{
				ParseMode: "markdown",
			})
		} else {
			resp := hist
			resp += "\n======\n"
			resp += "Length in runes: " + strconv.Itoa(utf8.RuneCountInString(hist)) + "\n"
			resp += "History size limit: " + strconv.Itoa(gCfg.HistorySize) + "\n"
			return c.Send(resp)
		}
	})

	theBot.Handle("/config", func(c tele.Context) error {
		resp := ""
		resp += "*CONFIG*\n"
		var cfgCopy = gCfg
		cfgCopy.OpenAiKey = ""
		cfgCopy.BotApiKey = ""
		bytes, _ := yaml.Marshal(cfgCopy)
		resp += "```\n"
		resp += string(bytes)
		resp += "```\n"
		return c.Send(resp, &tele.SendOptions{
			ParseMode: "markdown",
		})
	})

	theBot.Handle("/start", func(c tele.Context) error {
		var (
			user = c.Sender()
		)
		userName := strconv.Itoa(int(user.ID)) + "-" + user.Username
		clearHistory(userName)
		return c.Send("History cleared")
	})

	ctx := context.Background()
	theBot.Handle(tele.OnText, func(c tele.Context) error {
		// All the text messages that weren't
		// captured by existing handlers.

		var (
			user = c.Sender()
			text = c.Text()
		)

		log.Printf("[%s] %s", user.Username, text)
		c.Notify(tele.Typing)

		userName := strconv.Itoa(int(user.ID)) + "-" + user.Username

		history := getHistory(userName)
		humanPart := "Human: " + text + endSequence + "\n"

		log.Println("SENDING TO OPENAI API:")
		PfBlue(history + humanPart)
		req := gogpt.CompletionRequest{
			Model:       gogpt.GPT3TextDavinci003,
			MaxTokens:   gCfg.MaxTokens,
			Prompt:      history + humanPart + "AI: ",
			Temperature: gCfg.Temperature,
			Stop:        []string{endSequence},
		}

		response := ""

		completionResponse, err := openAiClient.CreateCompletion(ctx, req)
		if err != nil {
			log.Println("Error:", err)
			response = err.Error()
		} else {
			log.Println("OPENAI API RESPONSE:")
			PfGreen(strings.Trim(completionResponse.Choices[0].Text, "\n") + "\n")
			response = strings.Trim(completionResponse.Choices[0].Text, "\n")
		}

		PfYellow("USAGE: %+v \n", completionResponse.Usage)

		aiPart := "AI: " + response + endSequence + "\n"
		saveToHistory(userName, humanPart, aiPart)

		return c.Send(response, &tele.SendOptions{
			ReplyTo: c.Message(),
		})
	})

	theBot.Start()

}
