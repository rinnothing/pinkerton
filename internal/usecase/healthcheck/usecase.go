package healthcheck

import (
	"context"
	"sync"
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
	StoreStatus(url string, status int)
	AddParams(url string, period time.Duration) error    // returns model.ErrUrlExists if url already exists
	UpdateParams(url string, period time.Duration) error // returns model.ErrUrlNotExists if url doesn't exist
	GetModel(url string) (*model.Target, error)          // return model.ErrUrlNotExists if url doesn't exist
	RemoveModel(url string) error                        // returns model.ErrUrlNotExists if url doesn't exist
}

type usecaseImplementation struct {
	mx *sync.RWMutex
	em *emitter
	st Storage
}

func New(ctx context.Context, initModels []model.Target, threads int, tester Pinger, storage Storage) *usecaseImplementation {
	impl := usecaseImplementation{
		mx: new(sync.RWMutex),
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
	u.mx.Lock()
	defer u.mx.Unlock()

	err := u.st.AddParams(newTarget.URL, newTarget.Period)
	if err != nil {
		return err
	}

	u.em.AddEvent(constructTimeAndPeriod(newTarget.Period), newTarget.URL)
	return nil
}

func (u *usecaseImplementation) UpdateTarget(target *model.Target) error {
	u.mx.Lock()
	defer u.mx.Unlock()

	err := u.st.UpdateParams(target.URL, target.Period)
	if err != nil {
		return err
	}

	u.em.UpdateEvent(constructTimeAndPeriod(target.Period), target.URL)
	return nil
}

func (u *usecaseImplementation) GetTarget(url string) (*model.Target, error) {
	u.mx.RLock()
	defer u.mx.RUnlock()

	mdl, err := u.st.GetModel(url)
	if err != nil {
		return nil, err
	}

	return mdl, nil
}

// RemoveTarget implements Usecase.
func (u *usecaseImplementation) RemoveTarget(url string) error {
	u.mx.Lock()
	defer u.mx.Unlock()

	err := u.st.RemoveModel(url)
	if err != nil {
		return err
	}

	u.em.RemoveEvent(url)
	return nil
}
