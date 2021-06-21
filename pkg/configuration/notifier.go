package configuration

import (
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/notifier"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/notifier/slack"
)

type NotifierConfig struct {
	Slack struct {
		Webhook string `json:"webhook" yaml:"webhook" mapstructure:"webhook"`
	} `json:"slack" yaml:"slack" mapstructure:"slack"`
}

func CreateNotifier(conf NotifierConfig) notifier.Notifier {
	notifiers := make([]notifier.Notifier, 0)

	if len(conf.Slack.Webhook) != 0 {
		notifiers = append(notifiers, notifier.NewSlackNotifier(slack.NewSlack(conf.Slack.Webhook)))
	}

	return notifier.NewMultiNotifier(notifiers...)
}
