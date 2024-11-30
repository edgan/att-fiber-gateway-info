package main

import (
	"fmt"
	"strings"

	"edgan/att-fiber-gateway-info/internal/logging"
)

func checkForLoginFailure(body string, flags *flags, loginPath string) error {
	// Check for login failure
	if strings.Contains(body, loginPath) {
		logging.DebugLog(*flags.Debug, "LoginPath in body of page")
		return fmt.Errorf("Login failed. Password likely wrong.")
	}
	return nil
}
