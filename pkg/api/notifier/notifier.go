package notifier

import "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/notifier/slack"

type Request struct {
	Message string
}

type Notifier interface {
	Notify(Request) error
}

type SlackNotifier struct {
	slack slack.Slack
}

func Noop() NoopNotifier {
	return NoopNotifier{}
}

type NoopNotifier struct{}

func (n NoopNotifier) Notify(request Request) error {
	return nil
}

func NewSlackNotifier(slack slack.Slack) *SlackNotifier {
	return &SlackNotifier{slack: slack}
}

func (s SlackNotifier) Notify(request Request) error {
	return s.slack.Send(slack.SlackMessage{
		Message: request.Message,
	})
}

type MultiNotifier struct {
	notifiers []Notifier
}

func NewMultiNotifier(notifiers ...Notifier) *MultiNotifier {
	return &MultiNotifier{notifiers: notifiers}
}

func (m MultiNotifier) Notify(request Request) error {
	for _, notifier := range m.notifiers {
		if err := notifier.Notify(request); err != nil {
			return err
		}
	}
	return nil
}
