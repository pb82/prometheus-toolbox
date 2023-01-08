package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"github.com/pb82/prometheus-toolbox/api"
	"github.com/pb82/prometheus-toolbox/pkg/precalculated"
	"github.com/pb82/prometheus-toolbox/pkg/proxy"
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
	DefaultPrometheusUrl   = "http://localhost:9090"
	DefaultConfigFile      = "./config.yml"
	DefaultBatchSize       = 500
	DefaultProxyListenPort = 3241
)

var (
	prometheusUrl   *string
	configFile      *string
	printVersion    *bool
	batchSize       *int
	proxyListen     *bool
	proxyListenPort *int
	environment     *bool
	initialize      *bool
)

var (
	//go:embed environment.sh
	environmentSetupScript string

	//go:embed config.yml
	exampleConfig string
)

func main() {
	flag.Parse()

	if printVersion != nil && *printVersion {
		fmt.Printf("Prometheus toolbox v%v", version.Version)
		fmt.Println()
		os.Exit(0)
	}

	if environment != nil && *environment {
		fmt.Println(environmentSetupScript)
		os.Exit(0)
	}

	if initialize != nil && *initialize {
		fmt.Println(exampleConfig)
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
	if err != nil {
		fmt.Println(fmt.Sprintf("error sending samples: %v", err.Error()))
		os.Exit(1)
	}

	if len(requests) > 0 {
		log.Printf("sending %v samples in %v requests (max batch size is %v)", samples, len(requests), *batchSize)

		for i, request := range requests {
			err = remotewrite.SendWriteRequest(request, parsedPrometheusUrl)
			if err != nil {
				log.Fatalf("error sending batch: %v", err.Error())
			}
			log.Printf("successfully sent batch %v/%v", i+1, len(requests))
		}

		log.Printf("done sending precalculated series")
	} else {
		log.Println("no precalculated series")
	}

	wg := &sync.WaitGroup{}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGPIPE, syscall.SIGABRT)
	defer stop()

	err = stream.StartStreamWriters(ctx, config, parsedPrometheusUrl, wg)
	if err != nil {
		log.Fatalf("error starting stream writer: %v", err.Error())
	}

	if *proxyListen {
		err := proxy.StartListener(ctx, proxyListenPort)
		if err != nil {
			log.Fatalf("error starting proxy listener: %v", err.Error())
		}
	}

	wg.Wait()
}

func init() {
	prometheusUrl = flag.String("prometheus.url", "", "prometheus base url")
	configFile = flag.String("config.file", DefaultConfigFile, "config file location")
	batchSize = flag.Int("batch.size", DefaultBatchSize, "max number of samples per remote write request")
	printVersion = flag.Bool("version", false, "print version and exit")
	proxyListen = flag.Bool("proxy.listen", false, "receive remote write requests")
	proxyListenPort = flag.Int("proxy.listen.port", DefaultProxyListenPort, "port to receive remote write requests")
	environment = flag.Bool("environment", false, "print environment setup script and exit")
	initialize = flag.Bool("init", false, "print sample config file and exit")
}
