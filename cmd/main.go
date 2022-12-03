package main

import (
	"flag"
	"fmt"
	"os"
	"promtoolbox/api"
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

	fmt.Println(config)
}

func init() {
	prometheusUrl = flag.String("prometheus.url", "", "prometheus base url")
	configFile = flag.String("config.file", DefaultConfigFile, "config file location")
	batchSize = flag.Int("batch.size", DefaultBatchSize, "max number of samples per remote write request")
	printVersion = flag.Bool("version", false, "print version and exit")
}
