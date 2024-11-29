package main

import (
	"crypto/tls"
	"net/http"
)

type gatewayClient struct {
	client        *http.Client
	baseURL       string
	colorMode     bool
	cookieFile    string
	loadedCookies int
	loginPath     string
}

func newGatewayClient(configs configs, colorMode bool, flags *flags, loginPath string) (*gatewayClient, error) {
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

	return &gatewayClient{
		client:        client,
		baseURL:       configs.BaseURL,
		colorMode:     colorMode,
		cookieFile:    *flags.CookieFile,
		loadedCookies: loadedCookies,
		loginPath:     loginPath,
	}, nil
}

// createGatewayClient creates a new client for interacting with the gateway
func createGatewayClient(configs configs, colorMode bool, flags *flags) (*gatewayClient, error) {
	loginPath := returnPath("login")
	client, err := newGatewayClient(configs, colorMode, flags, loginPath)
	return client, err
}
