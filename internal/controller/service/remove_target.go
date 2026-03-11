package service

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/rinnothing/pinkerton/internal/model"
)

func (ctr *controller) RemoveTarget(resp http.ResponseWriter, req *http.Request) {
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

	err = ctr.uc.RemoveTarget(url)
	if errors.Is(err, model.ErrUrlNotExists) {
		notFound(resp)
		return
	} else if err != nil {
		slog.Error("error removing target by url", "url", url, "error", err)
		internalError(resp)
		return
	}

	statusOK(resp)
}
