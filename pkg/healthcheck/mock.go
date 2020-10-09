package healthcheck

import (
	"net/url"
)

type MockHealthCheck struct {
	state Health
	err   error
}

func Mock(state Health, err error) MockHealthCheck {
	return MockHealthCheck{
		state: state,
		err:   err,
	}
}

func (m MockHealthCheck) Check(url *url.URL) (Health, error) {
	return m.state, m.err
}
