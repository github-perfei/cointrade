package push

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const telegramPushAPI = "https://api.telegram.org/bot%s/sendMessage"

type telegramPush struct {
	botToken  string
	chatID    string
	httpProxy string
}

func NewTelegramPush(api, botToken, chatID, httpProxy string) pusher {
	if api == "" || botToken == "" || chatID == "" {
		panic("NewTelegramPush fail api or botToken or chatID is empty")
	}
	return &telegramPush{
		botToken,
		chatID,
		httpProxy,
	}
}

func (t *telegramPush) SendMessage(message string) error {
	data := make(map[string]string)
	data["chat_id"] = t.chatID
	data["text"] = message
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf(telegramPushAPI, t.botToken), "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Printf("Telegram发送消息状态码: %v", resp.StatusCode)
	return nil
}
