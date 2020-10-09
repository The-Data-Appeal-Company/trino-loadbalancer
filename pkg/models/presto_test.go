package models

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMarshalPrestoResponse(t *testing.T) {
	data := `{"id":"20200924_095706_01798_yi2gi","infoUri":"http://localhost:8080/ui/query.html?20200924_095706_01798_yi2gi","nextUri":"http://localhost:8080/v1/statement/20200924_095706_01798_yi2gi/1?slug=xc7951ca2b9124141a6baa68448edb219","stats":{"state":"QUEUED","queued":true,"scheduled":false,"nodes":0,"totalSplits":0,"queuedSplits":0,"runningSplits":0,"completedSplits":0,"cpuTimeMillis":0,"wallTimeMillis":0,"queuedTimeMillis":0,"elapsedTimeMillis":0,"processedRows":0,"processedBytes":0,"peakMemoryBytes":0,"spilledBytes":0},"warnings":[]}`

	var queryState PrestoQueryState
	err := json.Unmarshal([]byte(data), &queryState)

	require.NoError(t, err)
	require.Equal(t, "20200924_095706_01798_yi2gi", queryState.ID)
}
