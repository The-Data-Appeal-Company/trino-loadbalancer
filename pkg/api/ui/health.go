package ui

import "net/http"

func healthProbe(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
}
