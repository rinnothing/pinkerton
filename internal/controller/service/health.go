package service

import (
	"io"
	"log/slog"
	"net/http"
)

func (ctr *controller) Health(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	val, err := io.Copy(io.Discard, req.Body)
	if err != nil {
		slog.Error("error reading request body", "error", err)
		internalError(resp)
		return
	}
	if val != 0 {
		badRequest(resp, "request body should be empty")
		return
	}

	resp.WriteHeader(http.StatusOK)
}
