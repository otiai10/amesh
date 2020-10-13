package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func postMessage(message interface{}) error {
	body := bytes.NewBuffer(nil)
	if err := json.NewEncoder(body).Encode(message); err != nil {
		return err
	}
	// https://api.slack.com/methods/chat.postMessage
	req, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("SLACK_BOT_USER_OAUTH_ACCESS_TOKEN")))
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf(res.Status)
	}
	return nil
}
