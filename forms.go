package main

import (
	"fmt"
	"io"
	"log"
	"net/url"

	"edgan/att-fiber-gateway-info/internal/logging"
)

func (rc *gatewayClient) postForm(flags *flags, formData url.Values, path string) error {
	// Submit form
	resp, err1 := rc.client.PostForm(rc.baseURL+path, formData)

	if err1 != nil {
		return fmt.Errorf("failed to submit the form to %s: %v", path, err1)
	}

	defer resp.Body.Close()

	if path == rc.loginPath {
		// Save session cookies
		if err := rc.saveSessionCookies(); err != nil {
			log.Printf("Failed to save cookies: %v", err1)
		}
	}

	body, err2 := io.ReadAll(resp.Body)

	if err2 != nil {
		logging.DebugLog(*flags.Debug, "Failed to ReadAll")
		return fmt.Errorf("failed to read response for path %s: %v", path, err2)
	}

	bodyStr := string(body)

	err3 := checkForLoginFailure(bodyStr, flags, rc.loginPath)

	if err3 != nil {
		return err3
	}

	return nil
}

func (rc *gatewayClient) submitForm(
	action string, flags *flags, page string, path string,
	resetAction [actionAttributes]string,
) error {
	buttonName := resetAction[zero]
	buttonValue := resetAction[one]
	question := resetAction[two]
	warning := resetAction[three]

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
