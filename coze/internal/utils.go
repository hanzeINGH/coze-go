package internal

import "net/http"

const (
	logIDHeader = "x-tt-logid"
)

func getLogID(header http.Header) string {
	return header.Get(logIDHeader)
}
