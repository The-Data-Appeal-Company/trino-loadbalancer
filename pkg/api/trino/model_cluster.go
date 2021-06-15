package trino

type ClusterStatistics struct {
	RunningQueries   int32   `json:"runningQueries"`
	BlockedQueries   int32   `json:"blockedQueries"`
	QueuedQueries    int32   `json:"queuedQueries"`
	ActiveWorkers    int32   `json:"activeWorkers"`
	RunningDrivers   int32   `json:"runningDrivers"`
	ReservedMemory   float64 `json:"reservedMemory"`
	TotalInputRows   int64   `json:"totalInputRows"`
	TotalInputBytes  int64   `json:"totalInputBytes"`
	TotalCPUTimeSecs int32   `json:"totalCpuTimeSecs"`
}
