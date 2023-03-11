package user_state

import (
	"encoding/json"
	"log"
	"openai-tg-bot/includes/commands"
	"openai-tg-bot/includes/types"
	"os"
	"strconv"
	"unicode/utf8"
)

// UserState runtime user state (mode etc)
type UserState struct {
	Mode                commands.Mode `json:"mode"`
	ConversationHistory string        `json:"conversation-history"`
}

var AllUserStates map[int64]UserState

var gCfg types.Config

func Init(config types.Config) {
	gCfg = config
	AllUserStates = make(map[int64]UserState)
}

func getFilename(userId int64) string {
	return "state" + strconv.Itoa(int(userId)) + ".txt"
}

func Save(userId int64, state UserState) error {

	// set to map
	AllUserStates[userId] = state

	// save to file
	bytes, err := json.Marshal(state)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(getFilename(userId), os.O_WRONLY|os.O_CREATE, 0600)
	defer file.Close()
	if err != nil {
		return err
	}
	_, err = file.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func Load(userId int64) error {
	bytes, err := os.ReadFile(getFilename(userId))
	if os.IsNotExist(err) {
		AllUserStates[userId] = UserState{}
		return nil
	}
	if err != nil {
		return err
	}
	var state UserState
	err = json.Unmarshal(bytes, &state)
	if err != nil {
		return err
	}
	AllUserStates[userId] = state
	return nil
}

func FindById(userId int64) UserState {
	state, ok := AllUserStates[userId]
	// If user state exists
	if ok {
		return state
	}

	return UserState{}
}

func AppendToHistory(userId int64, userMessage string, aiResponse string) {
	hist := GetHistory(userId)

	newHist := hist + userMessage + aiResponse
	start := utf8.RuneCountInString(newHist) - gCfg.HistorySize
	if start < 0 {
		start = 0
	}
	newHist = string([]rune(newHist)[start:]) // convert to rune is important, because we don't want to break utf-8
	state := FindById(userId)
	state.ConversationHistory = newHist
	err := Save(userId, state)
	if err != nil {
		log.Println(err)
	}
}

func GetHistory(userId int64) string {
	state := FindById(userId)
	return state.ConversationHistory
}

func ClearHistory(userId int64) {
	state := FindById(userId)
	state.ConversationHistory = ""
	err := Save(userId, state)
	if err != nil {
		log.Println(err)
	}
}
