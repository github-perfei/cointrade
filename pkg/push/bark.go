package push

import (
	"fmt"
	"net/http"
	"net/url"
)

const barkPushAPI = "https://api.day.app/%s/%s"

type barkPush struct {
	appID string
}

func NewBarkPush(appID string) pusher {
	if appID == "" {
		panic("NewBarkPush fail appID is empty")
	}
	return &barkPush{
		appID,
	}
}

func (t *barkPush) SendMessage(message string) error {
	resp, err := http.PostForm(fmt.Sprintf(barkPushAPI, t.appID, url.QueryEscape(message)), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Printf("Bark发送消息状态码: %v", resp.StatusCode)
	return nil
}
