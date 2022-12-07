package main

import (
	"flag"
	"fmt"
	"github.com/pb82/prometheus-toolbox/api"
	"github.com/pb82/prometheus-toolbox/pkg/precalculated"
	"github.com/pb82/prometheus-toolbox/pkg/remotewrite"
	"github.com/pb82/prometheus-toolbox/pkg/stream"
	"github.com/pb82/prometheus-toolbox/version"
	"log"
	"net/url"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"
)

const (
	DefaultConfigFile = "./config.yml"
	DefaultBatchSize  = 500
)

var (
	prometheusUrl *string
	configFile    *string
	printVersion  *bool
	batchSize     *int
)

func main() {
	flag.Parse()

	if printVersion != nil && *printVersion {
		fmt.Printf("Prometheus toolbox v%v", version.Version)
		fmt.Println()
		os.Exit(0)
	}

	bytes, err := os.ReadFile(*configFile)
	if err != nil {
		fmt.Println(fmt.Sprintf("error reading config file: %v", err.Error()))
		os.Exit(1)
	}

	config, err := api.FromYaml(bytes)
	if err != nil {
		fmt.Println(fmt.Sprintf("error parsing config file: %v", err.Error()))
		os.Exit(1)
	}

	if prometheusUrl == nil || *prometheusUrl == "" {
		fmt.Println("missing prometheus base url, make sure to set --prometheus.url")
		os.Exit(1)
	}

	parsedPrometheusUrl, err := url.Parse(*prometheusUrl)
	if err != nil {
		fmt.Println(fmt.Sprintf("error parsing prometheus base url: %v", err.Error()))
		os.Exit(1)
	}
	parsedPrometheusUrl.Path = path.Join(parsedPrometheusUrl.Path, "/api/v1/write")

	requests, samples, err := precalculated.SchedulePrecalculatedRemoteWriteRequests(config, *batchSize)
	log.Printf("sending %v samples in %v requests (max batch size is %v)", samples, len(requests), *batchSize)

	for i, request := range requests {
		err = remotewrite.SendWriteRequest(request, parsedPrometheusUrl)
		if err != nil {
			log.Fatalf("error sending batch: %v", err.Error())
		}
		log.Printf("successfully sent batch %v/%v", i+1, len(requests))
	}

	log.Printf("done sending precalculated series")

	wg := &sync.WaitGroup{}
	stop := make(chan bool)
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGINT)
	go func() {
		sig := <-sigs
		switch sig {
		case syscall.SIGTERM, syscall.SIGABRT, syscall.SIGINT:
			log.Println("stop signal received")
			close(stop)
		}
	}()

	count, err := stream.StartStreamWriters(config, parsedPrometheusUrl, wg, stop)
	if err != nil {
		log.Fatalf("error starting stream writer: %v", err.Error())
	}
	if count > 0 {
		wg.Wait()
	} else {
		close(stop)
	}
}

func init() {
	prometheusUrl = flag.String("prometheus.url", "", "prometheus base url")
	configFile = flag.String("config.file", DefaultConfigFile, "config file location")
	batchSize = flag.Int("batch.size", DefaultBatchSize, "max number of samples per remote write request")
	printVersion = flag.Bool("version", false, "print version and exit")
}
