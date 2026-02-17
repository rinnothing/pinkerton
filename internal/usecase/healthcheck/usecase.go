package healthcheck

import (
	"context"
	"time"

	"github.com/rinnothing/pinkerton/internal/model"
)

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

type Storage interface {
	LoadStatus(url string, status int)
	LoadParams(url string, period time.Duration)
	GetModel(url string) *model.Target
	RemoveModel(url string)
}

type usecaseImplementation struct {
	em *emitter
	st Storage
}

func New(ctx context.Context, initModels []model.Target, threads int, tester Pinger, storage Storage) *usecaseImplementation {
	impl := usecaseImplementation{
		em: newEmitter(),
		st: storage,
	}

	for _, model := range initModels {
		impl.em.AddEvent(constructTimeAndPeriod(model.Period), model.URL)
	}

	ch := impl.em.Start(ctx)
	receiveEvents(threads, ch, tester, storage)
	return &impl
}

func (u *usecaseImplementation) AddTarget(newTarget *model.Target) error {
	u.st.LoadParams(newTarget.URL, newTarget.Period)
	u.em.AddEvent(constructTimeAndPeriod(newTarget.Period), newTarget.URL)
	return nil
}

// UpdateTarget implements Usecase.
func (u *usecaseImplementation) UpdateTarget(target *model.Target) error {
	u.st.LoadParams(target.URL, target.Period)
	u.em.UpdateEvent(constructTimeAndPeriod(target.Period), target.URL)
	return nil
}

func (u *usecaseImplementation) GetTarget(url string) (*model.Target, error) {
	mdl := u.st.GetModel(url)
	return mdl, nil
}

// RemoveTarget implements Usecase.
func (u *usecaseImplementation) RemoveTarget(url string) error {
	u.em.RemoveEvent(url)
	u.st.RemoveModel(url)
	return nil
}
