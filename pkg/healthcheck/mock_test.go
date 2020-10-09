package healthcheck

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMock(t *testing.T) {
	s := Health{}
	e := errors.New("generic error")

	m := MockHealthCheck{
		state: s,
		err:   e,
	}

	state, err := m.Check(nil)

	require.Equal(t, s, state)
	require.Equal(t, e, err)

}
