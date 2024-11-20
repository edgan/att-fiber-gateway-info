package main

import (
	"crypto/tls"
	"net/http"
)

type GatewayClient struct {
	client        *http.Client
	baseURL       string
	colorMode     bool
	cookieFile    string
	loadedCookies int
	loginPath     string
}

func NewGatewayClient(configs Configs, colorMode bool, flags *Flags, loginPath string) (*GatewayClient, error) {
	// Create and load cookies from file if applicable
	jar, loadedCookies, err := createAndLoadCookies(configs, flags)
	if err != nil {
		return nil, err
	}

	// Configure transport to skip SSL verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// Create HTTP client
	client := &http.Client{
		Transport: tr,
		Jar:       jar,
	}

	return &GatewayClient{
		client:        client,
		baseURL:       configs.BaseURL,
		colorMode:     colorMode,
		cookieFile:    *flags.CookieFile,
		loadedCookies: loadedCookies,
		loginPath:     loginPath,
	}, nil
}

// createGatewayClient creates a new client for interacting with the gateway
func createGatewayClient(configs Configs, colorMode bool, flags *Flags) (*GatewayClient, error) {
	loginPath := returnPath("login")
	client, err := NewGatewayClient(configs, colorMode, flags, loginPath)
	return client, err
}
