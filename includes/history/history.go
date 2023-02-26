package history

import (
	"log"
	"openai-tg-bot/includes/types"
	"os"
	"unicode/utf8"
)

var gCfg types.Config

func Init(config types.Config) {
	gCfg = config
}

func SaveToHistory(userName string, userMessage string, aiResponse string) {
	hist := GetHistory(userName)

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

func GetHistory(userName string) string {
	b, err := os.ReadFile(userName + ".txt")
	if err != nil {
		log.Println(err)
		return ""
	}

	str := string(b) // convert content to a 'string'
	return str
}

func ClearHistory(userName string) {
	file, err := os.OpenFile(userName+".txt", os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0600)
	defer file.Close()
	if err != nil {
		log.Println(err)
	}
}
