package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	MessageModel "github.com/taufikardiyan28/chat/model/messages"
)

type M map[string]interface{}

func (c *Connection) SendPushNotification(msg MessageModel.MessagePayload) error {
	fmt.Println("send push notif")
	url := c.Config.PushNotif.Url

	payload, err := json.Marshal(msg)
	//err := json.NewEncoder(payload).Encode(msg)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	client := new(http.Client)
	response, err := client.Do(request)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}
