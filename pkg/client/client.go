package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/rinnothing/pinkerton/internal/model"
)

type Target model.Target
type Targets []*Target
type TargetRequest model.TargetRequest

var ErrUrlExists = model.ErrUrlExists
var ErrUrlNotExists = model.ErrUrlNotExists
var ErrInternal = errors.New("internal server error")
var ErrUnexpected = errors.New("unexpected return code")

func errUnexpected(code int) error {
	return fmt.Errorf("%w: %d", ErrUnexpected, code)
}

func New(client http.Client, hostname, port string, timeout time.Duration) *Client {
	return &Client{
		client: client,
		path:   net.JoinHostPort(hostname, port),
	}
}

type Client struct {
	client  http.Client
	path    string
	timeout time.Duration
}

func decodeBody(body io.Reader, target any, queryUrl string) error {
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&target)
	if err != nil {
		return fmt.Errorf("error decoding response from url %q: %w", queryUrl, err)
	}

	return nil
}

func (c *Client) formAndDoRequest(ctx context.Context, method string, queryUrl string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, queryUrl, body)
	if err != nil {
		return nil, fmt.Errorf("can't form request for url %q: %w", queryUrl, err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't send request for url %q: %w", queryUrl, err)
	}

	return resp, nil
}

func (c *Client) constructURL(parts ...string) (string, error) {
	queryUrl, err := url.JoinPath(c.path, parts...)
	if err != nil {
		return "", fmt.Errorf("can't make query: %w", err)
	}

	return queryUrl, nil
}

func marshalBody(val any) (io.Reader, error) {
	valBytes, err := json.Marshal(val)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	return bytes.NewReader(valBytes), nil
}

func (c *Client) GetTarget(ctx context.Context, targetUrl string) (*Target, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	queryUrl, err := c.constructURL("targets", url.PathEscape(targetUrl))
	if err != nil {
		return nil, err
	}

	resp, err := c.formAndDoRequest(ctx, http.MethodGet, queryUrl, bytes.NewReader([]byte{}))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusNotFound:
		return nil, ErrUrlNotExists
	case http.StatusInternalServerError:
		return nil, ErrInternal
	default:
		return nil, errUnexpected(resp.StatusCode)
	}

	var tgt Target
	if err = decodeBody(resp.Body, &tgt, queryUrl); err != nil {
		return nil, err
	}

	return &tgt, nil
}

func (c *Client) GetAllTargets(ctx context.Context) (Targets, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	queryUrl, err := c.constructURL("targets")
	if err != nil {
		return nil, err
	}

	resp, err := c.formAndDoRequest(ctx, http.MethodGet, queryUrl, bytes.NewReader([]byte{}))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusInternalServerError:
		return nil, ErrInternal
	default:
		return nil, errUnexpected(resp.StatusCode)
	}

	var tgts Targets
	if err = decodeBody(resp.Body, &tgts, queryUrl); err != nil {
		return nil, err
	}

	return tgts, err
}

func (c *Client) AddTarget(ctx context.Context, req *TargetRequest) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	queryUrl, err := c.constructURL("targets")
	if err != nil {
		return err
	}

	bodyReader, err := marshalBody(req)
	if err != nil {
		return err
	}

	resp, err := c.formAndDoRequest(ctx, http.MethodPost, queryUrl, bodyReader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusConflict:
		return ErrUrlExists
	case http.StatusInternalServerError:
		return ErrInternal
	default:
		return errUnexpected(resp.StatusCode)
	}

	return nil
}

func (c *Client) UpdateTarget(ctx context.Context, req *TargetRequest) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	queryUrl, err := c.constructURL("targets")
	if err != nil {
		return err
	}

	bodyReader, err := marshalBody(req)
	if err != nil {
		return err
	}

	resp, err := c.formAndDoRequest(ctx, http.MethodPut, queryUrl, bodyReader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusNotFound:
		return ErrUrlNotExists
	case http.StatusInternalServerError:
		return ErrInternal
	default:
		return errUnexpected(resp.StatusCode)
	}

	return nil
}

func (c *Client) RemoveTarget(ctx context.Context, targetUrl string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	queryUrl, err := c.constructURL("targets", url.PathEscape(targetUrl))
	if err != nil {
		return err
	}

	resp, err := c.formAndDoRequest(ctx, http.MethodDelete, queryUrl, bytes.NewReader([]byte{}))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusNotFound:
		return ErrUrlNotExists
	case http.StatusInternalServerError:
		return ErrInternal
	default:
		return errUnexpected(resp.StatusCode)
	}

	return nil
}

func (c *Client) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	queryUrl, err := c.constructURL("health")
	if err != nil {
		return err
	}

	resp, err := c.formAndDoRequest(ctx, http.MethodGet, queryUrl, bytes.NewReader([]byte{}))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusInternalServerError:
		return ErrInternal
	default:
		return errUnexpected(resp.StatusCode)
	}

	return nil
}
