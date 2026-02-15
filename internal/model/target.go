package model

import (
	"time"
)

type Target struct {
	URL        string
	LastStatus int
	Timeout    time.Duration
	Period     time.Duration
}
