package service

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/rinnothing/pinkerton/internal/model"
	"github.com/rinnothing/pinkerton/internal/usecase/healthcheck"
	"github.com/rinnothing/pinkerton/pkg/checks"
)

type controller struct {
	*http.ServeMux
	uc healthcheck.Usecase
}

func New(uc healthcheck.Usecase) *controller {
	ctr := &controller{
		http.NewServeMux(),
		uc,
	}

	ctr.registerEndpoints()
	return ctr
}

func (ctr *controller) registerEndpoints() {
	ctr.HandleFunc("GET /targets/{url}", ctr.GetTarget)
	ctr.HandleFunc("GET /targets", ctr.GetAllTargets)
	ctr.HandleFunc("POST /targets", ctr.AddTarget)
	ctr.HandleFunc("PUT /targets", ctr.UpdateTarget)
	ctr.HandleFunc("DELETE /targets/{url}", ctr.RemoveTarget)

	ctr.HandleFunc("GET /health", ctr.Health)
}

func badRequest(resp http.ResponseWriter, comment string) {
	resp.WriteHeader(http.StatusBadRequest)
	_, err := strings.NewReader(comment).WriteTo(resp)

	if err != nil {
		slog.Error("error writing bad request response", "error", err)
	}
}

func internalError(resp http.ResponseWriter) {
	resp.WriteHeader(http.StatusInternalServerError)
}

func notFound(resp http.ResponseWriter) {
	resp.WriteHeader(http.StatusNotFound)
}

func statusConflict(resp http.ResponseWriter, comment string) {
	resp.WriteHeader(http.StatusConflict)
	_, err := strings.NewReader(comment).WriteTo(resp)

	if err != nil {
		slog.Error("error writing conflict response", "error", err)
	}
}

func statusOK(resp http.ResponseWriter) {
	resp.WriteHeader(http.StatusOK)
}

func checkUrl(resp http.ResponseWriter, url string) error {
	if err := checks.CheckUrl(url); err != nil {
		badRequest(resp, err.Error())
		return err
	}
	return nil
}

func checkPeriod(resp http.ResponseWriter, p time.Duration) error {
	if err := checks.CheckPeriod(p); err != nil {
		badRequest(resp, err.Error())
		return err
	}
	return nil
}

func checkRequest(resp http.ResponseWriter, t *model.TargetRequest) error {
	if err := checkUrl(resp, t.URL); err != nil {
		return err
	}
	if err := checkPeriod(resp, t.Period); err != nil {
		return err
	}
	return nil
}
