package main

import (
	"fmt"
	"io"
	"strings"

	"edgan/att-fiber-gateway-info/internal/logging"
)

func (rc *gatewayClient) getPath(flags *flags, path string) (string, error) {
	resp, err1 := rc.client.Get(rc.baseURL + path)

	if err1 != nil {
		logging.DebugLog(*flags.Debug, "Failed to Get")
		return empty, fmt.Errorf("failed to get path %s: %v", path, err1)
	}

	defer resp.Body.Close()

	body, err2 := io.ReadAll(resp.Body)

	if err2 != nil {
		logging.DebugLog(*flags.Debug, "Failed to ReadAll")
		return empty, fmt.Errorf("failed to read response for path %s: %v", path, err2)
	}

	bodyStr := string(body)

	err3 := checkForLoginFailure(bodyStr, flags, rc.loginPath)

	if err3 != nil {
		return empty, err3
	}

	return bodyStr, nil
}

func (rc *gatewayClient) getPage(
	action string, configs configs, flags *flags, model string, page string, returnFact string,
) (string, error) {
	fact := empty
	path := returnPath(page)

	if page == "reset" {
		mayWarning := "This may take down your internet immediately."
		willWarning := "Note this will take down your internet immediately."

		parts := strings.Split(action, dash)

		question := formatQuestion(parts[task], parts[resource])

		// The buttonNames and buttonValues have to be exact
		resetActions := map[string][actionAttributes]string{
			"reset-connection": {"ResetConn", "Reset Connection", question, mayWarning},
			"reset-device":     {"Reset", "Reset Device...", question, mayWarning},
			"reset-firewall":   {"FReset", "Reset Firewall Config", question, mayWarning},
			"reset-ip":         {"ResetIP", "Reset IP", question, mayWarning},
			"reset-wifi":       {"WReset", "Reset Wi-Fi Config", question, mayWarning},
			"restart-gateway":  {"Restart", "Restart", question, willWarning},
		}

		err := rc.submitForm(action, flags, page, path, resetActions[action])

		if err != nil {
			logging.LogFatalf("Submission failed: %v", err)
		}
	} else {
		// Get body using the new getPath function
		body, err := rc.getPath(flags, path)

		if err != nil {
			logging.DebugLog(*flags.Debug, "Failed to getPath")
			logging.LogFatal(err)
		}

		// Extract content-sub div
		content, err := extractContentSub(body)

		if err != nil {
			logging.DebugLog(*flags.Debug, "Failed to extractContentSub")
			logging.LogFatal(err)
		}

		fact, err = extractData(action, configs, content, flags, model, returnFact)
		if err != nil {
			logging.LogFatalf("Error extracting %s", action)
		}
	}

	return fact, nil
}
