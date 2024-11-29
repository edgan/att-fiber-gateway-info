package main

import (
	"fmt"
	"strings"
)

func checkForLoginFailure(body string, flags *flags, loginPath string) error {
	// Check for login failure
	if strings.Contains(body, loginPath) {
		debugLog(*flags.Debug, "LoginPath in body of page")
		return fmt.Errorf("Login failed. Password likely wrong.")
	}
	return nil
}
