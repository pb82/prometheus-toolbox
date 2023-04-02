package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"github.com/pb82/prometheus-toolbox/api"
	"github.com/pb82/prometheus-toolbox/internal"
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
	DefaultConfigFile          = "./config.yml"
	DefaultBatchSize           = 500
	DefaultProxyListenPort     = 3241
	DefaultRemoteWriteEndpoint = "/api/v1/write"
)

var (
	prometheusUrl     *string
	configFile        *string
	printVersion      *bool
	batchSize         *int
	proxyListen       *bool
	proxyListenPort   *int
	environment       *bool
	initialize        *bool
	oidcClientId      *string
	oidcClientSecret  *string
	oidcIssuerUrl     *string
	oidcAudience      *string
	oidcEnabled       *bool
	remoteWriteSuffix *string
	rules             *bool
)

var (
	//go:embed environment.sh
	environmentSetupScript string

	//go:embed config.yml
	exampleConfig string

	//go:embed rules.yml
	exampleRules string
)

func main() {
	flag.Parse()

	if *printVersion {
		fmt.Printf("Prometheus toolbox v%v", version.Version)
		fmt.Println()
		os.Exit(0)
	}

	if *environment {
		fmt.Println(environmentSetupScript)
		os.Exit(0)
	}

	if *initialize {
		fmt.Println(exampleConfig)
		os.Exit(0)
	}

	if *rules {
		fmt.Println(exampleRules)
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

	if *prometheusUrl == "" {
		fmt.Println("missing prometheus base url, make sure to set --prometheus.url")
		os.Exit(1)
	}

	parsedPrometheusUrl, err := url.Parse(*prometheusUrl)
	if err != nil {
		fmt.Println(fmt.Sprintf("error parsing prometheus base url: %v", err.Error()))
		os.Exit(1)
	}
	parsedPrometheusUrl.Path = path.Join(parsedPrometheusUrl.Path, *remoteWriteSuffix)

	requests, samples, err := precalculated.SchedulePrecalculatedRemoteWriteRequests(config, *batchSize)
	if err != nil {
		fmt.Println(fmt.Sprintf("error sending samples: %v", err.Error()))
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGPIPE, syscall.SIGABRT)
	defer stop()

	remoteWriter, err := buildRemoteWriter(ctx, parsedPrometheusUrl)
	if err != nil {
		panic(err)
	}

	if len(requests) > 0 {
		log.Printf("sending %v samples in %v requests (max batch size is %v)", samples, len(requests), *batchSize)

		for i, request := range requests {
			err = remoteWriter.SendWriteRequest(request)
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
	err = stream.StartStreamWriters(ctx, config, remoteWriter, wg)
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

func buildRemoteWriter(ctx context.Context, prometheusUrl *url.URL) (*remotewrite.RemoteWriter, error) {
	if *oidcEnabled {
		oidcConfig := internal.NewOIDCConfig(*oidcClientId, *oidcClientSecret, *oidcIssuerUrl, *oidcAudience)
		oidcConfig.Validate()
		return remotewrite.NewRemoteWriterWithOIDCTransport(ctx, prometheusUrl, oidcConfig)
	} else {
		return remotewrite.NewRemoteWriter(prometheusUrl), nil
	}
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
	oidcClientId = flag.String("oidc.clientId", "", "oidc client id")
	oidcClientSecret = flag.String("oidc.clientSecret", "", "oidc client secret")
	oidcIssuerUrl = flag.String("oidc.issuer", "", "oidc token issuer url")
	oidcAudience = flag.String("oidc.audience", "", "oidc audience")
	oidcEnabled = flag.Bool("oidc.enabled", false, "enable oidc token authentication")
	remoteWriteSuffix = flag.String("prometheus.url.suffix", DefaultRemoteWriteEndpoint, "allows alternate remote write endpoints")
	rules = flag.Bool("rules", false, "print sample alerting rules file and exit")
}
