package open_ai

import (
	gogpt "github.com/sashabaranov/go-gpt3"
	"openai-tg-bot/includes/config"
)

var OpenAiClient *gogpt.Client

func Init() {
	OpenAiClient = gogpt.NewClient(config.Data.OpenAiKey)
}
