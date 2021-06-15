package trino

import "time"

const (
	QueryFinished string = "FINISHED"
)

type QueryList []QueryListItem

type QueryListItem struct {
	QueryId string `json:"queryId"`
	Session struct {
		QueryId                  string        `json:"queryId"`
		TransactionId            string        `json:"transactionId"`
		ClientTransactionSupport bool          `json:"clientTransactionSupport"`
		User                     string        `json:"user"`
		Groups                   []interface{} `json:"groups"`
		Principal                string        `json:"principal"`
		Source                   string        `json:"source,omitempty"`
		Catalog                  string        `json:"catalog,omitempty"`
		Schema                   string        `json:"schema,omitempty"`
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
	ResourceGroupId []string `json:"resourceGroupId"`
	State           string   `json:"state"`
	MemoryPool      string   `json:"memoryPool"`
	Scheduled       bool     `json:"scheduled"`
	Self            string   `json:"self"`
	Query           string   `json:"query"`
	QueryStats      struct {
		CreateTime                 time.Time     `json:"createTime"`
		EndTime                    time.Time     `json:"endTime"`
		QueuedTime                 string        `json:"queuedTime"`
		ElapsedTime                string        `json:"elapsedTime"`
		ExecutionTime              string        `json:"executionTime"`
		TotalDrivers               int           `json:"totalDrivers"`
		QueuedDrivers              int           `json:"queuedDrivers"`
		RunningDrivers             int           `json:"runningDrivers"`
		CompletedDrivers           int           `json:"completedDrivers"`
		RawInputDataSize           string        `json:"rawInputDataSize"`
		RawInputPositions          int           `json:"rawInputPositions"`
		PhysicalInputDataSize      string        `json:"physicalInputDataSize"`
		CumulativeUserMemory       interface{}   `json:"cumulativeUserMemory"`
		UserMemoryReservation      string        `json:"userMemoryReservation"`
		TotalMemoryReservation     string        `json:"totalMemoryReservation"`
		PeakUserMemoryReservation  string        `json:"peakUserMemoryReservation"`
		PeakTotalMemoryReservation string        `json:"peakTotalMemoryReservation"`
		TotalCpuTime               string        `json:"totalCpuTime"`
		TotalScheduledTime         string        `json:"totalScheduledTime"`
		FullyBlocked               bool          `json:"fullyBlocked"`
		BlockedReasons             []interface{} `json:"blockedReasons"`
		ProgressPercentage         float64       `json:"progressPercentage"`
	} `json:"queryStats"`
	QueryType string `json:"queryType"`
	ErrorType string `json:"errorType,omitempty"`
	ErrorCode struct {
		Code int    `json:"code"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"errorCode,omitempty"`
}
