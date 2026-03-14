package service

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

func (ctr *controller) GetAllTargets(resp http.ResponseWriter, req *http.Request) {
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

	mdls := ctr.uc.GetAllTargets()

	err = json.NewEncoder(resp).Encode(mdls)
	if err != nil {
		slog.Error("error encoding value to json response", "value", mdls, "error", err)
	}
}
