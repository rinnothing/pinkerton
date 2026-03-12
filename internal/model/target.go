package model

import (
	"time"
)

type Target struct {
	URL          string        `json:"url"`
	LastStatus   int           `json:"last_status"`
	LastResponse time.Time     `json:"last_response"`
	Period       time.Duration `json:"period"`
}

type TargetRequest struct {
	URL    string        `json:"url"`
	Period time.Duration `json:"period"`
}
