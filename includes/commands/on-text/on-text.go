package on_text

import (
	"context"
	gogpt "github.com/sashabaranov/go-gpt3"
	tele "gopkg.in/telebot.v3"
	"log"
	"openai-tg-bot/includes/config"
	"openai-tg-bot/includes/history"
	open_ai "openai-tg-bot/includes/open-ai"
	"openai-tg-bot/includes/printing"
	"strconv"
	"strings"
)

func Handler(c tele.Context) error {
	var (
		user = c.Sender()
		text = c.Text()
	)

	log.Printf("[%s] %s", user.Username, text)
	c.Notify(tele.Typing)

	userName := strconv.Itoa(int(user.ID)) + "-" + user.Username

	hist := history.GetHistory(userName)
	humanPart := "Human: " + text + config.EndSequence + "\n"

	log.Println("SENDING TO OPENAI API:")
	printing.PfBlue(hist + humanPart)
	req := gogpt.CompletionRequest{
		Model:       gogpt.GPT3TextDavinci003,
		MaxTokens:   config.Data.MaxTokens,
		Prompt:      hist + humanPart + "AI: ",
		Temperature: config.Data.Temperature,
		Stop:        []string{config.EndSequence},
	}

	response := ""

	ctx := context.Background()
	completionResponse, err := open_ai.OpenAiClient.CreateCompletion(ctx, req)
	if err != nil {
		log.Println("Error:", err)
		response = err.Error()
	} else {
		log.Println("OPENAI API RESPONSE:")
		printing.PfGreen(strings.Trim(completionResponse.Choices[0].Text, "\n") + "\n")
		response = strings.Trim(completionResponse.Choices[0].Text, "\n")
	}

	printing.PfYellow("USAGE: %+v \n", completionResponse.Usage)

	aiPart := "AI: " + response + config.EndSequence + "\n"
	history.SaveToHistory(userName, humanPart, aiPart)

	return c.Send(response, &tele.SendOptions{
		ReplyTo: c.Message(),
	})
}
