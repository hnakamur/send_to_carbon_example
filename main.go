package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"strings"
	"time"

	graphite "github.com/marpaia/graphite-golang"
	"github.com/pkg/errors"
)

func newMetric(t time.Time, name, value string) graphite.Metric {
	return graphite.Metric{
		Name:      name,
		Value:     value,
		Timestamp: t.Unix(),
	}
}

func buildMetrics(t time.Time) []graphite.Metric {
	var metrics []graphite.Metric
	for i := 0; i < siteCount; i++ {
		domain := fmt.Sprintf("ex%d.example.jp", i+1)
		name := fmt.Sprintf("local.random.diceroll.%s.%s",
			strings.Replace(domain, ".", "_", -1), serverID)
		value := fmt.Sprintf("%d", randomDice())
		metrics = append(metrics, newMetric(t, name, value))
	}
	return metrics
}

func sendMetrics(t time.Time) error {
	g, err := graphite.NewGraphite(graphiteAddr, graphitePort)
	if err != nil {
		return errors.WithStack(err)
	}
	defer g.Disconnect()

	// log.Printf("t=%s", t)
	metrics := buildMetrics(t)

	if len(metrics) > 0 {
		err := g.SendMetrics(metrics)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func randomDice() int {
	n, err := rand.Int(rand.Reader, big.NewInt(6))
	if err != nil {
		log.Printf("error in rand.Int, err=%+v", errors.WithStack(err))
	}
	return int(n.Int64()) + 1
}

func nextRoundTime(t time.Time, d time.Duration) time.Time {
	next := t.Round(d)
	if next.Before(t) {
		next = next.Add(d)
	}
	return next
}

func run(ctx context.Context) error {
	now := time.Now()
	until := nextRoundTime(now, interval)
	d := until.Sub(now)
	log.Printf("%s: sleep %s until %s", serverID, d, until)
	select {
	case <-ctx.Done():
		log.Printf("%s: canceled, exiting", serverID)
		return nil
	case t := <-time.After(d):
		err := sendMetrics(t)
		if err != nil {
			log.Printf("%s: error in sendMetrics, err=%+v", serverID, err)
		} else {
			log.Printf("%s: elapsed for sendMetrics: %s", serverID, time.Since(t))
		}
	}
	for {
		select {
		case <-ctx.Done():
			log.Printf("%s: canceled, exiting", serverID)
			return nil
		case t := <-time.Tick(interval):
			err := sendMetrics(t)
			if err != nil {
				log.Printf("%s: error in sendMetrics, err=%+v", serverID, err)
			} else {
				log.Printf("%s: elapsed for sendMetrics: %s", serverID, time.Since(t))
			}
		}
	}
}

var interval time.Duration
var serverID string
var graphiteAddr string
var graphitePort int
var siteCount int

func main() {
	flag.DurationVar(&interval, "interval", time.Minute, "metrics interval")
	flag.StringVar(&graphiteAddr, "graphite-addr", "localhost", "graphite TCP address")
	flag.IntVar(&graphitePort, "graphite-port", 2003, "graphite TCP port")
	flag.StringVar(&serverID, "server-id", "sv01", "server ID")
	flag.IntVar(&siteCount, "site-count", 50, "site count")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Printf("%s: Got signal", serverID)
		cancel()
		log.Printf("%s: called cancel", serverID)
	}()

	err := run(ctx)
	if err != nil {
		panic(err)
	}
}
