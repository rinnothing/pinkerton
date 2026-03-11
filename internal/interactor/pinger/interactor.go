package pinger

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/rinnothing/pinkerton/internal/usecase/healthcheck"
)

var _ healthcheck.Pinger = &interactor{}

type interactor struct {
	ctx     context.Context
	client  *http.Client
	timeout time.Duration
}

func (i *interactor) Ping(url string) (int, error) {
	ctx, cancel := context.WithTimeout(i.ctx, i.timeout)
	defer cancel()

	emptyReader := bytes.NewReader([]byte{})
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, emptyReader)
	if err != nil {
		return 0, err
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}

func New(ctx context.Context, timeout time.Duration) *interactor {
	return &interactor{
		ctx:     ctx,
		client:  http.DefaultClient,
		timeout: timeout,
	}
}
