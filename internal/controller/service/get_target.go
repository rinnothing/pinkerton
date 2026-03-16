package service

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/rinnothing/pinkerton/internal/model"
)

func (ctr *controller) GetTarget(resp http.ResponseWriter, req *http.Request) {
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

	url := req.PathValue("url")
	if url == "" {
		badRequest(resp, "url should be encoded in path")
		return
	}

	err = checkUrl(resp, url)
	if err != nil {
		return
	}

	mdl, err := ctr.uc.GetTarget(url)
	if errors.Is(err, model.ErrUrlNotExists) {
		notFound(resp)
		return
	} else if err != nil {
		slog.Error("error getting target by url", "url", url, "error", err)
		internalError(resp)
		return
	}

	err = json.NewEncoder(resp).Encode(mdl)
	if err != nil {
		slog.Error("error encoding value to json response", "value", mdl, "error", err)
	}
}
