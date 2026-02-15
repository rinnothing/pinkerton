package healthcheck

import "github.com/rinnothing/pinkerton/internal/model"

type Usecase interface {
	AddTarget(newTarget *model.Target) error

	UpdateTarget(target *model.Target) error

	RemoveTarget(url string) error

	GetTarget(url string) (*model.Target, error)
}

var _ Usecase = &usecaseImplementation{}

type Pinger interface {
	// Ping returns status for the said url
	Ping(url string) (int, error)
}

type usecaseImplementation struct {
}

// AddTarget implements Usecase.
func (u *usecaseImplementation) AddTarget(newTarget *model.Target) error {
	panic("unimplemented")
}

// GetTarget implements Usecase.
func (u *usecaseImplementation) GetTarget(url string) (*model.Target, error) {
	panic("unimplemented")
}

// RemoveTarget implements Usecase.
func (u *usecaseImplementation) RemoveTarget(url string) error {
	panic("unimplemented")
}

// UpdateTarget implements Usecase.
func (u *usecaseImplementation) UpdateTarget(target *model.Target) error {
	panic("unimplemented")
}
