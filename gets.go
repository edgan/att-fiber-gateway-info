package main

import (
	"fmt"
	"io"
	"strings"
)

func (rc *GatewayClient) getPath(flags *Flags, path string) (string, error) {
	resp, err := rc.client.Get(rc.baseURL + path)

	if err != nil {
		debugLog(*flags.Debug, "Failed to Get")
		return "", fmt.Errorf("failed to get path %s: %v", path, err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		debugLog(*flags.Debug, "Failed to ReadAll")
		return "", fmt.Errorf("failed to read response for path %s: %v", path, err)
	}

	bodyStr := string(body)

	// Check for login failure
	if strings.Contains(bodyStr, rc.loginPath) {
		debugLog(*flags.Debug, "LoginPath in body of page")
		return "", fmt.Errorf("Login failed. Password likely wrong.")
	}

	return bodyStr, nil
}

func (rc *GatewayClient) getPage(action string, configs Configs, flags *Flags, model string, page string, returnFact string) (string, error) {
	fact := ""
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

		err := rc.submitForm(action, flags, page, path, resetActions[action])
		if err != nil {
			logFatalf("Submission failed: %v", err)
		}
	} else {
		// Get body using the new getPath function
		body, err := rc.getPath(flags, path)

		if err != nil {
			debugLog(*flags.Debug, "Failed to getPath")
			logFatal(err)
		}

		// Extract content-sub div
		content, err := extractContentSub(body)

		if err != nil {
			debugLog(*flags.Debug, "Failed to extractContentSub")
			logFatal(err)
		}

		fact, err = extractData(action, configs, content, flags, model, returnFact)
		if err != nil {
			logFatalf("Error extracting %s", action)
		}

	}

	return fact, nil
}
