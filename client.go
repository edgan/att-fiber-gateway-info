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

func NewGatewayClient(baseURL string, colorMode bool, cookieFile string, debug bool, freshCookies bool, loginPath string) (*GatewayClient, error) {
	// Create and load cookies from file if applicable
	jar, loadedCookies, err := createAndLoadCookies(baseURL, cookieFile, debug, freshCookies)
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
		baseURL:       baseURL,
		colorMode:     colorMode,
		cookieFile:    cookieFile,
		loadedCookies: loadedCookies,
		loginPath:     loginPath,
	}, nil
}

// createGatewayClient creates a new client for interacting with the gateway
func createGatewayClient(baseURL string, colorMode bool, cookieFile string, debug bool, freshCookies bool) (*GatewayClient, error) {
	loginPath := returnPath("login")
	client, err := NewGatewayClient(baseURL, colorMode, cookieFile, debug, freshCookies, loginPath)
	return client, err
}
