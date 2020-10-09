package models

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

type QueryInfo struct {
	User          string
	QueryID       string
	TransactionID string
}

type PrestoQueryState struct {
	ID      string  `json:"id"`
	InfoURI string  `json:"infoUri"`
	NextURI *string `json:"nextUri"`
	Columns []struct {
		Name          string `json:"name"`
		Type          string `json:"type"`
		TypeSignature struct {
			RawType          string        `json:"rawType"`
			TypeArguments    []interface{} `json:"typeArguments"`
			LiteralArguments []interface{} `json:"literalArguments"`
			Arguments        []struct {
				Kind  string      `json:"kind"`
				Value interface{} `json:"value"`
			} `json:"arguments"`
		} `json:"typeSignature"`
	} `json:"columns"`
	Data  interface{} `json:"data"`
	Stats struct {
		State             string `json:"state"`
		Queued            bool   `json:"queued"`
		Scheduled         bool   `json:"scheduled"`
		Nodes             int    `json:"nodes"`
		TotalSplits       int    `json:"totalSplits"`
		QueuedSplits      int    `json:"queuedSplits"`
		RunningSplits     int    `json:"runningSplits"`
		CompletedSplits   int    `json:"completedSplits"`
		CPUTimeMillis     int    `json:"cpuTimeMillis"`
		WallTimeMillis    int    `json:"wallTimeMillis"`
		QueuedTimeMillis  int    `json:"queuedTimeMillis"`
		ElapsedTimeMillis int    `json:"elapsedTimeMillis"`
		ProcessedRows     int    `json:"processedRows"`
		ProcessedBytes    int    `json:"processedBytes"`
		PeakMemoryBytes   int    `json:"peakMemoryBytes"`
		SpilledBytes      int    `json:"spilledBytes"`
		RootStage         struct {
			StageID         string        `json:"stageId"`
			State           string        `json:"state"`
			Done            bool          `json:"done"`
			Nodes           int           `json:"nodes"`
			TotalSplits     int           `json:"totalSplits"`
			QueuedSplits    int           `json:"queuedSplits"`
			RunningSplits   int           `json:"runningSplits"`
			CompletedSplits int           `json:"completedSplits"`
			CPUTimeMillis   int           `json:"cpuTimeMillis"`
			WallTimeMillis  int           `json:"wallTimeMillis"`
			ProcessedRows   int           `json:"processedRows"`
			ProcessedBytes  int           `json:"processedBytes"`
			SubStages       []interface{} `json:"subStages"`
		} `json:"rootStage"`
		ProgressPercentage float64 `json:"progressPercentage"`
	} `json:"stats"`
	Warnings []interface{} `json:"warnings"`
}
