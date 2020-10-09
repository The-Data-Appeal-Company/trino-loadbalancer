package lb

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestIsStatementRequest(t *testing.T) {
	valid := []string{
		"http://presto.cluster:8889/v1/statement/20200924_102554_02623_yi2gi/1?slug=xd9655d44c64d43119d3126cd47f2b6d0",
		"http://presto.cluster:8889/v1/statement/20200924_102554_02623_yi2gi/",
		"http://presto.cluster:8889/v1/statement/20200924_102554_02623_yi2gi",
	}

	invalid := []string{
		"http://presto.cluster:8889/v1/cluster",
		"http://presto.cluster:8889/v2/statement/20200924_102554_02623_yi2gi/",
		"http://presto.cluster:8889/v1/info",
	}

	for _, u := range valid {
		require.True(t, isStatementRequest(mustUrl(u)))
	}
	for _, u := range invalid {
		require.False(t, isStatementRequest(mustUrl(u)))
	}

}
func TestExtractQueryInfoFromRequest(t *testing.T) {

	urls := []string{
		"http://presto.local:8889/v1/statement/20200924_102554_02623_yi2gi/1?slug=xd9655d44c64d43119d3126cd47f2b6d0",
		"http://presto.local:8889/v1/statement/20200924_102554_02623_yi2gi/",
		"http://presto.local:8889/v1/statement/20200924_102554_02623_yi2gi",
	}

	for _, u := range urls {
		headers := http.Header{}
		headers.Add(PrestoHeaderUser, "test-user")
		headers.Add(PrestoHeaderTransaction, "test-tx")

		queryInfo, err := QueryInfoFromRequest(&http.Request{
			Method: "POST",
			URL:    mustUrl(u),
			Header: headers,
		})
		require.NoError(t, err)
		require.Equal(t, queryInfo.QueryID, "20200924_102554_02623_yi2gi")
		require.Equal(t, queryInfo.TransactionID, "test-tx")
		require.Equal(t, queryInfo.User, "test-user")
	}
}

func TestExtractQueryInfoFromResponse(t *testing.T) {

	body := `{"id":"20200924_095706_01798_yi2gi","infoUri":"http://localhost:8080/ui/query.html?20200924_095706_01798_yi2gi","nextUri":"http://localhost:8080/v1/statement/20200924_095706_01798_yi2gi/1?slug=xc7951ca2b9124141a6baa68448edb219","stats":{"state":"QUEUED","queued":true,"scheduled":false,"nodes":0,"totalSplits":0,"queuedSplits":0,"runningSplits":0,"completedSplits":0,"cpuTimeMillis":0,"wallTimeMillis":0,"queuedTimeMillis":0,"elapsedTimeMillis":0,"processedRows":0,"processedBytes":0,"peakMemoryBytes":0,"spilledBytes":0},"warnings":[]}`

	headers := http.Header{}
	headers.Add(PrestoHeaderUser, "test-user")
	headers.Add(PrestoHeaderTransaction, "test-tx")

	req := &http.Request{
		Method: "POST",
		URL:    mustUrl("http://localhost:1234"),
		Header: headers,
	}

	res := &http.Response{
		StatusCode: http.StatusOK,
		Body:       bodyReadCloser(body),
	}

	queryInfo, err := QueryInfoFromResponse(req, res)
	require.NoError(t, err)
	require.Equal(t, queryInfo.QueryID, "20200924_095706_01798_yi2gi")
	require.Equal(t, queryInfo.TransactionID, "test-tx")
	require.Equal(t, queryInfo.User, "test-user")

}

func mustUrl(raw string) *url.URL {
	parsed, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return parsed
}

func bodyReadCloser(val string) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewBuffer([]byte(val)))
}
