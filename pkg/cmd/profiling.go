package cmd

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

func enableProfilingServer(addr string) {
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

