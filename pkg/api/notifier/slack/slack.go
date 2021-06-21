package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var NotificationColor = "#F35A00"
var DefaultTimeout = 10 * time.Second

type Slack struct {
	Webhook string
	Client  *http.Client
}

type SlackMessage struct {
	Message     string       `json:"text"`
	Attachments []attachment `json:"attachments"`
}

type SlackField struct {
	Title string `yaml:"title,omitempty" json:"title,omitempty"`
	Value string `yaml:"value,omitempty" json:"value,omitempty"`
	Short *bool  `yaml:"short,omitempty" json:"short,omitempty"`
}

type attachment struct {
	Title      string       `json:"title,omitempty"`
	TitleLink  string       `json:"title_link,omitempty"`
	Pretext    string       `json:"pretext,omitempty"`
	Text       string       `json:"text"`
	Fallback   string       `json:"fallback"`
	CallbackID string       `json:"callback_id"`
	Fields     []SlackField `json:"fields,omitempty"`
	ImageURL   string       `json:"image_url,omitempty"`
	ThumbURL   string       `json:"thumb_url,omitempty"`
	Footer     string       `json:"footer"`
	Color      string       `json:"color,omitempty"`
	MrkdwnIn   []string     `json:"mrkdwn_in,omitempty"`
}

func NewSlack(webhook string) Slack {
	return Slack{
		Webhook: webhook,
		Client:  &http.Client{Timeout: DefaultTimeout},
	}
}

func (s Slack) Send(message SlackMessage) error {
	messageBody, err := json.Marshal(message)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, s.Webhook, bytes.NewBuffer(messageBody))

	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := s.Client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d from slack webhook %s", res.StatusCode, s.Webhook)
	}

	return nil
}
