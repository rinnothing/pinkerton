//go:build integration

package integration_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/rinnothing/pinkerton/config"
	"github.com/rinnothing/pinkerton/internal/controller"
	"github.com/rinnothing/pinkerton/pkg/client"
)

var cfg *config.Config
var cl *client.Client

func createHealthHandler(ctx context.Context, host, port string, status int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		io.Copy(io.Discard, r.Body)

		w.WriteHeader(status)
	})

	server := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		server.Close()
	}()
	go func() {
		server.ListenAndServe()
	}()
}

var targets = []*client.TargetRequest{
	{URL: "http://localhost:3235", Period: time.Second * 5},
	{URL: "http://localhost:3236", Period: time.Second * 5},
}

var targetsURLs []string

func init() {
	targetsURLs = make([]string, len(targets))
	for i, tgt := range targets {
		targetsURLs[i] = tgt.URL
	}
	slices.Sort(targetsURLs)
}

func TestMain(m *testing.M) {
	_, external := os.LookupEnv("EXTERNAL")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg = config.ReadConfig("integration_server.json")

	cl = client.New(http.DefaultClient, cfg.Host, cfg.Port, time.Second*5)

	if !external {
		go func() {
			controller.Run(ctx, cfg)
		}()
	}
	for _, tgt := range targets {
		createHealthHandler(ctx, "", strings.TrimSuffix(tgt.URL, "http://localhost:"), 200)
	}

	fmt.Println("ready to run")
	m.Run()
}

func TestHappyPath(t *testing.T) {
	t.Log("starting test")
	for cl.Health(t.Context()) != nil {
		t.Log("check health, bad")
	}
	t.Log("check health, good")

	for _, tgt := range targets {
		err := cl.AddTarget(t.Context(), tgt)
		if err != nil {
			t.Fatalf("can't add server %s to watchlist: %s", tgt.URL, err)
		}
	}

	err := cl.AddTarget(t.Context(), targets[0])
	if !errors.Is(err, client.ErrUrlExists) {
		t.Fatalf("shouldn't be able to add duplicate %s server to watchlist", targets[0].URL)
	}

	badReq := &client.TargetRequest{URL: "biba_and_boba.com", Period: time.Minute * 5}
	err = cl.AddTarget(t.Context(), badReq)
	if !errors.Is(err, client.ErrBadRequest) {
		t.Fatalf("shouldn't be able to add bad url %s server to watchlist", badReq.URL)
	}

	badReq = &client.TargetRequest{URL: "http://google.com", Period: time.Minute * 0}
	err = cl.AddTarget(t.Context(), badReq)
	if !errors.Is(err, client.ErrBadRequest) {
		t.Fatalf("shouldn't be able to add bad period %v server to watchlist", badReq.Period)
	}

	for _, tgt := range targets {
		resTgt, err := cl.GetTarget(t.Context(), tgt.URL)
		if err != nil {
			t.Fatalf("should be able to query added server %s: %s", tgt.URL, err)
		}

		if resTgt.URL != tgt.URL || resTgt.Period != tgt.Period {
			t.Fatalf("result and added url and period should be the same: %v != %v", resTgt, tgt)
		}
	}

	tgts, err := cl.GetAllTargets(t.Context())
	if err != nil {
		t.Fatal("should be able to retrieve all targets")
	}

	urlList := make([]string, len(tgts))
	for i, tgt := range tgts {
		urlList[i] = tgt.URL
	}
	slices.Sort(urlList)

	if !reflect.DeepEqual(targetsURLs, urlList) {
		t.Fatalf("retrieved and added urls should be equal: %v != %v", targetsURLs, urlList)
	}

	err = cl.RemoveTarget(t.Context(), targets[0].URL)
	if err != nil {
		t.Fatalf("should be able to remove added target %s: %s", targets[0].URL, err)
	}

	err = cl.RemoveTarget(t.Context(), targets[0].URL)
	if !errors.Is(err, client.ErrUrlNotExists) {
		t.Fatalf("shouldn't be able to remove already remove server %s", targets[0].URL)
	}

	targetCp := *targets[1]
	targetCp.Period = time.Minute * 5
	err = cl.UpdateTarget(t.Context(), &targetCp)
	if err != nil {
		t.Fatalf("should be able to update added target %s: %s", targetCp.URL, err)
	}

	resTgt, err := cl.GetTarget(t.Context(), targetCp.URL)
	if err != nil {
		t.Fatalf("should be able to get added target %s: %s", targetCp.URL, err)
	}

	if resTgt.URL != targetCp.URL || resTgt.Period != targetCp.Period {
		t.Fatalf("result and added url and period should be the same: %v != %v", resTgt, targetCp)
	}

	err = cl.UpdateTarget(t.Context(), targets[0])
	if !errors.Is(err, client.ErrUrlNotExists) {
		t.Fatalf("shouldn't be able to update removed target %s", targets[0].URL)
	}

	err = cl.RemoveTarget(t.Context(), targets[1].URL)
	if err != nil {
		t.Fatalf("should be able to remove added target %s: %s", targets[1].URL, err)
	}

	err = cl.Health(t.Context())
	if err != nil {
		t.Fatalf("health of working pinkerton should be ok")
	}
}
