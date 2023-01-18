package internal

import (
	"errors"
	"fmt"
	"net/url"
)

// SizeInfo contains information about the size of a remote write request
type SizeInfo struct {
	CompressedSize   float64
	UncompressedSize float64
	TimeseriesCount  float64
}

type OIDCConfig struct {
	ClientId     string
	ClientSecret string
	IssuerUrl    string
	Audience     string
}

func (s *SizeInfo) String() string {
	return fmt.Sprintf("compressed: %v bytes, uncompressed: %v bytes, times series count: %v",
		s.CompressedSize, s.UncompressedSize, s.TimeseriesCount)
}

func (c *OIDCConfig) Validate() {
	_, err := url.Parse(c.IssuerUrl)
	if err != nil {
		panic(err)
	}

	if c.ClientId == "" || c.ClientSecret == "" {
		panic(errors.New("clientId and/or clientSecret missing in OIDC config"))
	}
}

func NewOIDCConfig(clientId string, clientSecret string, issuerUrl string, audience string) *OIDCConfig {
	return &OIDCConfig{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		IssuerUrl:    issuerUrl,
		Audience:     audience,
	}
}
