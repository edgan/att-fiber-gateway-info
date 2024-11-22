package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"strings"
)

func calculateHash(configs Configs, nonce string) string {
	// Replicate JavaScript hex_md5(password + nonce)
	hasher := md5.New()
	io.WriteString(hasher, configs.Password+nonce)
	return hex.EncodeToString(hasher.Sum(nil))
}

func (rc *GatewayClient) login(configs Configs) error {
	page := "login"

	// Get nonce from page
	nonce, err := rc.getNonce(page)

	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
	}

	// Calculate hash
	hash := calculateHash(configs, nonce)

	// Prepare form data
	formData := url.Values{
		"nonce":        {nonce},
		"password":     {strings.Repeat("*", len(configs.Password))}, // Replicate JS behavior
		"hashpassword": {hash},
		"Continue":     {"Continue"}, // submit button
	}

	if err := rc.postForm(rc.loginPath, formData); err != nil {
		logFatalf("Submission to %s failed: %v", rc.loginPath, err)
	}

	return nil
}

// performLogin performs the login action for the client
func performLogin(client *GatewayClient, configs Configs) {
	if err := client.login(configs); err != nil {
		logFatalf("Login failed: %v", err)
	}
}
