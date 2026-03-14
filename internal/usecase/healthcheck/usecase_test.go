package healthcheck_test

import (
	"context"
	"errors"
	"slices"
	"strings"
	"sync"
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

type storageNotifier struct {
	healthcheck.Storage
	once   *sync.Once
	notify chan struct{}
}

func (sn *storageNotifier) StoreStatus(url string, status int) {
	sn.Storage.StoreStatus(url, status)

	sn.once.Do(func() {
		close(sn.notify)
	})
}

func newStorageNotifier() *storageNotifier {
	return &storageNotifier{
		storage.New(),
		new(sync.Once),
		make(chan struct{}),
	}
}

func newHealthcheck(initModels []model.Target, pinger healthcheck.Pinger) (healthcheck.Usecase, chan struct{}) {
	sn := newStorageNotifier()

	tgtMdls := make([]model.TargetRequest, len(initModels))
	for i, mdl := range initModels {
		tgtMdls[i] = model.TargetRequest{URL: mdl.URL, Period: mdl.Period}
	}
	return healthcheck.New(context.Background(), tgtMdls, num_threads, pinger, sn), sn.notify
}

var targets = []model.Target{
	{URL: "https://google.com", Period: time.Millisecond * 10},
	{URL: "https://youtube.com", Period: time.Minute * 30},
}

func TestInitialModels(t *testing.T) {
	t.Parallel()

	e, _ := newHealthcheck(targets, &constantPinger{code: 200})

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

	e, _ := newHealthcheck([]model.Target{}, &constantPinger{code: 200})

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

func TestGetAllTargets(t *testing.T) {
	t.Parallel()

	e, _ := newHealthcheck([]model.Target{}, &constantPinger{code: 200})

	addedTargets := make([]*model.Target, len(targets))
	for i, tgt := range targets {
		err := e.AddTarget(&tgt)
		if err != nil {
			t.Fatalf("adding new target shouldn't result in error %s", err)
		}
		addedTargets[i] = &tgt
	}

	tgtCmp := func(a, b *model.Target) int {
		return strings.Compare(a.URL, b.URL)
	}

	tgts := e.GetAllTargets()

	slices.SortFunc(tgts, tgtCmp)
	slices.SortFunc(addedTargets, tgtCmp)

	if len(addedTargets) != len(tgts) {
		t.Fatalf("lengths of target lists should be equal but added is %d and read is %d", len(addedTargets), len(tgts))
	}
	for i := range addedTargets {
		if *addedTargets[i] != *tgts[i] {
			t.Errorf("added and read targets should be equal, but they are %v and %v", *addedTargets[i], *tgts[i])
		}
	}
}

func TestUpdateTarget(t *testing.T) {
	t.Parallel()

	e, _ := newHealthcheck([]model.Target{}, &constantPinger{code: 200})

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

	e, _ := newHealthcheck([]model.Target{}, &constantPinger{code: 200})

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

func TestCodes(t *testing.T) {
	t.Parallel()

	errCode := 404
	e, notify := newHealthcheck([]model.Target{}, &constantPinger{code: errCode})

	err := e.AddTarget(&targets[0])
	if err != nil {
		t.Fatalf("adding new target shouldn't result in error %s", err)
	}

	t.Log("started sleep")
	<-notify
	t.Log("ended sleep")

	trg, err := e.GetTarget(targets[0].URL)
	if err != nil {
		t.Fatalf("quering existing target shouldn't result in error %s", err)
	}

	nw := time.Now()
	if trg.LastResponse.After(nw) {
		t.Fatalf("seem to get last response in future %s (now is %s)", trg.LastResponse, nw)
	}

	if trg.LastStatus != errCode {
		t.Fatalf("response seem to not update, should be %d, but is %d", errCode, trg.LastStatus)
	}
}
