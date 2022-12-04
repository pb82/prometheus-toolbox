package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"promtoolbox/api"
	"promtoolbox/pkg/precalculated"
	"promtoolbox/pkg/remotewrite"
	"promtoolbox/version"
)

const (
	DefaultConfigFile = "./example/config.yml"
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
}

func init() {
	prometheusUrl = flag.String("prometheus.url", "", "prometheus base url")
	configFile = flag.String("config.file", DefaultConfigFile, "config file location")
	batchSize = flag.Int("batch.size", DefaultBatchSize, "max number of samples per remote write request")
	printVersion = flag.Bool("version", false, "print version and exit")
}
