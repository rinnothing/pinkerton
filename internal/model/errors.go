package model

import "errors"

var (
	ErrUrlExists    = errors.New("url already exists")
	ErrUrlNotExists = errors.New("url doesn't exist")
)
