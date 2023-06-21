package trino

const (
	QueryFinished string = "FINISHED"
)

type QueryList []QueryListItem

type QueryListItem struct {
	QueryId          string     `json:"queryId"`
	SessionUser      string     `json:"sessionUser"`
	SessionPrincipal string     `json:"sessionPrincipal"`
	ResourceGroupId  []string   `json:"resourceGroupId"`
	State            string     `json:"state"`
	MemoryPool       string     `json:"memoryPool"`
	Scheduled        bool       `json:"scheduled"`
	Self             string     `json:"self"`
	Query            string     `json:"query"`
	QueryStats       QueryStats `json:"queryStats"`
	QueryType        string     `json:"queryType"`
	ErrorType        string     `json:"errorType,omitempty"`
	ErrorCode        struct {
		Code int    `json:"code"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"errorCode,omitempty"`
}
