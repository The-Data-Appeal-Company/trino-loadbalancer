package slack

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSlackNotificationForTerminationStatusOk(t *testing.T) {
	webhookPath := "/services/T07M5HUDA/BQ1U5VDGA/yhpIczRK0cZ3jDLK1U8qD634"

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		require.Equal(t, webhookPath, req.URL.Path)
		res.WriteHeader(200)
		_, err := res.Write([]byte("ok"))
		require.NoError(t, err)
	}))
	defer testServer.Close()

	slack := NewSlack(testServer.URL + webhookPath)
	err := slack.Send(SlackMessage{
		Message: "test message",
	})

	require.NoError(t, err)
}

func TestSlackNotificationForTerminationStatus500(t *testing.T) {
	webhookPath := "/services/T07M5HUDA/BQ1U5VDGA/yhpIczRK0cZ3jDLK1U8qD634"

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		require.Equal(t, webhookPath, req.URL.Path)
		res.WriteHeader(500)
		_, err := res.Write([]byte("ok"))
		require.NoError(t, err)
	}))
	defer testServer.Close()

 	slack := NewSlack(testServer.URL + webhookPath)
	err := slack.Send(SlackMessage{
		Message: "test message",
	})

	require.Error(t, err)
}
