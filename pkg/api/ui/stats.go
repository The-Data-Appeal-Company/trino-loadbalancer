package ui

import (
	"encoding/json"
	"net/http"
)

type StatsApiResponse struct {
	TotalWorkers   int32 `json:"total_workers"`
	RunningQueries int32 `json:"running_queries"`
	BlockedQueries int32 `json:"blocked_queries"`
	QueuedQueries  int32 `json:"queued_queries"`
}

func (a Api) statistics(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	clusters, err := a.discoveryStorage.All(ctx)
	if err != nil {
		if _, err := writer.Write([]byte(err.Error())); err != nil {
			a.logger.Error("error writing response: %w", err)
		}
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	var workers int32
	var runningQueries int32
	var blockedQueries int32
	var queuedQueries int32

	for _, cluster := range clusters {
		if !cluster.Enabled {
			continue
		}

		stats, err := a.statsRetriever.ClusterStatistics(cluster)
		if err != nil {
			a.logger.Warn("unable to get stats from %s", cluster.Name)
			continue
		}

		workers += stats.ActiveWorkers
		runningQueries += stats.RunningQueries
		blockedQueries += stats.BlockedQueries
		queuedQueries += stats.QueuedQueries
	}

	response := StatsApiResponse{
		TotalWorkers:   workers,
		RunningQueries: runningQueries,
		BlockedQueries: blockedQueries,
		QueuedQueries:  queuedQueries,
	}

	body, err := json.Marshal(response)
	if err != nil {
		if _, err := writer.Write([]byte(err.Error())); err != nil {
			a.logger.Error("error writing response: %w", err)
		}
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := writer.Write(body); err != nil {
		a.logger.Error("error writing response: %w", err)
	}

}
