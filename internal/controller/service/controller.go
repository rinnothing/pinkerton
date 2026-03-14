package service

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/rinnothing/pinkerton/internal/usecase/healthcheck"
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
	ctr.HandleFunc("GET /targets/", ctr.GetAllTargets)
	ctr.HandleFunc("POST /targets/", ctr.AddTarget)
	ctr.HandleFunc("PUT /targets/", ctr.UpdateTarget)
	ctr.HandleFunc("DELETE /targets/{url}", ctr.RemoveTarget)

	ctr.HandleFunc("GET /health/", ctr.Health)
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
	resp.WriteHeader(http.StatusBadRequest)
	_, err := strings.NewReader(comment).WriteTo(resp)

	if err != nil {
		slog.Error("error writing conflict response", "error", err)
	}
}

func statusOK(resp http.ResponseWriter) {
	resp.WriteHeader(http.StatusOK)
}
