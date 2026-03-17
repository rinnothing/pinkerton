package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rinnothing/pinkerton/internal/model"
	"github.com/rinnothing/pinkerton/pkg/checks"
	"github.com/rinnothing/pinkerton/pkg/client"
)

func writeTarget(w io.Writer, tgt *model.Target) error {
	_, err := fmt.Fprintf(w, "{URL: %s, LastStatus: %d, LastResponse: %v, Period: %v}", tgt.URL, tgt.LastStatus, tgt.LastResponse, tgt.Period)
	return err
}

func stringTarget(tgt *model.Target) string {
	b := new(strings.Builder)
	err := writeTarget(b, tgt)
	if err != nil {
		panic(fmt.Errorf("can't convert target to string: %w", err))
	}
	return b.String()
}

var allowedCalls = map[string]struct{}{"get": {}, "getall": {}, "add": {}, "update": {}, "remove": {}, "health": {}}
var urlCalls = map[string]struct{}{"get": {}, "add": {}, "update": {}, "remove": {}}
var periodCalls = map[string]struct{}{"add": {}, "update": {}}

func main() {
	host := flag.String("host", "localhost", "pinkerton host")
	port := flag.String("port", "3233", "pinkerton port")
	timeout := flag.Duration("timeout", time.Second*5, "pinkerton client timeout")

	var call string
	flag.Func("call", "type of call to pinkerton: get, getall, add, update, remove, health", func(s string) error {
		_, ok := allowedCalls[s]
		if !ok {
			return fmt.Errorf("no such call type %s", s)
		}

		call = s
		return nil
	})

	tgtUrl := flag.String("url", "", "target url (needed for get, add, update and remove)")
	period := flag.Duration("period", time.Second*10, "target update period (needed for add and update)")

	flag.Parse()

	if call == "" {
		fmt.Fprint(os.Stderr, "call type should be specified")
		os.Exit(1)
	}

	if _, ok := urlCalls[call]; ok {
		err := checks.CheckUrl(*tgtUrl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "incorrect (or not specified) url %q: %s\n", *tgtUrl, err)
			os.Exit(1)
		}
	}

	if _, ok := periodCalls[call]; ok {
		err := checks.CheckPeriod(*period)
		if err != nil {
			fmt.Fprintln(os.Stderr, *period)
			os.Exit(1)
		}
	}

	req := &model.TargetRequest{URL: *tgtUrl, Period: *period}

	ctx := context.Background()
	cl := client.New(http.DefaultClient, *host, *port, *timeout)
	switch call {
	case "health":
		err := cl.Health(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "health error: %s", err)
			return
		}
		fmt.Println("health is ok")
	case "getall":
		tgts, err := cl.GetAllTargets(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "getall error: %s", err)
			return
		}

		b := new(strings.Builder)
		for i, tgt := range tgts {
			if i != 0 {
				b.WriteString(",\n ")
			}
			if err := writeTarget(b, tgt); err != nil {
				panic(fmt.Errorf("can't convert array to string: %w", err))
			}
		}
		fmt.Printf("targets are:\n[%s]\n", b.String())
	case "get":
		tgt, err := cl.GetTarget(ctx, *tgtUrl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "get %q error: %s", *tgtUrl, err)
			return
		}
		fmt.Printf("target is: %s\n", stringTarget(tgt))
	case "remove":
		err := cl.RemoveTarget(ctx, *tgtUrl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "remove %q error: %s", *tgtUrl, err)
			return
		}
		fmt.Println("removed successfuly")
	case "add":
		err := cl.AddTarget(ctx, req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "add %v error: %s", req, err)
			return
		}
		fmt.Println("added successfuly")
	case "update":
		err := cl.UpdateTarget(ctx, req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "update %v error: %s", req, err)
			return
		}
		fmt.Println("updated successfuly")
	default:
		panic("no such call")
	}
}
