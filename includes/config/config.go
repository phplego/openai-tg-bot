package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"openai-tg-bot/includes/types"
)

var (
	Data types.Config = types.Config{
		// default config values here
		Temperature: 0.7,
		MaxTokens:   512,
		HistorySize: 1024,
	}
	EndSequence = " ␃"
)

func Init() {
	// load file
	err := cleanenv.ReadConfig("config.yaml", &Data)
	if err != nil {
		log.Println(err)
	}
	log.Printf("CONFIG: %+v", Data)
}
