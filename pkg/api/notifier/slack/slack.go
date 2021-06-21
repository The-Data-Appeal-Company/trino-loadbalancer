package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var NotificationColor = "#F35500"
var DefaultTimeout = 10 * time.Second

type Slack struct {
	Webhook string
	Client  *http.Client
}

type Message struct {
	Message     string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

type Field struct {
	Title string `yaml:"title,omitempty" json:"title,omitempty"`
	Value string `yaml:"value,omitempty" json:"value,omitempty"`
	Short *bool  `yaml:"short,omitempty" json:"short,omitempty"`
}

func FieldsFromMap(m map[string]string) []Field {
	fields := make([]Field, 0)
	for k, v := range m {
		short := len(v) < 20
		fields = append(fields, Field{
			Title: k,
			Value: v,
			Short: &short,
		})
	}
	return fields
}

type Attachment struct {
	Title      string   `json:"title,omitempty"`
	TitleLink  string   `json:"title_link,omitempty"`
	Pretext    string   `json:"pretext,omitempty"`
	Text       string   `json:"text"`
	Fallback   string   `json:"fallback"`
	CallbackID string   `json:"callback_id"`
	Fields     []Field  `json:"fields,omitempty"`
	ImageURL   string   `json:"image_url,omitempty"`
	ThumbURL   string   `json:"thumb_url,omitempty"`
	Footer     string   `json:"footer"`
	Color      string   `json:"color,omitempty"`
	MrkdwnIn   []string `json:"mrkdwn_in,omitempty"`
}

func NewSlack(webhook string) Slack {
	return Slack{
		Webhook: webhook,
		Client:  &http.Client{Timeout: DefaultTimeout},
	}
}

func (s Slack) Send(message Message) error {
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
