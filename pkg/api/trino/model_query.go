package trino

import "time"

type QueryInfo struct {
	User          string
	QueryID       string
	TransactionID string
}

type QueryState struct {
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

func (taskDetail Tasks) GetElapsedTime() (time.Duration, error) {
	return time.ParseDuration(taskDetail.Stats.ElapsedTime)
}

type QueryDetail struct {
	QueryID                       string                  `json:"queryId"`
	Session                       Session                 `json:"session"`
	State                         string                  `json:"state"`
	MemoryPool                    string                  `json:"memoryPool"`
	Scheduled                     bool                    `json:"scheduled"`
	Self                          string                  `json:"self"`
	FieldNames                    []string                `json:"fieldNames"`
	Query                         string                  `json:"query"`
	QueryStats                    QueryStats              `json:"queryStats"`
	SetSessionProperties          SetSessionProperties    `json:"setSessionProperties"`
	ResetSessionProperties        []interface{}           `json:"resetSessionProperties"`
	SetRoles                      SetRoles                `json:"setRoles"`
	AddedPreparedStatements       AddedPreparedStatements `json:"addedPreparedStatements"`
	DeallocatedPreparedStatements []interface{}           `json:"deallocatedPreparedStatements"`
	ClearTransactionID            bool                    `json:"clearTransactionId"`
	OutputStage                   OutputStage             `json:"outputStage"`
	Warnings                      []interface{}           `json:"warnings"`
	Routines                      []Routines              `json:"routines"`
	ResourceGroupID               []string                `json:"resourceGroupId"`
	QueryType                     string                  `json:"queryType"`
	FinalQueryInfo                bool                    `json:"finalQueryInfo"`
}
type Path struct {
}
type ResourceEstimates struct {
}
type SystemProperties struct {
}
type CatalogProperties struct {
}
type UnprocessedCatalogProperties struct {
}
type Roles struct {
}
type PreparedStatements struct {
}
type Session struct {
	QueryID                      string                       `json:"queryId"`
	TransactionID                string                       `json:"transactionId"`
	ClientTransactionSupport     bool                         `json:"clientTransactionSupport"`
	User                         string                       `json:"user"`
	Groups                       []interface{}                `json:"groups"`
	Principal                    string                       `json:"principal"`
	Source                       string                       `json:"source"`
	Catalog                      string                       `json:"catalog"`
	Schema                       string                       `json:"schema"`
	Path                         Path                         `json:"path"`
	TimeZoneKey                  int                          `json:"timeZoneKey"`
	Locale                       string                       `json:"locale"`
	RemoteUserAddress            string                       `json:"remoteUserAddress"`
	UserAgent                    string                       `json:"userAgent"`
	ClientTags                   []interface{}                `json:"clientTags"`
	ClientCapabilities           []string                     `json:"clientCapabilities"`
	ResourceEstimates            ResourceEstimates            `json:"resourceEstimates"`
	Start                        time.Time                    `json:"start"`
	SystemProperties             SystemProperties             `json:"systemProperties"`
	CatalogProperties            CatalogProperties            `json:"catalogProperties"`
	UnprocessedCatalogProperties UnprocessedCatalogProperties `json:"unprocessedCatalogProperties"`
	Roles                        Roles                        `json:"roles"`
	PreparedStatements           PreparedStatements           `json:"preparedStatements"`
	ProtocolName                 string                       `json:"protocolName"`
}
type StageGcStatistics struct {
	StageID          int `json:"stageId"`
	Tasks            int `json:"tasks"`
	FullGcTasks      int `json:"fullGcTasks"`
	MinFullGcSec     int `json:"minFullGcSec"`
	MaxFullGcSec     int `json:"maxFullGcSec"`
	TotalFullGcSec   int `json:"totalFullGcSec"`
	AverageFullGcSec int `json:"averageFullGcSec"`
}
type DynamicFilterDomainStats struct {
	DynamicFilterID     string `json:"dynamicFilterId"`
	SimplifiedDomain    string `json:"simplifiedDomain"`
	RangeCount          int    `json:"rangeCount"`
	DiscreteValuesCount int    `json:"discreteValuesCount"`
	CollectionDuration  string `json:"collectionDuration"`
}
type DynamicFiltersStats struct {
	DynamicFilterDomainStats []DynamicFilterDomainStats `json:"dynamicFilterDomainStats"`
	LazyDynamicFilters       int                        `json:"lazyDynamicFilters"`
	ReplicatedDynamicFilters int                        `json:"replicatedDynamicFilters"`
	TotalDynamicFilters      int                        `json:"totalDynamicFilters"`
	DynamicFiltersCompleted  int                        `json:"dynamicFiltersCompleted"`
}

type QueryStats struct {
	CreateTime                        time.Time           `json:"createTime"`
	ExecutionStartTime                time.Time           `json:"executionStartTime"`
	LastHeartbeat                     time.Time           `json:"lastHeartbeat"`
	EndTime                           time.Time           `json:"endTime"`
	ElapsedTime                       string              `json:"elapsedTime"`
	QueuedTime                        string              `json:"queuedTime"`
	ResourceWaitingTime               string              `json:"resourceWaitingTime"`
	DispatchingTime                   string              `json:"dispatchingTime"`
	ExecutionTime                     string              `json:"executionTime"`
	AnalysisTime                      string              `json:"analysisTime"`
	PlanningTime                      string              `json:"planningTime"`
	FinishingTime                     string              `json:"finishingTime"`
	TotalTasks                        int                 `json:"totalTasks"`
	RunningTasks                      int                 `json:"runningTasks"`
	CompletedTasks                    int                 `json:"completedTasks"`
	TotalDrivers                      int                 `json:"totalDrivers"`
	QueuedDrivers                     int                 `json:"queuedDrivers"`
	RunningDrivers                    int                 `json:"runningDrivers"`
	BlockedDrivers                    int                 `json:"blockedDrivers"`
	CompletedDrivers                  int                 `json:"completedDrivers"`
	CumulativeUserMemory              interface{}         `json:"cumulativeUserMemory"`
	UserMemoryReservation             string              `json:"userMemoryReservation"`
	RevocableMemoryReservation        string              `json:"revocableMemoryReservation"`
	TotalMemoryReservation            string              `json:"totalMemoryReservation"`
	PeakUserMemoryReservation         string              `json:"peakUserMemoryReservation"`
	PeakRevocableMemoryReservation    string              `json:"peakRevocableMemoryReservation"`
	PeakNonRevocableMemoryReservation string              `json:"peakNonRevocableMemoryReservation"`
	PeakTotalMemoryReservation        string              `json:"peakTotalMemoryReservation"`
	PeakTaskUserMemory                string              `json:"peakTaskUserMemory"`
	PeakTaskRevocableMemory           string              `json:"peakTaskRevocableMemory"`
	PeakTaskTotalMemory               string              `json:"peakTaskTotalMemory"`
	Scheduled                         bool                `json:"scheduled"`
	TotalScheduledTime                string              `json:"totalScheduledTime"`
	TotalCPUTime                      string              `json:"totalCpuTime"`
	TotalBlockedTime                  string              `json:"totalBlockedTime"`
	FullyBlocked                      bool                `json:"fullyBlocked"`
	BlockedReasons                    []interface{}       `json:"blockedReasons"`
	PhysicalInputDataSize             string              `json:"physicalInputDataSize"`
	PhysicalInputPositions            int                 `json:"physicalInputPositions"`
	PhysicalInputReadTime             string              `json:"physicalInputReadTime"`
	InternalNetworkInputDataSize      string              `json:"internalNetworkInputDataSize"`
	InternalNetworkInputPositions     int                 `json:"internalNetworkInputPositions"`
	RawInputDataSize                  string              `json:"rawInputDataSize"`
	RawInputPositions                 int                 `json:"rawInputPositions"`
	ProcessedInputDataSize            string              `json:"processedInputDataSize"`
	ProcessedInputPositions           int                 `json:"processedInputPositions"`
	OutputDataSize                    string              `json:"outputDataSize"`
	OutputPositions                   int                 `json:"outputPositions"`
	PhysicalWrittenDataSize           string              `json:"physicalWrittenDataSize"`
	StageGcStatistics                 []StageGcStatistics `json:"stageGcStatistics"`
	DynamicFiltersStats               DynamicFiltersStats `json:"dynamicFiltersStats"`
	OperatorSummaries                 []OperatorSummaries `json:"operatorSummaries"`
	LogicalWrittenDataSize            string              `json:"logicalWrittenDataSize"`
	WrittenPositions                  int                 `json:"writtenPositions"`
	SpilledDataSize                   string              `json:"spilledDataSize"`
	ProgressPercentage                float64             `json:"progressPercentage"`
}
type SetSessionProperties struct {
}
type SetRoles struct {
}
type AddedPreparedStatements struct {
}
type Orderings struct {
	DateTrunc string `json:"date_trunc"`
}
type OrderingScheme struct {
	OrderBy   []string  `json:"orderBy"`
	Orderings Orderings `json:"orderings"`
}
type Source struct {
	Type              string         `json:"@type"`
	ID                string         `json:"id"`
	SourceFragmentIds []string       `json:"sourceFragmentIds"`
	Outputs           []string       `json:"outputs"`
	OrderingScheme    OrderingScheme `json:"orderingScheme"`
	ExchangeType      string         `json:"exchangeType"`
}
type Root struct {
	Type    string   `json:"@type"`
	ID      string   `json:"id"`
	Source  Source   `json:"source"`
	Columns []string `json:"columns"`
	Outputs []string `json:"outputs"`
}
type Symbols struct {
	Count            string `json:"count"`
	DateTrunc        string `json:"date_trunc"`
	ApproxPercentile string `json:"approx_percentile"`
	Round            string `json:"round"`
}
type ConnectorHandle struct {
	Type         string `json:"@type"`
	Partitioning string `json:"partitioning"`
	Function     string `json:"function"`
}

type Handle struct {
	ConnectorHandle ConnectorHandle `json:"connectorHandle"`
}

type StageExecutionDescriptor struct {
	Strategy                  string        `json:"strategy"`
	GroupedExecutionScanNodes []interface{} `json:"groupedExecutionScanNodes"`
}

type Costs struct {
}
type StatsAndCosts struct {
	Stats Stats `json:"stats"`
	Costs Costs `json:"costs"`
}

type GetSplitDistribution struct {
	Count interface{} `json:"count"`
	Total interface{} `json:"total"`
	P01   interface{} `json:"p01"`
	P05   interface{} `json:"p05"`
	P10   interface{} `json:"p10"`
	P25   interface{} `json:"p25"`
	P50   interface{} `json:"p50"`
	P75   interface{} `json:"p75"`
	P90   interface{} `json:"p90"`
	P95   interface{} `json:"p95"`
	P99   interface{} `json:"p99"`
	Min   interface{} `json:"min"`
	Max   interface{} `json:"max"`
	Avg   interface{} `json:"avg"`
}
type GcInfo struct {
	StageID          int `json:"stageId"`
	Tasks            int `json:"tasks"`
	FullGcTasks      int `json:"fullGcTasks"`
	MinFullGcSec     int `json:"minFullGcSec"`
	MaxFullGcSec     int `json:"maxFullGcSec"`
	TotalFullGcSec   int `json:"totalFullGcSec"`
	AverageFullGcSec int `json:"averageFullGcSec"`
}
type OperatorSummaries struct {
	StageID                        int     `json:"stageId"`
	PipelineID                     int     `json:"pipelineId"`
	OperatorID                     int     `json:"operatorId"`
	PlanNodeID                     string  `json:"planNodeId"`
	OperatorType                   string  `json:"operatorType"`
	TotalDrivers                   int     `json:"totalDrivers"`
	AddInputCalls                  int     `json:"addInputCalls"`
	AddInputWall                   string  `json:"addInputWall"`
	AddInputCPU                    string  `json:"addInputCpu"`
	PhysicalInputDataSize          string  `json:"physicalInputDataSize"`
	PhysicalInputPositions         int     `json:"physicalInputPositions"`
	InternalNetworkInputDataSize   string  `json:"internalNetworkInputDataSize"`
	InternalNetworkInputPositions  int     `json:"internalNetworkInputPositions"`
	RawInputDataSize               string  `json:"rawInputDataSize"`
	InputDataSize                  string  `json:"inputDataSize"`
	InputPositions                 int     `json:"inputPositions"`
	SumSquaredInputPositions       float64 `json:"sumSquaredInputPositions"`
	GetOutputCalls                 int     `json:"getOutputCalls"`
	GetOutputWall                  string  `json:"getOutputWall"`
	GetOutputCPU                   string  `json:"getOutputCpu"`
	OutputDataSize                 string  `json:"outputDataSize"`
	OutputPositions                int     `json:"outputPositions"`
	DynamicFilterSplitsProcessed   int     `json:"dynamicFilterSplitsProcessed"`
	PhysicalWrittenDataSize        string  `json:"physicalWrittenDataSize"`
	BlockedWall                    string  `json:"blockedWall"`
	FinishCalls                    int     `json:"finishCalls"`
	FinishWall                     string  `json:"finishWall"`
	FinishCPU                      string  `json:"finishCpu"`
	UserMemoryReservation          string  `json:"userMemoryReservation"`
	RevocableMemoryReservation     string  `json:"revocableMemoryReservation"`
	SystemMemoryReservation        string  `json:"systemMemoryReservation"`
	PeakUserMemoryReservation      string  `json:"peakUserMemoryReservation"`
	PeakSystemMemoryReservation    string  `json:"peakSystemMemoryReservation"`
	PeakRevocableMemoryReservation string  `json:"peakRevocableMemoryReservation"`
	PeakTotalMemoryReservation     string  `json:"peakTotalMemoryReservation"`
	SpilledDataSize                string  `json:"spilledDataSize"`
}
type StageStats struct {
	SchedulingComplete             time.Time            `json:"schedulingComplete"`
	GetSplitDistribution           GetSplitDistribution `json:"getSplitDistribution"`
	TotalTasks                     int                  `json:"totalTasks"`
	RunningTasks                   int                  `json:"runningTasks"`
	CompletedTasks                 int                  `json:"completedTasks"`
	TotalDrivers                   int                  `json:"totalDrivers"`
	QueuedDrivers                  int                  `json:"queuedDrivers"`
	RunningDrivers                 int                  `json:"runningDrivers"`
	BlockedDrivers                 int                  `json:"blockedDrivers"`
	CompletedDrivers               int                  `json:"completedDrivers"`
	CumulativeUserMemory           interface{}          `json:"cumulativeUserMemory"`
	UserMemoryReservation          string               `json:"userMemoryReservation"`
	RevocableMemoryReservation     string               `json:"revocableMemoryReservation"`
	TotalMemoryReservation         string               `json:"totalMemoryReservation"`
	PeakUserMemoryReservation      string               `json:"peakUserMemoryReservation"`
	PeakRevocableMemoryReservation string               `json:"peakRevocableMemoryReservation"`
	TotalScheduledTime             string               `json:"totalScheduledTime"`
	TotalCPUTime                   string               `json:"totalCpuTime"`
	TotalBlockedTime               string               `json:"totalBlockedTime"`
	FullyBlocked                   bool                 `json:"fullyBlocked"`
	BlockedReasons                 []interface{}        `json:"blockedReasons"`
	PhysicalInputDataSize          string               `json:"physicalInputDataSize"`
	PhysicalInputPositions         int                  `json:"physicalInputPositions"`
	PhysicalInputReadTime          string               `json:"physicalInputReadTime"`
	InternalNetworkInputDataSize   string               `json:"internalNetworkInputDataSize"`
	InternalNetworkInputPositions  int                  `json:"internalNetworkInputPositions"`
	RawInputDataSize               string               `json:"rawInputDataSize"`
	RawInputPositions              int                  `json:"rawInputPositions"`
	ProcessedInputDataSize         string               `json:"processedInputDataSize"`
	ProcessedInputPositions        int                  `json:"processedInputPositions"`
	BufferedDataSize               string               `json:"bufferedDataSize"`
	OutputDataSize                 string               `json:"outputDataSize"`
	OutputPositions                int                  `json:"outputPositions"`
	PhysicalWrittenDataSize        string               `json:"physicalWrittenDataSize"`
	GcInfo                         GcInfo               `json:"gcInfo"`
	OperatorSummaries              []OperatorSummaries  `json:"operatorSummaries"`
}
type TaskStatus struct {
	TaskID                     string        `json:"taskId"`
	TaskInstanceID             string        `json:"taskInstanceId"`
	Version                    int           `json:"version"`
	State                      string        `json:"state"`
	Self                       string        `json:"self"`
	NodeID                     string        `json:"nodeId"`
	CompletedDriverGroups      []interface{} `json:"completedDriverGroups"`
	Failures                   []interface{} `json:"failures"`
	QueuedPartitionedDrivers   int           `json:"queuedPartitionedDrivers"`
	RunningPartitionedDrivers  int           `json:"runningPartitionedDrivers"`
	OutputBufferOverutilized   bool          `json:"outputBufferOverutilized"`
	PhysicalWrittenDataSize    string        `json:"physicalWrittenDataSize"`
	MemoryReservation          string        `json:"memoryReservation"`
	SystemMemoryReservation    string        `json:"systemMemoryReservation"`
	RevocableMemoryReservation string        `json:"revocableMemoryReservation"`
	FullGcCount                int           `json:"fullGcCount"`
	FullGcTime                 string        `json:"fullGcTime"`
	DynamicFiltersVersion      int           `json:"dynamicFiltersVersion"`
}
type OutputBuffers struct {
	Type               string        `json:"type"`
	State              string        `json:"state"`
	CanAddBuffers      bool          `json:"canAddBuffers"`
	CanAddPages        bool          `json:"canAddPages"`
	TotalBufferedBytes int           `json:"totalBufferedBytes"`
	TotalBufferedPages int           `json:"totalBufferedPages"`
	TotalRowsSent      int           `json:"totalRowsSent"`
	TotalPagesSent     int           `json:"totalPagesSent"`
	Buffers            []interface{} `json:"buffers"`
}
type QueuedTime struct {
	Count interface{} `json:"count"`
	Total interface{} `json:"total"`
	P01   interface{} `json:"p01"`
	P05   interface{} `json:"p05"`
	P10   interface{} `json:"p10"`
	P25   interface{} `json:"p25"`
	P50   interface{} `json:"p50"`
	P75   interface{} `json:"p75"`
	P90   interface{} `json:"p90"`
	P95   interface{} `json:"p95"`
	P99   interface{} `json:"p99"`
	Min   interface{} `json:"min"`
	Max   interface{} `json:"max"`
	Avg   interface{} `json:"avg"`
}
type ElapsedTime struct {
	Count float64     `json:"count"`
	Total interface{} `json:"total"`
	P01   interface{} `json:"p01"`
	P05   interface{} `json:"p05"`
	P10   interface{} `json:"p10"`
	P25   interface{} `json:"p25"`
	P50   interface{} `json:"p50"`
	P75   interface{} `json:"p75"`
	P90   interface{} `json:"p90"`
	P95   interface{} `json:"p95"`
	P99   interface{} `json:"p99"`
	Min   interface{} `json:"min"`
	Max   interface{} `json:"max"`
	Avg   interface{} `json:"avg"`
}
type Pipelines struct {
	PipelineID                    int                 `json:"pipelineId"`
	FirstStartTime                time.Time           `json:"firstStartTime"`
	LastStartTime                 time.Time           `json:"lastStartTime"`
	LastEndTime                   time.Time           `json:"lastEndTime"`
	InputPipeline                 bool                `json:"inputPipeline"`
	OutputPipeline                bool                `json:"outputPipeline"`
	TotalDrivers                  int                 `json:"totalDrivers"`
	QueuedDrivers                 int                 `json:"queuedDrivers"`
	QueuedPartitionedDrivers      int                 `json:"queuedPartitionedDrivers"`
	RunningDrivers                int                 `json:"runningDrivers"`
	RunningPartitionedDrivers     int                 `json:"runningPartitionedDrivers"`
	BlockedDrivers                int                 `json:"blockedDrivers"`
	CompletedDrivers              int                 `json:"completedDrivers"`
	UserMemoryReservation         string              `json:"userMemoryReservation"`
	RevocableMemoryReservation    string              `json:"revocableMemoryReservation"`
	SystemMemoryReservation       string              `json:"systemMemoryReservation"`
	QueuedTime                    QueuedTime          `json:"queuedTime"`
	ElapsedTime                   ElapsedTime         `json:"elapsedTime"`
	TotalScheduledTime            string              `json:"totalScheduledTime"`
	TotalCPUTime                  string              `json:"totalCpuTime"`
	TotalBlockedTime              string              `json:"totalBlockedTime"`
	FullyBlocked                  bool                `json:"fullyBlocked"`
	BlockedReasons                []interface{}       `json:"blockedReasons"`
	PhysicalInputDataSize         string              `json:"physicalInputDataSize"`
	PhysicalInputPositions        int                 `json:"physicalInputPositions"`
	PhysicalInputReadTime         string              `json:"physicalInputReadTime"`
	InternalNetworkInputDataSize  string              `json:"internalNetworkInputDataSize"`
	InternalNetworkInputPositions int                 `json:"internalNetworkInputPositions"`
	RawInputDataSize              string              `json:"rawInputDataSize"`
	RawInputPositions             int                 `json:"rawInputPositions"`
	ProcessedInputDataSize        string              `json:"processedInputDataSize"`
	ProcessedInputPositions       int                 `json:"processedInputPositions"`
	OutputDataSize                string              `json:"outputDataSize"`
	OutputPositions               int                 `json:"outputPositions"`
	PhysicalWrittenDataSize       string              `json:"physicalWrittenDataSize"`
	OperatorSummaries             []OperatorSummaries `json:"operatorSummaries"`
	Drivers                       []interface{}       `json:"drivers"`
}
type Stats struct {
	CreateTime                    time.Time     `json:"createTime"`
	FirstStartTime                time.Time     `json:"firstStartTime"`
	LastStartTime                 time.Time     `json:"lastStartTime"`
	LastEndTime                   time.Time     `json:"lastEndTime"`
	EndTime                       time.Time     `json:"endTime"`
	ElapsedTime                   string        `json:"elapsedTime"`
	QueuedTime                    string        `json:"queuedTime"`
	TotalDrivers                  int           `json:"totalDrivers"`
	QueuedDrivers                 int           `json:"queuedDrivers"`
	QueuedPartitionedDrivers      int           `json:"queuedPartitionedDrivers"`
	RunningDrivers                int           `json:"runningDrivers"`
	RunningPartitionedDrivers     int           `json:"runningPartitionedDrivers"`
	BlockedDrivers                int           `json:"blockedDrivers"`
	CompletedDrivers              int           `json:"completedDrivers"`
	CumulativeUserMemory          interface{}   `json:"cumulativeUserMemory"`
	UserMemoryReservation         string        `json:"userMemoryReservation"`
	RevocableMemoryReservation    string        `json:"revocableMemoryReservation"`
	SystemMemoryReservation       string        `json:"systemMemoryReservation"`
	TotalScheduledTime            string        `json:"totalScheduledTime"`
	TotalCPUTime                  string        `json:"totalCpuTime"`
	TotalBlockedTime              string        `json:"totalBlockedTime"`
	FullyBlocked                  bool          `json:"fullyBlocked"`
	BlockedReasons                []interface{} `json:"blockedReasons"`
	PhysicalInputDataSize         string        `json:"physicalInputDataSize"`
	PhysicalInputPositions        int           `json:"physicalInputPositions"`
	PhysicalInputReadTime         string        `json:"physicalInputReadTime"`
	InternalNetworkInputDataSize  string        `json:"internalNetworkInputDataSize"`
	InternalNetworkInputPositions int           `json:"internalNetworkInputPositions"`
	RawInputDataSize              string        `json:"rawInputDataSize"`
	RawInputPositions             int           `json:"rawInputPositions"`
	ProcessedInputDataSize        string        `json:"processedInputDataSize"`
	ProcessedInputPositions       int           `json:"processedInputPositions"`
	OutputDataSize                string        `json:"outputDataSize"`
	OutputPositions               int           `json:"outputPositions"`
	PhysicalWrittenDataSize       string        `json:"physicalWrittenDataSize"`
	FullGcCount                   int           `json:"fullGcCount"`
	FullGcTime                    string        `json:"fullGcTime"`
	Pipelines                     []Pipelines   `json:"pipelines"`
}
type Tasks struct {
	TaskStatus    TaskStatus    `json:"taskStatus"`
	LastHeartbeat time.Time     `json:"lastHeartbeat"`
	OutputBuffers OutputBuffers `json:"outputBuffers"`
	NoMoreSplits  []string      `json:"noMoreSplits"`
	Stats         Stats         `json:"stats"`
	NeedsPlan     bool          `json:"needsPlan"`
}
type Tables struct {
}
type SubStages struct {
	StageID    string      `json:"stageId"`
	State      string      `json:"state"`
	Types      []string    `json:"types"`
	StageStats StageStats  `json:"stageStats"`
	Tasks      []Tasks     `json:"tasks"`
	SubStages  []SubStages `json:"subStages"`
	Tables     Tables      `json:"tables"`
}
type OutputStage struct {
	StageID    string      `json:"stageId"`
	State      string      `json:"state"`
	Types      []string    `json:"types"`
	StageStats StageStats  `json:"stageStats"`
	Tasks      []Tasks     `json:"tasks"`
	SubStages  []SubStages `json:"subStages"`
	Tables     Tables      `json:"tables"`
}
type ConnectorInfo struct {
	PartitionIds []string `json:"partitionIds"`
	Truncated    bool     `json:"truncated"`
}
type Routines struct {
	Routine       string `json:"routine"`
	Authorization string `json:"authorization"`
}
