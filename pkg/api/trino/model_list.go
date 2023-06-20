package trino

import "time"

const (
	QueryFinished string = "FINISHED"
)

type QueryList []QueryListItem

type QueryListItem struct {
	QueryId         string           `json:"queryId"`
	Session         QueryItemSession `json:"session"`
	ResourceGroupId []string         `json:"resourceGroupId"`
	State           string           `json:"state"`
	MemoryPool      string           `json:"memoryPool"`
	Scheduled       bool             `json:"scheduled"`
	Self            string           `json:"self"`
	Query           string           `json:"query"`
	QueryStats      QueryStats       `json:"queryStats"`
	QueryType       string           `json:"queryType"`
	ErrorType       string           `json:"errorType,omitempty"`
	ErrorCode       struct {
		Code int    `json:"code"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"errorCode,omitempty"`
}

type QueryItemSession struct {
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
}
