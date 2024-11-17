package main

import (
	"fmt"
	"io"
	"log"
	"strings"
)

// getPath is a generic function that fetches and returns the response body as a string.
func (rc *GatewayClient) getPath(path string) (string, error) {
	resp, err := rc.client.Get(rc.baseURL + path)

	if err != nil {
		return "", fmt.Errorf("failed to get path %s: %v", path, err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", fmt.Errorf("failed to read response for path %s: %v", path, err)
	}

	bodyStr := string(body)

	// Check for login failure
	if strings.Contains(bodyStr, rc.loginPath) {
		return "", fmt.Errorf("Login failed. Password likely wrong.")
	}

	return bodyStr, nil
}

func (rc *GatewayClient) getPage(action string, answerYes bool, filter string, natActionPrefix string, page string, pretty bool) error {
	path := returnPath(page)

	if page == "reset" {
		mayWarning := "This may take down your internet immediately."
		willWarning := "Note this will take down your internet immediately."

		parts := strings.Split(action, "-")
		task := parts[0]
		resource := parts[1]

		question := formatQuestion(task, resource)

		// The buttonNames and buttonValues have to be exact
		resetActions := map[string][4]string{
			"reset-connection": {"ResetConn", "Reset Connection", question, mayWarning},
			"reset-device":     {"Reset", "Reset Device...", question, mayWarning},
			"reset-firewall":   {"FReset", "Reset Firewall Config", question, mayWarning},
			"reset-ip":         {"ResetIP", "Reset IP", question, mayWarning},
			"reset-wifi":       {"WReset", "Reset Wi-Fi Config", question, mayWarning},
			"restart-gateway":  {"Restart", "Restart", question, willWarning},
		}

		err := rc.submitForm(action, answerYes, page, path, resetActions[action])
		if err != nil {
			log.Fatalf("Submission failed: %v", err)
		}
	} else {

		// Get body using the new getPath function
		body, err := rc.getPath(path)

		if err != nil {
			return err
		}

		// Extract content-sub div
		content, err := extractContentSub(body)

		if err != nil {
			log.Fatal(err)
		}

		if err := extractData(action, content, filter, natActionPrefix, pretty); err != nil {
			log.Fatalf("Error extracting %s", action)
		}

	}

	return nil
}
