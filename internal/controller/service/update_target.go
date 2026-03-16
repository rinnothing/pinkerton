package service

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/rinnothing/pinkerton/internal/model"
)

func (ctr *controller) UpdateTarget(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()

	var tgtReq model.TargetRequest
	err := dec.Decode(&tgtReq)
	if err != nil {
		badRequest(resp, "can't decode request")
		return
	}

	err = checkRequest(resp, &tgtReq)
	if err != nil {
		return
	}

	err = ctr.uc.UpdateTarget(&model.Target{URL: tgtReq.URL, Period: tgtReq.Period})
	if errors.Is(err, model.ErrUrlNotExists) {
		notFound(resp)
		return
	} else if err != nil {
		slog.Error("error updating target", "target", tgtReq, "error", err)
		internalError(resp)
		return
	}

	statusOK(resp)
}
