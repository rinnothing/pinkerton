package model

import (
	"time"
)

type Target struct {
	URL          string
	LastStatus   int
	LastResponse time.Time
	Period       time.Duration
}
