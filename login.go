package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"
)

func calculateHash(password, nonce string) string {
	// Replicate JavaScript hex_md5(password + nonce)
	hasher := md5.New()
	io.WriteString(hasher, password+nonce)
	return hex.EncodeToString(hasher.Sum(nil))
}

func (rc *GatewayClient) login(password string) error {
	page := "login"

	// Get nonce from page
	nonce, err := rc.getNonce(page)

	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
	}

	// Calculate hash
	hash := calculateHash(password, nonce)

	// Prepare form data
	formData := url.Values{
		"nonce":        {nonce},
		"password":     {strings.Repeat("*", len(password))}, // Replicate JS behavior
		"hashpassword": {hash},
		"Continue":     {"Continue"}, // submit button
	}

	if err := rc.postForm(rc.loginPath, formData); err != nil {
		log.Fatalf("Submission to %s failed: %v", rc.loginPath, err)
	}

	return nil
}

// performLogin performs the login action for the client
func performLogin(client *GatewayClient, password string) {
	if err := client.login(password); err != nil {
		log.Fatalf("Login failed: %v", err)
	}
}
