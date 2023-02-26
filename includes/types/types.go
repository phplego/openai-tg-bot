package types

// Config is an application configuration structure
type Config struct {
	OpenAiKey   string  `yaml:"open-ai-key,omitempty"`
	BotApiKey   string  `yaml:"bot-api-key,omitempty"`
	Temperature float32 `yaml:"temperature"`
	MaxTokens   int     `yaml:"max-tokens"`
	HistorySize int     `yaml:"history-size"`
}
