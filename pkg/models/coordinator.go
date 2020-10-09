package models

import "net/url"

type PrestoDist string

var (
	PrestoDistSql PrestoDist = "prestosql"
	PrestoDistDb  PrestoDist = "prestodb"
)

type Coordinator struct {
	Name         string
	URL          *url.URL
	Tags         map[string]string
	Enabled      bool
	Distribution PrestoDist
}
