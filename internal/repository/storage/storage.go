package storage

import (
	"log/slog"
	"sync"
	"time"

	"github.com/rinnothing/pinkerton/internal/model"
	"github.com/rinnothing/pinkerton/internal/usecase/healthcheck"
)

var _ healthcheck.Storage = (*storage)(nil)

type storage struct {
	paramsMx  *sync.RWMutex
	paramsMap map[string]time.Duration

	statusMap *sync.Map
}

func New() *storage {
	return &storage{
		paramsMx:  new(sync.RWMutex),
		paramsMap: make(map[string]time.Duration),
		statusMap: new(sync.Map),
	}
}

type statusAndTime struct {
	status int
	moment time.Time
}

func (s *storage) StoreStatus(url string, status int) {
	st := statusAndTime{
		status: status,
		// it won't be the latest time actually, because some stuck thread could put another thread with status before, but since we measure often it shouldn't be a big problem
		moment: time.Now(),
	}

	slog.Debug("storing new status and time on url", "url", url, "status", st.status, "time", st.moment)
	s.statusMap.Store(url, st)
}

func (s *storage) AddParams(url string, period time.Duration) error {
	s.paramsMx.Lock()
	defer s.paramsMx.Unlock()

	_, ok := s.paramsMap[url]
	if ok {
		return model.ErrUrlExists
	}

	slog.Debug("add params on url", "url", url, "period", period)
	s.paramsMap[url] = period
	return nil
}

func (s *storage) UpdateParams(url string, period time.Duration) error {
	s.paramsMx.Lock()
	defer s.paramsMx.Unlock()

	_, ok := s.paramsMap[url]
	if !ok {
		return model.ErrUrlNotExists
	}

	slog.Debug("update params on url", "url", url, "period", period)
	s.paramsMap[url] = period
	return nil
}

func (s *storage) GetTarget(url string) (*model.Target, error) {
	s.paramsMx.RLock()
	defer s.paramsMx.RUnlock()

	period, ok := s.paramsMap[url]
	if !ok {
		return nil, model.ErrUrlNotExists
	}

	var status statusAndTime
	val, ok := s.statusMap.Load(url)
	if ok {
		status = val.(statusAndTime)
	}

	slog.Debug("return info on url", "url", url, "period", period, "last status", status.status, "last response", status.moment)
	return &model.Target{
		URL:          url,
		LastStatus:   status.status,
		LastResponse: status.moment,
		Period:       period,
	}, nil
}

func (s *storage) GetAllTargets() []*model.Target {
	s.paramsMx.RLock()
	defer s.paramsMx.RUnlock()

	var targets []*model.Target
	for url, period := range s.paramsMap {
		var status statusAndTime
		val, ok := s.statusMap.Load(url)
		if ok {
			status = val.(statusAndTime)
		}

		targets = append(targets, &model.Target{
			URL:          url,
			LastStatus:   status.status,
			LastResponse: status.moment,
			Period:       period,
		})
	}

	slog.Debug("return info on all targets")
	return targets
}

func (s *storage) RemoveTarget(url string) error {
	s.paramsMx.Lock()
	defer s.paramsMx.Unlock()

	_, ok := s.paramsMap[url]
	if !ok {
		return model.ErrUrlNotExists
	}

	slog.Debug("remove target by url", "url", url)
	delete(s.paramsMap, url)
	return nil
}
