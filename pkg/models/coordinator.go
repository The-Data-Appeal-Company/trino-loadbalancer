package models

import (
	"fmt"
	"net/url"
)

type PrestoDist string

const prestoSql = "prestosql"
const prestoDb = "prestodb"

func ParsePrestoDist(distRaw string) (PrestoDist, error) {
	switch distRaw {

	case prestoSql:
		return PrestoDistSql, nil
	case prestoDb:
		return PrestoDistDb, nil
	}

	return "", fmt.Errorf("cannot parse presto distribution type %s", distRaw)

}

var (
	PrestoDistSql PrestoDist = prestoSql
	PrestoDistDb  PrestoDist = prestoDb
)

type Coordinator struct {
	Name         string
	URL          *url.URL
	Tags         map[string]string
	Enabled      bool
	Distribution PrestoDist
}
