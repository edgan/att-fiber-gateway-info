package main

import (
	"fmt"
	"io"
	"log"
	"net/url"
)

func (rc *gatewayClient) postForm(flags *flags, formData url.Values, path string) error {
	// Submit form
	resp, err := rc.client.PostForm(rc.baseURL+path, formData)

	if err != nil {
		return fmt.Errorf("failed to submit the form to %s: %v", path, err)
	}

	defer resp.Body.Close()

	if path == rc.loginPath {
		// Save session cookies
		if err := rc.saveSessionCookies(); err != nil {
			log.Printf("Failed to save cookies: %v", err)
		}
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		debugLog(*flags.Debug, "Failed to ReadAll")
		return fmt.Errorf("failed to read response for path %s: %v", path, err)
	}

	bodyStr := string(body)

	error := checkForLoginFailure(bodyStr, flags, rc.loginPath)

	if error != nil {
		return error
	}

	return nil
}

func (rc *gatewayClient) submitForm(action string, flags *flags, page string, path string, resetAction [4]string) error {
	buttonName := resetAction[0]
	buttonValue := resetAction[1]
	question := resetAction[2]
	warning := resetAction[3]

	if askYesNo(rc.colorMode, flags, question, warning) {
		nonce, err := rc.getNonce(page)
		if err != nil {
			return fmt.Errorf("failed to get nonce: %v", err)
		}

		formData := url.Values{
			"nonce":    {nonce},
			buttonName: {buttonValue}, // Dynamically use the submit button
		}

		if err := rc.postForm(flags, formData, path); err != nil {
			return fmt.Errorf("%s in %s failed\n%v", action, path, err)
		}
	}

	return nil
}
