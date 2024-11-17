package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func (rc *GatewayClient) postForm(path string, formData url.Values) error {
	// Submit form
	resp, err := rc.client.PostForm(rc.baseURL+path, formData)

	if err != nil {
		return fmt.Errorf("failed to submit the form to %s: %v", path, err)
	}

	defer resp.Body.Close()

	// Check if submission was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("submission to %s failed with status: %d", path, resp.StatusCode)
	}

	if path == rc.loginPath {
		// Save session cookies
		if err := rc.saveSessionCookies(); err != nil {
			log.Printf("Failed to save cookies: %v", err)
		}
	}

	return nil
}

func (rc *GatewayClient) submitForm(action string, answerNo bool, answerYes bool, page string, path string, resetAction [4]string) error {
	buttonName := resetAction[0]
	buttonValue := resetAction[1]
	question := resetAction[2]
	warning := resetAction[3]

	if askYesNo(answerNo, answerYes, rc.colorMode, question, warning) {
		nonce, err := rc.getNonce(page)
		if err != nil {
			return fmt.Errorf("failed to get nonce: %v", err)
		}

		formData := url.Values{
			"nonce":    {nonce},
			buttonName: {buttonValue}, // Dynamically use the submit button
		}

		if err := rc.postForm(path, formData); err != nil {
			return fmt.Errorf("submission to %s in %s failed: %v", action, path, err)
		}
	}

	return nil
}
