package statistics

import (
	"fmt"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/models"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/tests"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApiRetrieverPrestoDB(t *testing.T) {
	const serverResponse = `
{
  "runningQueries": 1,
  "blockedQueries": 0,
  "queuedQueries": 0,
  "activeWorkers": 5,
  "runningDrivers": 0,
  "reservedMemory": 0,
  "totalInputRows": 360508662922,
  "totalInputBytes": 5799000097802,
  "totalCpuTimeSecs": 264808
}
`
	apiSrv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/v1/cluster" {
			_, err := writer.Write([]byte(serverResponse))
			require.NoError(t, err)
			return
		}

		writer.WriteHeader(http.StatusNotFound)
	}))
	defer apiSrv.Close()

	retriever := NewPrestoClusterApi()

	stats, err := retriever.GetStatistics(models.Coordinator{
		Name:         "test",
		URL:          tests.MustUrl(apiSrv.URL),
		Enabled:      true,
		Distribution: models.PrestoDistDb,
	})
	require.NoError(t, err)

	require.Equal(t, stats.RunningQueries, int32(1))
	require.Equal(t, stats.BlockedQueries, int32(0))
	require.Equal(t, stats.QueuedQueries, int32(0))
	require.Equal(t, stats.ActiveWorkers, int32(5))
	require.Equal(t, stats.RunningDrivers, int32(0))
	require.Equal(t, stats.ReservedMemory, float64(0))
	require.Equal(t, stats.TotalInputRows, int64(360508662922))
	require.Equal(t, stats.TotalInputBytes, int64(5799000097802))
	require.Equal(t, stats.TotalCPUTimeSecs, int32(264808))
}

func TestApiRetrieverPrestoSQL(t *testing.T) {
	const serverResponse = `
{
  "runningQueries": 1,
  "blockedQueries": 0,
  "queuedQueries": 0,
  "activeWorkers": 5,
  "runningDrivers": 0,
  "reservedMemory": 0,
  "totalInputRows": 360508662922,
  "totalInputBytes": 5799000097802,
  "totalCpuTimeSecs": 264808
}
`
	apiSrv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		const prestoUser = "prestolb"
		const prestoAuthCookie = "test"

		if request.URL.Path != "/ui/login" && request.Header.Get("Cookie") != prestoAuthCookie {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		if request.URL.Path == "/ui/api/stats" {
			_, err := writer.Write([]byte(serverResponse))
			require.NoError(t, err)
			return
		}

		if request.URL.Path == "/ui/login" {
			body, err := ioutil.ReadAll(request.Body)
			require.NoError(t, err)
			defer request.Body.Close()

			require.Equal(t, fmt.Sprintf("username=%s&password=&redirectPath=", prestoUser), string(body))

			writer.Header().Set("Content-Type", "application/x-www-form-urlencoded")
			writer.Header().Set("Set-Cookie", prestoAuthCookie)
			writer.Header().Set("Location", "http://presto.local")
			writer.WriteHeader(http.StatusSeeOther)
			return
		}

		writer.WriteHeader(http.StatusNotFound)
	}))
	defer apiSrv.Close()

	retriever := NewPrestoClusterApi()

	stats, err := retriever.GetStatistics(models.Coordinator{
		Name:         "test",
		URL:          tests.MustUrl(apiSrv.URL),
		Enabled:      true,
		Distribution: models.PrestoDistSql,
	})
	require.NoError(t, err)

	require.Equal(t, stats.RunningQueries, int32(1))
	require.Equal(t, stats.BlockedQueries, int32(0))
	require.Equal(t, stats.QueuedQueries, int32(0))
	require.Equal(t, stats.ActiveWorkers, int32(5))
	require.Equal(t, stats.RunningDrivers, int32(0))
	require.Equal(t, stats.ReservedMemory, float64(0))
	require.Equal(t, stats.TotalInputRows, int64(360508662922))
	require.Equal(t, stats.TotalInputBytes, int64(5799000097802))
	require.Equal(t, stats.TotalCPUTimeSecs, int32(264808))
}

func TestApiRetrieverFailOn404(t *testing.T) {
	apiSrv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNotFound)
	}))
	defer apiSrv.Close()

	retriever := NewPrestoClusterApi()

	_, err := retriever.GetStatistics(models.Coordinator{
		Name:         "test",
		URL:          tests.MustUrl(apiSrv.URL),
		Enabled:      true,
		Distribution: models.PrestoDistDb,
	})
	require.Error(t, err)
}

func TestApiRetrieverFailOnMalformedJson(t *testing.T) {
	const serverResponse = `
{
  "test: true
}
`
	apiSrv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/v1/cluster" {
			_, err := writer.Write([]byte(serverResponse))
			require.NoError(t, err)
			return
		}

		writer.WriteHeader(http.StatusNotFound)
	}))
	defer apiSrv.Close()

	retriever := NewPrestoClusterApi()

	_, err := retriever.GetStatistics(models.Coordinator{
		Name:         "test",
		URL:          tests.MustUrl(apiSrv.URL),
		Enabled:      true,
		Distribution: models.PrestoDistDb,
	})

	require.Error(t, err)

}
