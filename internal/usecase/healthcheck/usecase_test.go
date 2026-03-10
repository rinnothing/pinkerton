package healthcheck_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rinnothing/pinkerton/internal/model"
	"github.com/rinnothing/pinkerton/internal/repository/storage"
	"github.com/rinnothing/pinkerton/internal/usecase/healthcheck"
)

const num_threads = 10

var _ healthcheck.Pinger = &constantPinger{}

type constantPinger struct {
	code int
	err  error
}

func (p *constantPinger) Ping(url string) (int, error) {
	return p.code, p.err
}

func newHealthcheck(initModels []model.Target, pinger healthcheck.Pinger) healthcheck.Usecase {
	return healthcheck.New(context.Background(), initModels, num_threads, pinger, storage.New())
}

var targets = []model.Target{
	{URL: "https://google.com", Period: time.Millisecond * 10},
	{URL: "https://youtube.com", Period: time.Minute * 30},
}

func TestInitialModels(t *testing.T) {
	t.Parallel()

	e := newHealthcheck(targets, &constantPinger{code: 200})

	for _, trg := range targets {
		resTrg, err := e.GetTarget(trg.URL)
		if err != nil {
			t.Fatalf("retrieving target %s shouldn't result in error %s", trg.URL, err)
		}

		if resTrg.URL != trg.URL {
			t.Fatalf("retrieved model url %s not equal to requested %s", resTrg.URL, trg.URL)
		}

		if resTrg.Period != trg.Period {
			t.Fatalf("retrieved model has different period %s from original %s", resTrg.Period, trg.Period)
		}
	}
}

func TestAddTarget(t *testing.T) {
	t.Parallel()

	e := newHealthcheck([]model.Target{}, &constantPinger{code: 200})

	err := e.AddTarget(&targets[0])
	if err != nil {
		t.Fatalf("adding new target shouldn't result in error %s", err)
	}

	trg, err := e.GetTarget(targets[0].URL)
	if err != nil {
		t.Fatalf("retrieving target %s shouldn't result in error %s", targets[0].URL, err)
	}

	if trg.URL != targets[0].URL {
		t.Fatalf("retrieved target url %s not equal to requested %s", trg.URL, targets[0].URL)
	}

	if trg.Period != targets[0].Period {
		t.Fatalf("retrieved target has different period %s from original %s", trg.Period, targets[0].Period)
	}

	_, err = e.GetTarget(targets[1].URL)
	if !errors.Is(err, model.ErrUrlNotExists) {
		t.Fatal("quering nonexisting target should result in error")
	}

	err = e.AddTarget(&targets[0])
	if !errors.Is(err, model.ErrUrlExists) {
		t.Fatal("adding duplicate target should result in error")
	}
}

func TestUpdateTarget(t *testing.T) {
	t.Parallel()

	e := newHealthcheck([]model.Target{}, &constantPinger{code: 200})

	err := e.AddTarget(&targets[0])
	if err != nil {
		t.Fatalf("adding new target shouldn't result in error %s", err)
	}

	var trgCp model.Target = targets[0]
	trgCp.Period = time.Hour

	err = e.UpdateTarget(&trgCp)
	if err != nil {
		t.Fatalf("updating existing target shouldn't result in error %s", err)
	}

	resTrg, err := e.GetTarget(targets[0].URL)
	if err != nil {
		t.Fatalf("quering existing target %s shouldn't result in error %s", targets[0].URL, err)
	}

	if resTrg.Period != trgCp.Period {
		t.Fatalf("should update target period to %s, but it remained %s", trgCp.Period, resTrg.Period)
	}

	trgCp = targets[1]
	err = e.UpdateTarget(&trgCp)
	if !errors.Is(err, model.ErrUrlNotExists) {
		t.Fatalf("updating nonexisting target should result in error")
	}
}

func TestRemoveTarget(t *testing.T) {
	t.Parallel()

	e := newHealthcheck([]model.Target{}, &constantPinger{code: 200})

	err := e.AddTarget(&targets[0])
	if err != nil {
		t.Fatalf("adding new target shouldn't result in error %s", err)
	}

	_, err = e.GetTarget(targets[0].URL)
	if err != nil {
		t.Fatalf("quering existing target %s shouldn't result in error %s", targets[0].URL, err)
	}

	err = e.RemoveTarget(targets[0].URL)
	if err != nil {
		t.Fatalf("removing existing target %s shouldn't result in error %s", targets[0].URL, err)
	}

	err = e.RemoveTarget(targets[1].URL)
	if !errors.Is(err, model.ErrUrlNotExists) {
		t.Fatalf("removing nonexisting target should result in error")
	}
}

// func TestCodes(t *testing.T) {
// 	errCode := 404
// 	e := newHealthcheck([]model.Target{}, &constantPinger{code: errCode})

// 	cpTgt := targets[0]
// 	cpTgt.Period = 0
// 	err := e.AddTarget(&cpTgt)
// 	if err != nil {
// 		t.Fatalf("adding new target shouldn't result in error %s", err)
// 	}

// 	time.Sleep(time.Millisecond * 20)

// 	trg, err := e.GetTarget(targets[0].URL)
// 	if err != nil {
// 		t.Fatalf("quering existing target shouldn't result in error %s", err)
// 	}

// 	nw := time.Now()
// 	if trg.LastResponse.After(nw) {
// 		t.Fatalf("seem to get last response in future %s (now is %s)", trg.LastResponse, nw)
// 	}

// 	if trg.LastStatus != errCode {
// 		t.Fatalf("response seem to not update, should be %d, but is %d", errCode, trg.LastStatus)
// 	}
// }
