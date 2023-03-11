package on_text

import (
	"context"
	gogpt "github.com/sashabaranov/go-gpt3"
	tele "gopkg.in/telebot.v3"
	"log"
	"openai-tg-bot/includes/commands"
	"openai-tg-bot/includes/config"
	open_ai "openai-tg-bot/includes/open-ai"
	"openai-tg-bot/includes/printing"
	user_state "openai-tg-bot/includes/user-state"
	"strings"
)

func Handler(c tele.Context) error {
	var (
		user = c.Sender()
	)

	user_state.Load(user.ID)
	userState := user_state.FindById(user.ID)

	switch userState.Mode {
	case commands.TextMode:
		return HandleTextMode(c)
	case commands.ImageMode:
		return HandleImageMode(c)
	}

	return nil
}

func HandleTextMode(c tele.Context) error {
	var (
		user = c.Sender()
		text = c.Text()
	)

	// load state
	err := user_state.Load(user.ID)
	if err != nil {
		return err
	}

	log.Printf("[%s] %s", user.Username, text)
	c.Notify(tele.Typing)

	hist := user_state.GetHistory(user.ID)
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
	user_state.AppendToHistory(user.ID, humanPart, aiPart)

	return c.Send(response, &tele.SendOptions{
		ReplyTo: c.Message(),
	})
}

func HandleImageMode(c tele.Context) error {
	c.Notify(tele.Typing)

	log.Println("SENDING TO OPENAI API:")
	printing.PfBlue(c.Text() + "\n")
	req := gogpt.ImageRequest{
		Prompt: c.Text(),
	}

	response := ""

	ctx := context.Background()
	imageResponse, err := open_ai.OpenAiClient.CreateImage(ctx, req)
	if err != nil {
		log.Println("Error:", err)
		response = err.Error()
	} else {
		log.Println("OPENAI API RESPONSE:")
		printing.PfGreen(strings.Trim(imageResponse.Data[0].URL, "\n") + "\n")
		response = strings.Trim(imageResponse.Data[0].URL, "\n")
	}

	return c.Send(response, &tele.SendOptions{
		ReplyTo: c.Message(),
	})
}
