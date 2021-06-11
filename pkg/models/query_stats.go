package models

import "time"

type QueryStats struct {
	QueryId string `json:"queryId"`
	Session struct {
		QueryId                  string        `json:"queryId"`
		TransactionId            string        `json:"transactionId"`
		ClientTransactionSupport bool          `json:"clientTransactionSupport"`
		User                     string        `json:"user"`
		Groups                   []interface{} `json:"groups"`
		Principal                string        `json:"principal"`
		Source                   string        `json:"source"`
		Catalog                  string        `json:"catalog"`
		Schema                   string        `json:"schema"`
		Path                     struct {
		} `json:"path"`
		TimeZoneKey        int           `json:"timeZoneKey"`
		Locale             string        `json:"locale"`
		RemoteUserAddress  string        `json:"remoteUserAddress"`
		UserAgent          string        `json:"userAgent"`
		ClientTags         []interface{} `json:"clientTags"`
		ClientCapabilities []string      `json:"clientCapabilities"`
		ResourceEstimates  struct {
		} `json:"resourceEstimates"`
		Start            time.Time `json:"start"`
		SystemProperties struct {
		} `json:"systemProperties"`
		CatalogProperties struct {
		} `json:"catalogProperties"`
		UnprocessedCatalogProperties struct {
		} `json:"unprocessedCatalogProperties"`
		Roles struct {
		} `json:"roles"`
		PreparedStatements struct {
		} `json:"preparedStatements"`
		ProtocolName string `json:"protocolName"`
	} `json:"session"`
	State      string   `json:"state"`
	MemoryPool string   `json:"memoryPool"`
	Scheduled  bool     `json:"scheduled"`
	Self       string   `json:"self"`
	FieldNames []string `json:"fieldNames"`
	Query      string   `json:"query"`
	QueryStats struct {
		CreateTime                        time.Time     `json:"createTime"`
		ExecutionStartTime                time.Time     `json:"executionStartTime"`
		LastHeartbeat                     time.Time     `json:"lastHeartbeat"`
		EndTime                           time.Time     `json:"endTime"`
		ElapsedTime                       string        `json:"elapsedTime"`
		QueuedTime                        string        `json:"queuedTime"`
		ResourceWaitingTime               string        `json:"resourceWaitingTime"`
		DispatchingTime                   string        `json:"dispatchingTime"`
		ExecutionTime                     string        `json:"executionTime"`
		AnalysisTime                      string        `json:"analysisTime"`
		PlanningTime                      string        `json:"planningTime"`
		FinishingTime                     string        `json:"finishingTime"`
		TotalTasks                        int           `json:"totalTasks"`
		RunningTasks                      int           `json:"runningTasks"`
		CompletedTasks                    int           `json:"completedTasks"`
		TotalDrivers                      int           `json:"totalDrivers"`
		QueuedDrivers                     int           `json:"queuedDrivers"`
		RunningDrivers                    int           `json:"runningDrivers"`
		BlockedDrivers                    int           `json:"blockedDrivers"`
		CompletedDrivers                  int           `json:"completedDrivers"`
		CumulativeUserMemory              float64       `json:"cumulativeUserMemory"`
		UserMemoryReservation             string        `json:"userMemoryReservation"`
		RevocableMemoryReservation        string        `json:"revocableMemoryReservation"`
		TotalMemoryReservation            string        `json:"totalMemoryReservation"`
		PeakUserMemoryReservation         string        `json:"peakUserMemoryReservation"`
		PeakRevocableMemoryReservation    string        `json:"peakRevocableMemoryReservation"`
		PeakNonRevocableMemoryReservation string        `json:"peakNonRevocableMemoryReservation"`
		PeakTotalMemoryReservation        string        `json:"peakTotalMemoryReservation"`
		PeakTaskUserMemory                string        `json:"peakTaskUserMemory"`
		PeakTaskRevocableMemory           string        `json:"peakTaskRevocableMemory"`
		PeakTaskTotalMemory               string        `json:"peakTaskTotalMemory"`
		Scheduled                         bool          `json:"scheduled"`
		TotalScheduledTime                string        `json:"totalScheduledTime"`
		TotalCpuTime                      string        `json:"totalCpuTime"`
		TotalBlockedTime                  string        `json:"totalBlockedTime"`
		FullyBlocked                      bool          `json:"fullyBlocked"`
		BlockedReasons                    []interface{} `json:"blockedReasons"`
		PhysicalInputDataSize             string        `json:"physicalInputDataSize"`
		PhysicalInputPositions            int           `json:"physicalInputPositions"`
		PhysicalInputReadTime             string        `json:"physicalInputReadTime"`
		InternalNetworkInputDataSize      string        `json:"internalNetworkInputDataSize"`
		InternalNetworkInputPositions     int           `json:"internalNetworkInputPositions"`
		RawInputDataSize                  string        `json:"rawInputDataSize"`
		RawInputPositions                 int           `json:"rawInputPositions"`
		ProcessedInputDataSize            string        `json:"processedInputDataSize"`
		ProcessedInputPositions           int           `json:"processedInputPositions"`
		OutputDataSize                    string        `json:"outputDataSize"`
		OutputPositions                   int           `json:"outputPositions"`
		PhysicalWrittenDataSize           string        `json:"physicalWrittenDataSize"`
		StageGcStatistics                 []struct {
			StageId          int `json:"stageId"`
			Tasks            int `json:"tasks"`
			FullGcTasks      int `json:"fullGcTasks"`
			MinFullGcSec     int `json:"minFullGcSec"`
			MaxFullGcSec     int `json:"maxFullGcSec"`
			TotalFullGcSec   int `json:"totalFullGcSec"`
			AverageFullGcSec int `json:"averageFullGcSec"`
		} `json:"stageGcStatistics"`
		DynamicFiltersStats struct {
			DynamicFilterDomainStats []interface{} `json:"dynamicFilterDomainStats"`
			LazyDynamicFilters       int           `json:"lazyDynamicFilters"`
			ReplicatedDynamicFilters int           `json:"replicatedDynamicFilters"`
			TotalDynamicFilters      int           `json:"totalDynamicFilters"`
			DynamicFiltersCompleted  int           `json:"dynamicFiltersCompleted"`
		} `json:"dynamicFiltersStats"`
		OperatorSummaries []struct {
			StageId                        int     `json:"stageId"`
			PipelineId                     int     `json:"pipelineId"`
			OperatorId                     int     `json:"operatorId"`
			PlanNodeId                     string  `json:"planNodeId"`
			OperatorType                   string  `json:"operatorType"`
			TotalDrivers                   int     `json:"totalDrivers"`
			AddInputCalls                  int     `json:"addInputCalls"`
			AddInputWall                   string  `json:"addInputWall"`
			AddInputCpu                    string  `json:"addInputCpu"`
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
			GetOutputCpu                   string  `json:"getOutputCpu"`
			OutputDataSize                 string  `json:"outputDataSize"`
			OutputPositions                int     `json:"outputPositions"`
			DynamicFilterSplitsProcessed   int     `json:"dynamicFilterSplitsProcessed"`
			PhysicalWrittenDataSize        string  `json:"physicalWrittenDataSize"`
			BlockedWall                    string  `json:"blockedWall"`
			FinishCalls                    int     `json:"finishCalls"`
			FinishWall                     string  `json:"finishWall"`
			FinishCpu                      string  `json:"finishCpu"`
			UserMemoryReservation          string  `json:"userMemoryReservation"`
			RevocableMemoryReservation     string  `json:"revocableMemoryReservation"`
			SystemMemoryReservation        string  `json:"systemMemoryReservation"`
			PeakUserMemoryReservation      string  `json:"peakUserMemoryReservation"`
			PeakSystemMemoryReservation    string  `json:"peakSystemMemoryReservation"`
			PeakRevocableMemoryReservation string  `json:"peakRevocableMemoryReservation"`
			PeakTotalMemoryReservation     string  `json:"peakTotalMemoryReservation"`
			SpilledDataSize                string  `json:"spilledDataSize"`
			Info                           struct {
				Type                        string        `json:"@type"`
				BufferedBytes               int           `json:"bufferedBytes,omitempty"`
				MaxBufferedBytes            int           `json:"maxBufferedBytes,omitempty"`
				AverageBytesPerRequest      int           `json:"averageBytesPerRequest,omitempty"`
				SuccessfulRequestsCount     int           `json:"successfulRequestsCount,omitempty"`
				BufferedPages               int           `json:"bufferedPages,omitempty"`
				NoMoreLocations             bool          `json:"noMoreLocations,omitempty"`
				PageBufferClientStatuses    []interface{} `json:"pageBufferClientStatuses,omitempty"`
				RowsAdded                   int           `json:"rowsAdded,omitempty"`
				PagesAdded                  int           `json:"pagesAdded,omitempty"`
				OutputBufferPeakMemoryUsage int           `json:"outputBufferPeakMemoryUsage,omitempty"`
			} `json:"info,omitempty"`
		} `json:"operatorSummaries"`
		LogicalWrittenDataSize string  `json:"logicalWrittenDataSize"`
		WrittenPositions       int     `json:"writtenPositions"`
		SpilledDataSize        string  `json:"spilledDataSize"`
		ProgressPercentage     float64 `json:"progressPercentage"`
	} `json:"queryStats"`
	SetSessionProperties struct {
	} `json:"setSessionProperties"`
	ResetSessionProperties []interface{} `json:"resetSessionProperties"`
	SetRoles               struct {
	} `json:"setRoles"`
	AddedPreparedStatements struct {
	} `json:"addedPreparedStatements"`
	DeallocatedPreparedStatements []interface{} `json:"deallocatedPreparedStatements"`
	ClearTransactionId            bool          `json:"clearTransactionId"`
	OutputStage                   struct {
		StageId string `json:"stageId"`
		State   string `json:"state"`
		Plan    struct {
			Id   string `json:"id"`
			Root struct {
				Type   string `json:"@type"`
				Id     string `json:"id"`
				Source struct {
					Type              string   `json:"@type"`
					Id                string   `json:"id"`
					SourceFragmentIds []string `json:"sourceFragmentIds"`
					Outputs           []string `json:"outputs"`
					ExchangeType      string   `json:"exchangeType"`
				} `json:"source"`
				Columns []string `json:"columns"`
				Outputs []string `json:"outputs"`
			} `json:"root"`
			Symbols struct {
				Sum string `json:"sum"`
				Avg string `json:"avg"`
			} `json:"symbols"`
			Partitioning struct {
				ConnectorHandle struct {
					Type         string `json:"@type"`
					Partitioning string `json:"partitioning"`
					Function     string `json:"function"`
				} `json:"connectorHandle"`
			} `json:"partitioning"`
			PartitionedSources []interface{} `json:"partitionedSources"`
			PartitioningScheme struct {
				Partitioning struct {
					Handle struct {
						ConnectorHandle struct {
							Type         string `json:"@type"`
							Partitioning string `json:"partitioning"`
							Function     string `json:"function"`
						} `json:"connectorHandle"`
					} `json:"handle"`
					Arguments []interface{} `json:"arguments"`
				} `json:"partitioning"`
				OutputLayout         []string `json:"outputLayout"`
				ReplicateNullsAndAny bool     `json:"replicateNullsAndAny"`
				BucketToPartition    []int    `json:"bucketToPartition"`
			} `json:"partitioningScheme"`
			StageExecutionDescriptor struct {
				Strategy                  string        `json:"strategy"`
				GroupedExecutionScanNodes []interface{} `json:"groupedExecutionScanNodes"`
			} `json:"stageExecutionDescriptor"`
			StatsAndCosts struct {
				Stats struct {
				} `json:"stats"`
				Costs struct {
				} `json:"costs"`
			} `json:"statsAndCosts"`
			JsonRepresentation string `json:"jsonRepresentation"`
		} `json:"plan"`
		Types      []string `json:"types"`
		StageStats struct {
			SchedulingComplete   time.Time `json:"schedulingComplete"`
			GetSplitDistribution struct {
				Count float64 `json:"count"`
				Total float64 `json:"total"`
				P01   string  `json:"p01"`
				P05   string  `json:"p05"`
				P10   string  `json:"p10"`
				P25   string  `json:"p25"`
				P50   string  `json:"p50"`
				P75   string  `json:"p75"`
				P90   string  `json:"p90"`
				P95   string  `json:"p95"`
				P99   string  `json:"p99"`
				Min   string  `json:"min"`
				Max   string  `json:"max"`
				Avg   string  `json:"avg"`
			} `json:"getSplitDistribution"`
			TotalTasks                     int           `json:"totalTasks"`
			RunningTasks                   int           `json:"runningTasks"`
			CompletedTasks                 int           `json:"completedTasks"`
			TotalDrivers                   int           `json:"totalDrivers"`
			QueuedDrivers                  int           `json:"queuedDrivers"`
			RunningDrivers                 int           `json:"runningDrivers"`
			BlockedDrivers                 int           `json:"blockedDrivers"`
			CompletedDrivers               int           `json:"completedDrivers"`
			CumulativeUserMemory           float64       `json:"cumulativeUserMemory"`
			UserMemoryReservation          string        `json:"userMemoryReservation"`
			RevocableMemoryReservation     string        `json:"revocableMemoryReservation"`
			TotalMemoryReservation         string        `json:"totalMemoryReservation"`
			PeakUserMemoryReservation      string        `json:"peakUserMemoryReservation"`
			PeakRevocableMemoryReservation string        `json:"peakRevocableMemoryReservation"`
			TotalScheduledTime             string        `json:"totalScheduledTime"`
			TotalCpuTime                   string        `json:"totalCpuTime"`
			TotalBlockedTime               string        `json:"totalBlockedTime"`
			FullyBlocked                   bool          `json:"fullyBlocked"`
			BlockedReasons                 []interface{} `json:"blockedReasons"`
			PhysicalInputDataSize          string        `json:"physicalInputDataSize"`
			PhysicalInputPositions         int           `json:"physicalInputPositions"`
			PhysicalInputReadTime          string        `json:"physicalInputReadTime"`
			InternalNetworkInputDataSize   string        `json:"internalNetworkInputDataSize"`
			InternalNetworkInputPositions  int           `json:"internalNetworkInputPositions"`
			RawInputDataSize               string        `json:"rawInputDataSize"`
			RawInputPositions              int           `json:"rawInputPositions"`
			ProcessedInputDataSize         string        `json:"processedInputDataSize"`
			ProcessedInputPositions        int           `json:"processedInputPositions"`
			BufferedDataSize               string        `json:"bufferedDataSize"`
			OutputDataSize                 string        `json:"outputDataSize"`
			OutputPositions                int           `json:"outputPositions"`
			PhysicalWrittenDataSize        string        `json:"physicalWrittenDataSize"`
			GcInfo                         struct {
				StageId          int `json:"stageId"`
				Tasks            int `json:"tasks"`
				FullGcTasks      int `json:"fullGcTasks"`
				MinFullGcSec     int `json:"minFullGcSec"`
				MaxFullGcSec     int `json:"maxFullGcSec"`
				TotalFullGcSec   int `json:"totalFullGcSec"`
				AverageFullGcSec int `json:"averageFullGcSec"`
			} `json:"gcInfo"`
			OperatorSummaries []struct {
				StageId                        int     `json:"stageId"`
				PipelineId                     int     `json:"pipelineId"`
				OperatorId                     int     `json:"operatorId"`
				PlanNodeId                     string  `json:"planNodeId"`
				OperatorType                   string  `json:"operatorType"`
				TotalDrivers                   int     `json:"totalDrivers"`
				AddInputCalls                  int     `json:"addInputCalls"`
				AddInputWall                   string  `json:"addInputWall"`
				AddInputCpu                    string  `json:"addInputCpu"`
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
				GetOutputCpu                   string  `json:"getOutputCpu"`
				OutputDataSize                 string  `json:"outputDataSize"`
				OutputPositions                int     `json:"outputPositions"`
				DynamicFilterSplitsProcessed   int     `json:"dynamicFilterSplitsProcessed"`
				PhysicalWrittenDataSize        string  `json:"physicalWrittenDataSize"`
				BlockedWall                    string  `json:"blockedWall"`
				FinishCalls                    int     `json:"finishCalls"`
				FinishWall                     string  `json:"finishWall"`
				FinishCpu                      string  `json:"finishCpu"`
				UserMemoryReservation          string  `json:"userMemoryReservation"`
				RevocableMemoryReservation     string  `json:"revocableMemoryReservation"`
				SystemMemoryReservation        string  `json:"systemMemoryReservation"`
				PeakUserMemoryReservation      string  `json:"peakUserMemoryReservation"`
				PeakSystemMemoryReservation    string  `json:"peakSystemMemoryReservation"`
				PeakRevocableMemoryReservation string  `json:"peakRevocableMemoryReservation"`
				PeakTotalMemoryReservation     string  `json:"peakTotalMemoryReservation"`
				SpilledDataSize                string  `json:"spilledDataSize"`
				Info                           struct {
					Type                     string        `json:"@type"`
					BufferedBytes            int           `json:"bufferedBytes"`
					MaxBufferedBytes         int           `json:"maxBufferedBytes"`
					AverageBytesPerRequest   int           `json:"averageBytesPerRequest"`
					SuccessfulRequestsCount  int           `json:"successfulRequestsCount"`
					BufferedPages            int           `json:"bufferedPages"`
					NoMoreLocations          bool          `json:"noMoreLocations"`
					PageBufferClientStatuses []interface{} `json:"pageBufferClientStatuses"`
				} `json:"info,omitempty"`
			} `json:"operatorSummaries"`
		} `json:"stageStats"`
		Tasks []struct {
			TaskStatus struct {
				TaskId                     string        `json:"taskId"`
				TaskInstanceId             string        `json:"taskInstanceId"`
				Version                    int           `json:"version"`
				State                      string        `json:"state"`
				Self                       string        `json:"self"`
				NodeId                     string        `json:"nodeId"`
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
			} `json:"taskStatus"`
			LastHeartbeat time.Time `json:"lastHeartbeat"`
			OutputBuffers struct {
				Type               string        `json:"type"`
				State              string        `json:"state"`
				CanAddBuffers      bool          `json:"canAddBuffers"`
				CanAddPages        bool          `json:"canAddPages"`
				TotalBufferedBytes int           `json:"totalBufferedBytes"`
				TotalBufferedPages int           `json:"totalBufferedPages"`
				TotalRowsSent      int           `json:"totalRowsSent"`
				TotalPagesSent     int           `json:"totalPagesSent"`
				Buffers            []interface{} `json:"buffers"`
			} `json:"outputBuffers"`
			NoMoreSplits []string `json:"noMoreSplits"`
			Stats        struct {
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
				CumulativeUserMemory          float64       `json:"cumulativeUserMemory"`
				UserMemoryReservation         string        `json:"userMemoryReservation"`
				RevocableMemoryReservation    string        `json:"revocableMemoryReservation"`
				SystemMemoryReservation       string        `json:"systemMemoryReservation"`
				TotalScheduledTime            string        `json:"totalScheduledTime"`
				TotalCpuTime                  string        `json:"totalCpuTime"`
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
				Pipelines                     []struct {
					PipelineId                 int       `json:"pipelineId"`
					FirstStartTime             time.Time `json:"firstStartTime"`
					LastStartTime              time.Time `json:"lastStartTime"`
					LastEndTime                time.Time `json:"lastEndTime"`
					InputPipeline              bool      `json:"inputPipeline"`
					OutputPipeline             bool      `json:"outputPipeline"`
					TotalDrivers               int       `json:"totalDrivers"`
					QueuedDrivers              int       `json:"queuedDrivers"`
					QueuedPartitionedDrivers   int       `json:"queuedPartitionedDrivers"`
					RunningDrivers             int       `json:"runningDrivers"`
					RunningPartitionedDrivers  int       `json:"runningPartitionedDrivers"`
					BlockedDrivers             int       `json:"blockedDrivers"`
					CompletedDrivers           int       `json:"completedDrivers"`
					UserMemoryReservation      string    `json:"userMemoryReservation"`
					RevocableMemoryReservation string    `json:"revocableMemoryReservation"`
					SystemMemoryReservation    string    `json:"systemMemoryReservation"`
					QueuedTime                 struct {
						Count float64 `json:"count"`
						Total float64 `json:"total"`
						P01   float64 `json:"p01"`
						P05   float64 `json:"p05"`
						P10   float64 `json:"p10"`
						P25   float64 `json:"p25"`
						P50   float64 `json:"p50"`
						P75   float64 `json:"p75"`
						P90   float64 `json:"p90"`
						P95   float64 `json:"p95"`
						P99   float64 `json:"p99"`
						Min   float64 `json:"min"`
						Max   float64 `json:"max"`
						Avg   float64 `json:"avg"`
					} `json:"queuedTime"`
					ElapsedTime struct {
						Count float64 `json:"count"`
						Total float64 `json:"total"`
						P01   float64 `json:"p01"`
						P05   float64 `json:"p05"`
						P10   float64 `json:"p10"`
						P25   float64 `json:"p25"`
						P50   float64 `json:"p50"`
						P75   float64 `json:"p75"`
						P90   float64 `json:"p90"`
						P95   float64 `json:"p95"`
						P99   float64 `json:"p99"`
						Min   float64 `json:"min"`
						Max   float64 `json:"max"`
						Avg   float64 `json:"avg"`
					} `json:"elapsedTime"`
					TotalScheduledTime            string        `json:"totalScheduledTime"`
					TotalCpuTime                  string        `json:"totalCpuTime"`
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
					OperatorSummaries             []struct {
						StageId                        int     `json:"stageId"`
						PipelineId                     int     `json:"pipelineId"`
						OperatorId                     int     `json:"operatorId"`
						PlanNodeId                     string  `json:"planNodeId"`
						OperatorType                   string  `json:"operatorType"`
						TotalDrivers                   int     `json:"totalDrivers"`
						AddInputCalls                  int     `json:"addInputCalls"`
						AddInputWall                   string  `json:"addInputWall"`
						AddInputCpu                    string  `json:"addInputCpu"`
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
						GetOutputCpu                   string  `json:"getOutputCpu"`
						OutputDataSize                 string  `json:"outputDataSize"`
						OutputPositions                int     `json:"outputPositions"`
						DynamicFilterSplitsProcessed   int     `json:"dynamicFilterSplitsProcessed"`
						PhysicalWrittenDataSize        string  `json:"physicalWrittenDataSize"`
						BlockedWall                    string  `json:"blockedWall"`
						FinishCalls                    int     `json:"finishCalls"`
						FinishWall                     string  `json:"finishWall"`
						FinishCpu                      string  `json:"finishCpu"`
						UserMemoryReservation          string  `json:"userMemoryReservation"`
						RevocableMemoryReservation     string  `json:"revocableMemoryReservation"`
						SystemMemoryReservation        string  `json:"systemMemoryReservation"`
						PeakUserMemoryReservation      string  `json:"peakUserMemoryReservation"`
						PeakSystemMemoryReservation    string  `json:"peakSystemMemoryReservation"`
						PeakRevocableMemoryReservation string  `json:"peakRevocableMemoryReservation"`
						PeakTotalMemoryReservation     string  `json:"peakTotalMemoryReservation"`
						SpilledDataSize                string  `json:"spilledDataSize"`
						Info                           struct {
							Type                     string        `json:"@type"`
							BufferedBytes            int           `json:"bufferedBytes"`
							MaxBufferedBytes         int           `json:"maxBufferedBytes"`
							AverageBytesPerRequest   int           `json:"averageBytesPerRequest"`
							SuccessfulRequestsCount  int           `json:"successfulRequestsCount"`
							BufferedPages            int           `json:"bufferedPages"`
							NoMoreLocations          bool          `json:"noMoreLocations"`
							PageBufferClientStatuses []interface{} `json:"pageBufferClientStatuses"`
						} `json:"info,omitempty"`
					} `json:"operatorSummaries"`
					Drivers []interface{} `json:"drivers"`
				} `json:"pipelines"`
			} `json:"stats"`
			NeedsPlan bool `json:"needsPlan"`
		} `json:"tasks"`
		SubStages []struct {
			StageId string `json:"stageId"`
			State   string `json:"state"`
			Plan    struct {
				Id   string `json:"id"`
				Root struct {
					Type   string `json:"@type"`
					Id     string `json:"id"`
					Source struct {
						Type   string `json:"@type"`
						Id     string `json:"id"`
						Source struct {
							Type               string `json:"@type"`
							Id                 string `json:"id"`
							Type1              string `json:"type"`
							Scope              string `json:"scope"`
							PartitioningScheme struct {
								Partitioning struct {
									Handle struct {
										ConnectorHandle struct {
											Type         string `json:"@type"`
											Partitioning string `json:"partitioning"`
											Function     string `json:"function"`
										} `json:"connectorHandle"`
									} `json:"handle"`
									Arguments []struct {
										Expression string `json:"expression"`
									} `json:"arguments"`
								} `json:"partitioning"`
								OutputLayout         []string `json:"outputLayout"`
								HashColumn           string   `json:"hashColumn"`
								ReplicateNullsAndAny bool     `json:"replicateNullsAndAny"`
							} `json:"partitioningScheme"`
							Sources []struct {
								Type              string   `json:"@type"`
								Id                string   `json:"id"`
								SourceFragmentIds []string `json:"sourceFragmentIds"`
								Outputs           []string `json:"outputs"`
								ExchangeType      string   `json:"exchangeType"`
							} `json:"sources"`
							Inputs [][]string `json:"inputs"`
						} `json:"source"`
						Aggregations struct {
							Avg struct {
								ResolvedFunction struct {
									Signature struct {
										Name          string   `json:"name"`
										ReturnType    string   `json:"returnType"`
										ArgumentTypes []string `json:"argumentTypes"`
									} `json:"signature"`
									Id               string `json:"id"`
									TypeDependencies struct {
									} `json:"typeDependencies"`
									FunctionDependencies []interface{} `json:"functionDependencies"`
								} `json:"resolvedFunction"`
								Arguments []string `json:"arguments"`
								Distinct  bool     `json:"distinct"`
							} `json:"avg"`
							Sum struct {
								ResolvedFunction struct {
									Signature struct {
										Name          string   `json:"name"`
										ReturnType    string   `json:"returnType"`
										ArgumentTypes []string `json:"argumentTypes"`
									} `json:"signature"`
									Id               string `json:"id"`
									TypeDependencies struct {
									} `json:"typeDependencies"`
									FunctionDependencies []interface{} `json:"functionDependencies"`
								} `json:"resolvedFunction"`
								Arguments []string `json:"arguments"`
								Distinct  bool     `json:"distinct"`
							} `json:"sum"`
						} `json:"aggregations"`
						GroupingSets struct {
							GroupingKeys       []string      `json:"groupingKeys"`
							GroupingSetCount   int           `json:"groupingSetCount"`
							GlobalGroupingSets []interface{} `json:"globalGroupingSets"`
						} `json:"groupingSets"`
						PreGroupedSymbols []interface{} `json:"preGroupedSymbols"`
						Step              string        `json:"step"`
						HashSymbol        string        `json:"hashSymbol"`
					} `json:"source"`
					Assignments struct {
						Assignments struct {
							Avg string `json:"avg"`
							Sum string `json:"sum"`
						} `json:"assignments"`
					} `json:"assignments"`
				} `json:"root"`
				Symbols struct {
					Sum         string `json:"sum"`
					Avg1        string `json:"avg_1"`
					Hashvalue3  string `json:"$hashvalue_3"`
					Avg         string `json:"avg"`
					UserCountry string `json:"user_country"`
					Sum2        string `json:"sum_2"`
					Hashvalue   string `json:"$hashvalue"`
				} `json:"symbols"`
				Partitioning struct {
					ConnectorHandle struct {
						Type         string `json:"@type"`
						Partitioning string `json:"partitioning"`
						Function     string `json:"function"`
					} `json:"connectorHandle"`
				} `json:"partitioning"`
				PartitionedSources []interface{} `json:"partitionedSources"`
				PartitioningScheme struct {
					Partitioning struct {
						Handle struct {
							ConnectorHandle struct {
								Type         string `json:"@type"`
								Partitioning string `json:"partitioning"`
								Function     string `json:"function"`
							} `json:"connectorHandle"`
						} `json:"handle"`
						Arguments []interface{} `json:"arguments"`
					} `json:"partitioning"`
					OutputLayout         []string `json:"outputLayout"`
					ReplicateNullsAndAny bool     `json:"replicateNullsAndAny"`
					BucketToPartition    []int    `json:"bucketToPartition"`
				} `json:"partitioningScheme"`
				StageExecutionDescriptor struct {
					Strategy                  string        `json:"strategy"`
					GroupedExecutionScanNodes []interface{} `json:"groupedExecutionScanNodes"`
				} `json:"stageExecutionDescriptor"`
				StatsAndCosts struct {
					Stats struct {
					} `json:"stats"`
					Costs struct {
					} `json:"costs"`
				} `json:"statsAndCosts"`
				JsonRepresentation string `json:"jsonRepresentation"`
			} `json:"plan"`
			Types      []string `json:"types"`
			StageStats struct {
				SchedulingComplete   time.Time `json:"schedulingComplete"`
				GetSplitDistribution struct {
					Count float64 `json:"count"`
					Total float64 `json:"total"`
					P01   string  `json:"p01"`
					P05   string  `json:"p05"`
					P10   string  `json:"p10"`
					P25   string  `json:"p25"`
					P50   string  `json:"p50"`
					P75   string  `json:"p75"`
					P90   string  `json:"p90"`
					P95   string  `json:"p95"`
					P99   string  `json:"p99"`
					Min   string  `json:"min"`
					Max   string  `json:"max"`
					Avg   string  `json:"avg"`
				} `json:"getSplitDistribution"`
				TotalTasks                     int           `json:"totalTasks"`
				RunningTasks                   int           `json:"runningTasks"`
				CompletedTasks                 int           `json:"completedTasks"`
				TotalDrivers                   int           `json:"totalDrivers"`
				QueuedDrivers                  int           `json:"queuedDrivers"`
				RunningDrivers                 int           `json:"runningDrivers"`
				BlockedDrivers                 int           `json:"blockedDrivers"`
				CompletedDrivers               int           `json:"completedDrivers"`
				CumulativeUserMemory           float64       `json:"cumulativeUserMemory"`
				UserMemoryReservation          string        `json:"userMemoryReservation"`
				RevocableMemoryReservation     string        `json:"revocableMemoryReservation"`
				TotalMemoryReservation         string        `json:"totalMemoryReservation"`
				PeakUserMemoryReservation      string        `json:"peakUserMemoryReservation"`
				PeakRevocableMemoryReservation string        `json:"peakRevocableMemoryReservation"`
				TotalScheduledTime             string        `json:"totalScheduledTime"`
				TotalCpuTime                   string        `json:"totalCpuTime"`
				TotalBlockedTime               string        `json:"totalBlockedTime"`
				FullyBlocked                   bool          `json:"fullyBlocked"`
				BlockedReasons                 []interface{} `json:"blockedReasons"`
				PhysicalInputDataSize          string        `json:"physicalInputDataSize"`
				PhysicalInputPositions         int           `json:"physicalInputPositions"`
				PhysicalInputReadTime          string        `json:"physicalInputReadTime"`
				InternalNetworkInputDataSize   string        `json:"internalNetworkInputDataSize"`
				InternalNetworkInputPositions  int           `json:"internalNetworkInputPositions"`
				RawInputDataSize               string        `json:"rawInputDataSize"`
				RawInputPositions              int           `json:"rawInputPositions"`
				ProcessedInputDataSize         string        `json:"processedInputDataSize"`
				ProcessedInputPositions        int           `json:"processedInputPositions"`
				BufferedDataSize               string        `json:"bufferedDataSize"`
				OutputDataSize                 string        `json:"outputDataSize"`
				OutputPositions                int           `json:"outputPositions"`
				PhysicalWrittenDataSize        string        `json:"physicalWrittenDataSize"`
				GcInfo                         struct {
					StageId          int `json:"stageId"`
					Tasks            int `json:"tasks"`
					FullGcTasks      int `json:"fullGcTasks"`
					MinFullGcSec     int `json:"minFullGcSec"`
					MaxFullGcSec     int `json:"maxFullGcSec"`
					TotalFullGcSec   int `json:"totalFullGcSec"`
					AverageFullGcSec int `json:"averageFullGcSec"`
				} `json:"gcInfo"`
				OperatorSummaries []struct {
					StageId                        int     `json:"stageId"`
					PipelineId                     int     `json:"pipelineId"`
					OperatorId                     int     `json:"operatorId"`
					PlanNodeId                     string  `json:"planNodeId"`
					OperatorType                   string  `json:"operatorType"`
					TotalDrivers                   int     `json:"totalDrivers"`
					AddInputCalls                  int     `json:"addInputCalls"`
					AddInputWall                   string  `json:"addInputWall"`
					AddInputCpu                    string  `json:"addInputCpu"`
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
					GetOutputCpu                   string  `json:"getOutputCpu"`
					OutputDataSize                 string  `json:"outputDataSize"`
					OutputPositions                int     `json:"outputPositions"`
					DynamicFilterSplitsProcessed   int     `json:"dynamicFilterSplitsProcessed"`
					PhysicalWrittenDataSize        string  `json:"physicalWrittenDataSize"`
					BlockedWall                    string  `json:"blockedWall"`
					FinishCalls                    int     `json:"finishCalls"`
					FinishWall                     string  `json:"finishWall"`
					FinishCpu                      string  `json:"finishCpu"`
					UserMemoryReservation          string  `json:"userMemoryReservation"`
					RevocableMemoryReservation     string  `json:"revocableMemoryReservation"`
					SystemMemoryReservation        string  `json:"systemMemoryReservation"`
					PeakUserMemoryReservation      string  `json:"peakUserMemoryReservation"`
					PeakSystemMemoryReservation    string  `json:"peakSystemMemoryReservation"`
					PeakRevocableMemoryReservation string  `json:"peakRevocableMemoryReservation"`
					PeakTotalMemoryReservation     string  `json:"peakTotalMemoryReservation"`
					SpilledDataSize                string  `json:"spilledDataSize"`
					Info                           struct {
						Type                     string        `json:"@type"`
						BufferedBytes            int           `json:"bufferedBytes"`
						MaxBufferedBytes         int           `json:"maxBufferedBytes"`
						AverageBytesPerRequest   int           `json:"averageBytesPerRequest"`
						SuccessfulRequestsCount  int           `json:"successfulRequestsCount"`
						BufferedPages            int           `json:"bufferedPages"`
						NoMoreLocations          bool          `json:"noMoreLocations"`
						PageBufferClientStatuses []interface{} `json:"pageBufferClientStatuses"`
					} `json:"info,omitempty"`
				} `json:"operatorSummaries"`
			} `json:"stageStats"`
			Tasks []struct {
				TaskStatus struct {
					TaskId                     string        `json:"taskId"`
					TaskInstanceId             string        `json:"taskInstanceId"`
					Version                    int           `json:"version"`
					State                      string        `json:"state"`
					Self                       string        `json:"self"`
					NodeId                     string        `json:"nodeId"`
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
				} `json:"taskStatus"`
				LastHeartbeat time.Time `json:"lastHeartbeat"`
				OutputBuffers struct {
					Type               string        `json:"type"`
					State              string        `json:"state"`
					CanAddBuffers      bool          `json:"canAddBuffers"`
					CanAddPages        bool          `json:"canAddPages"`
					TotalBufferedBytes int           `json:"totalBufferedBytes"`
					TotalBufferedPages int           `json:"totalBufferedPages"`
					TotalRowsSent      int           `json:"totalRowsSent"`
					TotalPagesSent     int           `json:"totalPagesSent"`
					Buffers            []interface{} `json:"buffers"`
				} `json:"outputBuffers"`
				NoMoreSplits []string `json:"noMoreSplits"`
				Stats        struct {
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
					CumulativeUserMemory          float64       `json:"cumulativeUserMemory"`
					UserMemoryReservation         string        `json:"userMemoryReservation"`
					RevocableMemoryReservation    string        `json:"revocableMemoryReservation"`
					SystemMemoryReservation       string        `json:"systemMemoryReservation"`
					TotalScheduledTime            string        `json:"totalScheduledTime"`
					TotalCpuTime                  string        `json:"totalCpuTime"`
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
					Pipelines                     []struct {
						PipelineId                 int       `json:"pipelineId"`
						FirstStartTime             time.Time `json:"firstStartTime"`
						LastStartTime              time.Time `json:"lastStartTime"`
						LastEndTime                time.Time `json:"lastEndTime"`
						InputPipeline              bool      `json:"inputPipeline"`
						OutputPipeline             bool      `json:"outputPipeline"`
						TotalDrivers               int       `json:"totalDrivers"`
						QueuedDrivers              int       `json:"queuedDrivers"`
						QueuedPartitionedDrivers   int       `json:"queuedPartitionedDrivers"`
						RunningDrivers             int       `json:"runningDrivers"`
						RunningPartitionedDrivers  int       `json:"runningPartitionedDrivers"`
						BlockedDrivers             int       `json:"blockedDrivers"`
						CompletedDrivers           int       `json:"completedDrivers"`
						UserMemoryReservation      string    `json:"userMemoryReservation"`
						RevocableMemoryReservation string    `json:"revocableMemoryReservation"`
						SystemMemoryReservation    string    `json:"systemMemoryReservation"`
						QueuedTime                 struct {
							Count float64 `json:"count"`
							Total float64 `json:"total"`
							P01   float64 `json:"p01"`
							P05   float64 `json:"p05"`
							P10   float64 `json:"p10"`
							P25   float64 `json:"p25"`
							P50   float64 `json:"p50"`
							P75   float64 `json:"p75"`
							P90   float64 `json:"p90"`
							P95   float64 `json:"p95"`
							P99   float64 `json:"p99"`
							Min   float64 `json:"min"`
							Max   float64 `json:"max"`
							Avg   float64 `json:"avg"`
						} `json:"queuedTime"`
						ElapsedTime struct {
							Count float64 `json:"count"`
							Total float64 `json:"total"`
							P01   float64 `json:"p01"`
							P05   float64 `json:"p05"`
							P10   float64 `json:"p10"`
							P25   float64 `json:"p25"`
							P50   float64 `json:"p50"`
							P75   float64 `json:"p75"`
							P90   float64 `json:"p90"`
							P95   float64 `json:"p95"`
							P99   float64 `json:"p99"`
							Min   float64 `json:"min"`
							Max   float64 `json:"max"`
							Avg   float64 `json:"avg"`
						} `json:"elapsedTime"`
						TotalScheduledTime            string        `json:"totalScheduledTime"`
						TotalCpuTime                  string        `json:"totalCpuTime"`
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
						OperatorSummaries             []struct {
							StageId                        int     `json:"stageId"`
							PipelineId                     int     `json:"pipelineId"`
							OperatorId                     int     `json:"operatorId"`
							PlanNodeId                     string  `json:"planNodeId"`
							OperatorType                   string  `json:"operatorType"`
							TotalDrivers                   int     `json:"totalDrivers"`
							AddInputCalls                  int     `json:"addInputCalls"`
							AddInputWall                   string  `json:"addInputWall"`
							AddInputCpu                    string  `json:"addInputCpu"`
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
							GetOutputCpu                   string  `json:"getOutputCpu"`
							OutputDataSize                 string  `json:"outputDataSize"`
							OutputPositions                int     `json:"outputPositions"`
							DynamicFilterSplitsProcessed   int     `json:"dynamicFilterSplitsProcessed"`
							PhysicalWrittenDataSize        string  `json:"physicalWrittenDataSize"`
							BlockedWall                    string  `json:"blockedWall"`
							FinishCalls                    int     `json:"finishCalls"`
							FinishWall                     string  `json:"finishWall"`
							FinishCpu                      string  `json:"finishCpu"`
							UserMemoryReservation          string  `json:"userMemoryReservation"`
							RevocableMemoryReservation     string  `json:"revocableMemoryReservation"`
							SystemMemoryReservation        string  `json:"systemMemoryReservation"`
							PeakUserMemoryReservation      string  `json:"peakUserMemoryReservation"`
							PeakSystemMemoryReservation    string  `json:"peakSystemMemoryReservation"`
							PeakRevocableMemoryReservation string  `json:"peakRevocableMemoryReservation"`
							PeakTotalMemoryReservation     string  `json:"peakTotalMemoryReservation"`
							SpilledDataSize                string  `json:"spilledDataSize"`
							Info                           struct {
								Type                     string        `json:"@type"`
								BufferedBytes            int           `json:"bufferedBytes"`
								MaxBufferedBytes         int           `json:"maxBufferedBytes"`
								AverageBytesPerRequest   int           `json:"averageBytesPerRequest"`
								SuccessfulRequestsCount  int           `json:"successfulRequestsCount"`
								BufferedPages            int           `json:"bufferedPages"`
								NoMoreLocations          bool          `json:"noMoreLocations"`
								PageBufferClientStatuses []interface{} `json:"pageBufferClientStatuses"`
							} `json:"info,omitempty"`
						} `json:"operatorSummaries"`
						Drivers []interface{} `json:"drivers"`
					} `json:"pipelines"`
				} `json:"stats"`
				NeedsPlan bool `json:"needsPlan"`
			} `json:"tasks"`
			SubStages []struct {
				StageId string `json:"stageId"`
				State   string `json:"state"`
				Plan    struct {
					Id   string `json:"id"`
					Root struct {
						Type   string `json:"@type"`
						Id     string `json:"id"`
						Source struct {
							Type   string `json:"@type"`
							Id     string `json:"id"`
							Source struct {
								Type  string `json:"@type"`
								Id    string `json:"id"`
								Table struct {
									CatalogName     string `json:"catalogName"`
									ConnectorHandle struct {
										Type             string `json:"@type"`
										SchemaName       string `json:"schemaName"`
										TableName        string `json:"tableName"`
										PartitionColumns []struct {
											Type                string `json:"@type"`
											BaseColumnName      string `json:"baseColumnName"`
											BaseHiveColumnIndex int    `json:"baseHiveColumnIndex"`
											BaseHiveType        string `json:"baseHiveType"`
											BaseType            string `json:"baseType"`
											ColumnType          string `json:"columnType"`
										} `json:"partitionColumns"`
										DataColumns []struct {
											Type                string `json:"@type"`
											BaseColumnName      string `json:"baseColumnName"`
											BaseHiveColumnIndex int    `json:"baseHiveColumnIndex"`
											BaseHiveType        string `json:"baseHiveType"`
											BaseType            string `json:"baseType"`
											ColumnType          string `json:"columnType"`
										} `json:"dataColumns"`
										CompactEffectivePredicate struct {
											ColumnDomains []interface{} `json:"columnDomains"`
										} `json:"compactEffectivePredicate"`
										EnforcedConstraint struct {
											ColumnDomains []interface{} `json:"columnDomains"`
										} `json:"enforcedConstraint"`
										Transaction struct {
											Operation     string `json:"operation"`
											TransactionId int    `json:"transactionId"`
											WriteId       int    `json:"writeId"`
										} `json:"transaction"`
									} `json:"connectorHandle"`
									Transaction struct {
										Type string `json:"@type"`
										Uuid string `json:"uuid"`
									} `json:"transaction"`
								} `json:"table"`
								OutputSymbols []string `json:"outputSymbols"`
								Assignments   struct {
									Sentiment struct {
										Type                string `json:"@type"`
										BaseColumnName      string `json:"baseColumnName"`
										BaseHiveColumnIndex int    `json:"baseHiveColumnIndex"`
										BaseHiveType        string `json:"baseHiveType"`
										BaseType            string `json:"baseType"`
										ColumnType          string `json:"columnType"`
									} `json:"sentiment"`
									UserCountry struct {
										Type                string `json:"@type"`
										BaseColumnName      string `json:"baseColumnName"`
										BaseHiveColumnIndex int    `json:"baseHiveColumnIndex"`
										BaseHiveType        string `json:"baseHiveType"`
										BaseType            string `json:"baseType"`
										ColumnType          string `json:"columnType"`
									} `json:"user_country"`
									ContentsPositive struct {
										Type                string `json:"@type"`
										BaseColumnName      string `json:"baseColumnName"`
										BaseHiveColumnIndex int    `json:"baseHiveColumnIndex"`
										BaseHiveType        string `json:"baseHiveType"`
										BaseType            string `json:"baseType"`
										ColumnType          string `json:"columnType"`
									} `json:"contents_positive"`
								} `json:"assignments"`
								UpdateTarget                 bool `json:"updateTarget"`
								UseConnectorNodePartitioning bool `json:"useConnectorNodePartitioning"`
							} `json:"source"`
							Assignments struct {
								Assignments struct {
									ContentsPositive string `json:"contents_positive"`
									UserCountry      string `json:"user_country"`
									Sentiment0       string `json:"sentiment_0"`
									Hashvalue4       string `json:"$hashvalue_4"`
								} `json:"assignments"`
							} `json:"assignments"`
						} `json:"source"`
						Aggregations struct {
							Sum2 struct {
								ResolvedFunction struct {
									Signature struct {
										Name          string   `json:"name"`
										ReturnType    string   `json:"returnType"`
										ArgumentTypes []string `json:"argumentTypes"`
									} `json:"signature"`
									Id               string `json:"id"`
									TypeDependencies struct {
									} `json:"typeDependencies"`
									FunctionDependencies []interface{} `json:"functionDependencies"`
								} `json:"resolvedFunction"`
								Arguments []string `json:"arguments"`
								Distinct  bool     `json:"distinct"`
							} `json:"sum_2"`
							Avg1 struct {
								ResolvedFunction struct {
									Signature struct {
										Name          string   `json:"name"`
										ReturnType    string   `json:"returnType"`
										ArgumentTypes []string `json:"argumentTypes"`
									} `json:"signature"`
									Id               string `json:"id"`
									TypeDependencies struct {
									} `json:"typeDependencies"`
									FunctionDependencies []interface{} `json:"functionDependencies"`
								} `json:"resolvedFunction"`
								Arguments []string `json:"arguments"`
								Distinct  bool     `json:"distinct"`
							} `json:"avg_1"`
						} `json:"aggregations"`
						GroupingSets struct {
							GroupingKeys       []string      `json:"groupingKeys"`
							GroupingSetCount   int           `json:"groupingSetCount"`
							GlobalGroupingSets []interface{} `json:"globalGroupingSets"`
						} `json:"groupingSets"`
						PreGroupedSymbols []interface{} `json:"preGroupedSymbols"`
						Step              string        `json:"step"`
						HashSymbol        string        `json:"hashSymbol"`
					} `json:"root"`
					Symbols struct {
						Sentiment        string `json:"sentiment"`
						Hashvalue4       string `json:"$hashvalue_4"`
						Avg1             string `json:"avg_1"`
						UserCountry      string `json:"user_country"`
						ContentsPositive string `json:"contents_positive"`
						Sum2             string `json:"sum_2"`
						Sentiment0       string `json:"sentiment_0"`
					} `json:"symbols"`
					Partitioning struct {
						ConnectorHandle struct {
							Type         string `json:"@type"`
							Partitioning string `json:"partitioning"`
							Function     string `json:"function"`
						} `json:"connectorHandle"`
					} `json:"partitioning"`
					PartitionedSources []string `json:"partitionedSources"`
					PartitioningScheme struct {
						Partitioning struct {
							Handle struct {
								ConnectorHandle struct {
									Type         string `json:"@type"`
									Partitioning string `json:"partitioning"`
									Function     string `json:"function"`
								} `json:"connectorHandle"`
							} `json:"handle"`
							Arguments []struct {
								Expression string `json:"expression"`
							} `json:"arguments"`
						} `json:"partitioning"`
						OutputLayout         []string `json:"outputLayout"`
						HashColumn           string   `json:"hashColumn"`
						ReplicateNullsAndAny bool     `json:"replicateNullsAndAny"`
						BucketToPartition    []int    `json:"bucketToPartition"`
					} `json:"partitioningScheme"`
					StageExecutionDescriptor struct {
						Strategy                  string        `json:"strategy"`
						GroupedExecutionScanNodes []interface{} `json:"groupedExecutionScanNodes"`
					} `json:"stageExecutionDescriptor"`
					StatsAndCosts struct {
						Stats struct {
						} `json:"stats"`
						Costs struct {
						} `json:"costs"`
					} `json:"statsAndCosts"`
					JsonRepresentation string `json:"jsonRepresentation"`
				} `json:"plan"`
				Types      []string `json:"types"`
				StageStats struct {
					SchedulingComplete   time.Time `json:"schedulingComplete"`
					GetSplitDistribution struct {
						Count float64 `json:"count"`
						Total float64 `json:"total"`
						P01   float64 `json:"p01"`
						P05   float64 `json:"p05"`
						P10   float64 `json:"p10"`
						P25   float64 `json:"p25"`
						P50   float64 `json:"p50"`
						P75   float64 `json:"p75"`
						P90   float64 `json:"p90"`
						P95   float64 `json:"p95"`
						P99   float64 `json:"p99"`
						Min   float64 `json:"min"`
						Max   float64 `json:"max"`
						Avg   float64 `json:"avg"`
					} `json:"getSplitDistribution"`
					TotalTasks                     int           `json:"totalTasks"`
					RunningTasks                   int           `json:"runningTasks"`
					CompletedTasks                 int           `json:"completedTasks"`
					TotalDrivers                   int           `json:"totalDrivers"`
					QueuedDrivers                  int           `json:"queuedDrivers"`
					RunningDrivers                 int           `json:"runningDrivers"`
					BlockedDrivers                 int           `json:"blockedDrivers"`
					CompletedDrivers               int           `json:"completedDrivers"`
					CumulativeUserMemory           float64       `json:"cumulativeUserMemory"`
					UserMemoryReservation          string        `json:"userMemoryReservation"`
					RevocableMemoryReservation     string        `json:"revocableMemoryReservation"`
					TotalMemoryReservation         string        `json:"totalMemoryReservation"`
					PeakUserMemoryReservation      string        `json:"peakUserMemoryReservation"`
					PeakRevocableMemoryReservation string        `json:"peakRevocableMemoryReservation"`
					TotalScheduledTime             string        `json:"totalScheduledTime"`
					TotalCpuTime                   string        `json:"totalCpuTime"`
					TotalBlockedTime               string        `json:"totalBlockedTime"`
					FullyBlocked                   bool          `json:"fullyBlocked"`
					BlockedReasons                 []interface{} `json:"blockedReasons"`
					PhysicalInputDataSize          string        `json:"physicalInputDataSize"`
					PhysicalInputPositions         int           `json:"physicalInputPositions"`
					PhysicalInputReadTime          string        `json:"physicalInputReadTime"`
					InternalNetworkInputDataSize   string        `json:"internalNetworkInputDataSize"`
					InternalNetworkInputPositions  int           `json:"internalNetworkInputPositions"`
					RawInputDataSize               string        `json:"rawInputDataSize"`
					RawInputPositions              int           `json:"rawInputPositions"`
					ProcessedInputDataSize         string        `json:"processedInputDataSize"`
					ProcessedInputPositions        int           `json:"processedInputPositions"`
					BufferedDataSize               string        `json:"bufferedDataSize"`
					OutputDataSize                 string        `json:"outputDataSize"`
					OutputPositions                int           `json:"outputPositions"`
					PhysicalWrittenDataSize        string        `json:"physicalWrittenDataSize"`
					GcInfo                         struct {
						StageId          int `json:"stageId"`
						Tasks            int `json:"tasks"`
						FullGcTasks      int `json:"fullGcTasks"`
						MinFullGcSec     int `json:"minFullGcSec"`
						MaxFullGcSec     int `json:"maxFullGcSec"`
						TotalFullGcSec   int `json:"totalFullGcSec"`
						AverageFullGcSec int `json:"averageFullGcSec"`
					} `json:"gcInfo"`
					OperatorSummaries []struct {
						StageId                        int     `json:"stageId"`
						PipelineId                     int     `json:"pipelineId"`
						OperatorId                     int     `json:"operatorId"`
						PlanNodeId                     string  `json:"planNodeId"`
						OperatorType                   string  `json:"operatorType"`
						TotalDrivers                   int     `json:"totalDrivers"`
						AddInputCalls                  int     `json:"addInputCalls"`
						AddInputWall                   string  `json:"addInputWall"`
						AddInputCpu                    string  `json:"addInputCpu"`
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
						GetOutputCpu                   string  `json:"getOutputCpu"`
						OutputDataSize                 string  `json:"outputDataSize"`
						OutputPositions                int     `json:"outputPositions"`
						DynamicFilterSplitsProcessed   int     `json:"dynamicFilterSplitsProcessed"`
						PhysicalWrittenDataSize        string  `json:"physicalWrittenDataSize"`
						BlockedWall                    string  `json:"blockedWall"`
						FinishCalls                    int     `json:"finishCalls"`
						FinishWall                     string  `json:"finishWall"`
						FinishCpu                      string  `json:"finishCpu"`
						UserMemoryReservation          string  `json:"userMemoryReservation"`
						RevocableMemoryReservation     string  `json:"revocableMemoryReservation"`
						SystemMemoryReservation        string  `json:"systemMemoryReservation"`
						PeakUserMemoryReservation      string  `json:"peakUserMemoryReservation"`
						PeakSystemMemoryReservation    string  `json:"peakSystemMemoryReservation"`
						PeakRevocableMemoryReservation string  `json:"peakRevocableMemoryReservation"`
						PeakTotalMemoryReservation     string  `json:"peakTotalMemoryReservation"`
						SpilledDataSize                string  `json:"spilledDataSize"`
						Info                           struct {
							Type                        string `json:"@type"`
							RowsAdded                   int    `json:"rowsAdded"`
							PagesAdded                  int    `json:"pagesAdded"`
							OutputBufferPeakMemoryUsage int    `json:"outputBufferPeakMemoryUsage"`
						} `json:"info,omitempty"`
					} `json:"operatorSummaries"`
				} `json:"stageStats"`
				Tasks []struct {
					TaskStatus struct {
						TaskId                     string        `json:"taskId"`
						TaskInstanceId             string        `json:"taskInstanceId"`
						Version                    int           `json:"version"`
						State                      string        `json:"state"`
						Self                       string        `json:"self"`
						NodeId                     string        `json:"nodeId"`
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
					} `json:"taskStatus"`
					LastHeartbeat time.Time `json:"lastHeartbeat"`
					OutputBuffers struct {
						Type               string        `json:"type"`
						State              string        `json:"state"`
						CanAddBuffers      bool          `json:"canAddBuffers"`
						CanAddPages        bool          `json:"canAddPages"`
						TotalBufferedBytes int           `json:"totalBufferedBytes"`
						TotalBufferedPages int           `json:"totalBufferedPages"`
						TotalRowsSent      int           `json:"totalRowsSent"`
						TotalPagesSent     int           `json:"totalPagesSent"`
						Buffers            []interface{} `json:"buffers"`
					} `json:"outputBuffers"`
					NoMoreSplits []string `json:"noMoreSplits"`
					Stats        struct {
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
						CumulativeUserMemory          float64       `json:"cumulativeUserMemory"`
						UserMemoryReservation         string        `json:"userMemoryReservation"`
						RevocableMemoryReservation    string        `json:"revocableMemoryReservation"`
						SystemMemoryReservation       string        `json:"systemMemoryReservation"`
						TotalScheduledTime            string        `json:"totalScheduledTime"`
						TotalCpuTime                  string        `json:"totalCpuTime"`
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
						Pipelines                     []struct {
							PipelineId                 int       `json:"pipelineId"`
							FirstStartTime             time.Time `json:"firstStartTime"`
							LastStartTime              time.Time `json:"lastStartTime"`
							LastEndTime                time.Time `json:"lastEndTime"`
							InputPipeline              bool      `json:"inputPipeline"`
							OutputPipeline             bool      `json:"outputPipeline"`
							TotalDrivers               int       `json:"totalDrivers"`
							QueuedDrivers              int       `json:"queuedDrivers"`
							QueuedPartitionedDrivers   int       `json:"queuedPartitionedDrivers"`
							RunningDrivers             int       `json:"runningDrivers"`
							RunningPartitionedDrivers  int       `json:"runningPartitionedDrivers"`
							BlockedDrivers             int       `json:"blockedDrivers"`
							CompletedDrivers           int       `json:"completedDrivers"`
							UserMemoryReservation      string    `json:"userMemoryReservation"`
							RevocableMemoryReservation string    `json:"revocableMemoryReservation"`
							SystemMemoryReservation    string    `json:"systemMemoryReservation"`
							QueuedTime                 struct {
								Count float64 `json:"count"`
								Total float64 `json:"total"`
								P01   float64 `json:"p01"`
								P05   float64 `json:"p05"`
								P10   float64 `json:"p10"`
								P25   float64 `json:"p25"`
								P50   float64 `json:"p50"`
								P75   float64 `json:"p75"`
								P90   float64 `json:"p90"`
								P95   float64 `json:"p95"`
								P99   float64 `json:"p99"`
								Min   float64 `json:"min"`
								Max   float64 `json:"max"`
								Avg   float64 `json:"avg"`
							} `json:"queuedTime"`
							ElapsedTime struct {
								Count float64 `json:"count"`
								Total float64 `json:"total"`
								P01   float64 `json:"p01"`
								P05   float64 `json:"p05"`
								P10   float64 `json:"p10"`
								P25   float64 `json:"p25"`
								P50   float64 `json:"p50"`
								P75   float64 `json:"p75"`
								P90   float64 `json:"p90"`
								P95   float64 `json:"p95"`
								P99   float64 `json:"p99"`
								Min   float64 `json:"min"`
								Max   float64 `json:"max"`
								Avg   float64 `json:"avg"`
							} `json:"elapsedTime"`
							TotalScheduledTime            string        `json:"totalScheduledTime"`
							TotalCpuTime                  string        `json:"totalCpuTime"`
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
							OperatorSummaries             []struct {
								StageId                        int     `json:"stageId"`
								PipelineId                     int     `json:"pipelineId"`
								OperatorId                     int     `json:"operatorId"`
								PlanNodeId                     string  `json:"planNodeId"`
								OperatorType                   string  `json:"operatorType"`
								TotalDrivers                   int     `json:"totalDrivers"`
								AddInputCalls                  int     `json:"addInputCalls"`
								AddInputWall                   string  `json:"addInputWall"`
								AddInputCpu                    string  `json:"addInputCpu"`
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
								GetOutputCpu                   string  `json:"getOutputCpu"`
								OutputDataSize                 string  `json:"outputDataSize"`
								OutputPositions                int     `json:"outputPositions"`
								DynamicFilterSplitsProcessed   int     `json:"dynamicFilterSplitsProcessed"`
								PhysicalWrittenDataSize        string  `json:"physicalWrittenDataSize"`
								BlockedWall                    string  `json:"blockedWall"`
								FinishCalls                    int     `json:"finishCalls"`
								FinishWall                     string  `json:"finishWall"`
								FinishCpu                      string  `json:"finishCpu"`
								UserMemoryReservation          string  `json:"userMemoryReservation"`
								RevocableMemoryReservation     string  `json:"revocableMemoryReservation"`
								SystemMemoryReservation        string  `json:"systemMemoryReservation"`
								PeakUserMemoryReservation      string  `json:"peakUserMemoryReservation"`
								PeakSystemMemoryReservation    string  `json:"peakSystemMemoryReservation"`
								PeakRevocableMemoryReservation string  `json:"peakRevocableMemoryReservation"`
								PeakTotalMemoryReservation     string  `json:"peakTotalMemoryReservation"`
								SpilledDataSize                string  `json:"spilledDataSize"`
								Info                           struct {
									Type                        string `json:"@type"`
									RowsAdded                   int    `json:"rowsAdded"`
									PagesAdded                  int    `json:"pagesAdded"`
									OutputBufferPeakMemoryUsage int    `json:"outputBufferPeakMemoryUsage"`
								} `json:"info,omitempty"`
							} `json:"operatorSummaries"`
							Drivers []interface{} `json:"drivers"`
						} `json:"pipelines"`
					} `json:"stats"`
					NeedsPlan bool `json:"needsPlan"`
				} `json:"tasks"`
				SubStages []interface{} `json:"subStages"`
				Tables    struct {
					Field1 struct {
						TableName string `json:"tableName"`
						Predicate struct {
							ColumnDomains []struct {
								Column struct {
									Type                string `json:"@type"`
									BaseColumnName      string `json:"baseColumnName"`
									BaseHiveColumnIndex int    `json:"baseHiveColumnIndex"`
									BaseHiveType        string `json:"baseHiveType"`
									BaseType            string `json:"baseType"`
									ColumnType          string `json:"columnType"`
								} `json:"column"`
								Domain struct {
									Values struct {
										Type         string `json:"@type"`
										Type1        string `json:"type"`
										Inclusive    []bool `json:"inclusive"`
										SortedRanges string `json:"sortedRanges"`
									} `json:"values"`
									NullAllowed bool `json:"nullAllowed"`
								} `json:"domain"`
							} `json:"columnDomains"`
						} `json:"predicate"`
					} `json:"0"`
				} `json:"tables"`
			} `json:"subStages"`
			Tables struct {
			} `json:"tables"`
		} `json:"subStages"`
		Tables struct {
		} `json:"tables"`
	} `json:"outputStage"`
	Warnings []interface{} `json:"warnings"`
	Inputs   []struct {
		CatalogName string `json:"catalogName"`
		Schema      string `json:"schema"`
		Table       string `json:"table"`
		Columns     []struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"columns"`
		FragmentId string `json:"fragmentId"`
		PlanNodeId string `json:"planNodeId"`
	} `json:"inputs"`
	ReferencedTables []struct {
		Catalog       string        `json:"catalog"`
		Schema        string        `json:"schema"`
		Table         string        `json:"table"`
		Authorization string        `json:"authorization"`
		Filters       []interface{} `json:"filters"`
		Columns       []struct {
			Column string        `json:"column"`
			Masks  []interface{} `json:"masks"`
		} `json:"columns"`
		DirectlyReferenced bool `json:"directlyReferenced"`
	} `json:"referencedTables"`
	Routines []struct {
		Routine       string `json:"routine"`
		Authorization string `json:"authorization"`
	} `json:"routines"`
	ResourceGroupId []string `json:"resourceGroupId"`
	QueryType       string   `json:"queryType"`
	FinalQueryInfo  bool     `json:"finalQueryInfo"`
}
