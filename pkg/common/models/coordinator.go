package models

import (
	"net/url"
)

type Coordinator struct {
	Name    string
	URL     *url.URL
	Tags    map[string]string
	Enabled bool
}
