package api

import (
	"github.com/ghodss/yaml"
)

type TimeseriesConfig struct {
	Series string `json:"series"`
	Values string `json:"values"`
	Stream string `json:"stream"`
}

type Config struct {
	Interval string             `json:"interval"`
	Series   []TimeseriesConfig `json:"time_series"`
}

func FromYaml(raw []byte) (*Config, error) {
	config := &Config{}
	err := yaml.Unmarshal(raw, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
