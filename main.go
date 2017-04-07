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

func sendMetrics() error {
	g, err := graphite.NewGraphite(graphiteAddr, graphitePort)
	if err != nil {
		return errors.WithStack(err)
	}
	defer g.Disconnect()

	t := time.Now()
	metrics := buildMetrics(t)

	return g.SendMetrics(metrics)
}

func randomDice() int {
	n, err := rand.Int(rand.Reader, big.NewInt(6))
	if err != nil {
		log.Printf("error in rand.Int, err=%+v", errors.WithStack(err))
	}
	return int(n.Int64()) + 1
}

func run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Printf("canceled, exiting")
			return nil
		case <-time.Tick(5 * time.Second):
			err := sendMetrics()
			if err != nil {
				log.Printf("error in sendMetrics, err=%+v", errors.WithStack(err))
			}
		}
	}
}

var serverID string
var graphiteAddr string
var graphitePort int
var siteCount int

func main() {
	flag.StringVar(&graphiteAddr, "graphite-addr", "localhost", "graphite TCP address")
	flag.IntVar(&graphitePort, "graphite-port", 2003, "graphite TCP port")
	flag.StringVar(&serverID, "server-id", "sv01", "server ID")
	flag.IntVar(&siteCount, "site-count", 50, "site count")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Printf("Got signal")
		cancel()
		log.Printf("called cancel")
	}()

	err := run(ctx)
	if err != nil {
		panic(err)
	}
}
