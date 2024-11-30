package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"strings"

	"edgan/att-fiber-gateway-info/internal/logging"
)

func calculateHash(configs configs, nonce string) string {
	// Replicate JavaScript hex_md5(password + nonce)
	hasher := md5.New()
	io.WriteString(hasher, configs.Password+nonce)
	return hex.EncodeToString(hasher.Sum(nil))
}

func (rc *gatewayClient) login(configs configs, flags *flags) error {
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
		"password":     {strings.Repeat(star, len(configs.Password))}, // Replicate JS behavior
		"hashpassword": {hash},
		"Continue":     {"Continue"}, // submit button
	}

	if err := rc.postForm(flags, formData, rc.loginPath); err != nil {
		logging.LogFatalf("Submission to %s failed\n%v", rc.loginPath, err)
	}

	return nil
}

// performLogin performs the login action for the client
func performLogin(client *gatewayClient, configs configs, flags *flags) {
	if err := client.login(configs, flags); err != nil {
		logging.LogFatalf("Login failed: %v", err)
	}
}
