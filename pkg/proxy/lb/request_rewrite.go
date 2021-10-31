package lb

import "net/http"

type RequestRewriter interface {
	Rewrite(*http.Request) (*http.Request, error)
}

type MultiRequestRewriter struct {
	rewriters []RequestRewriter
}

func RequestReWriters(rrws ...RequestRewriter) RequestRewriter {
	return MultiRequestRewriter{rewriters: rrws}
}

func (m MultiRequestRewriter) Rewrite(request *http.Request) (*http.Request, error) {
	for _, rewriter := range m.rewriters {
		rewritten, err := rewriter.Rewrite(request)
		if err != nil {
			return nil, err
		}
		request = rewritten
	}
	return request, nil
}
