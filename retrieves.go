package main

import (
	"log"
)

func (rc *GatewayClient) retrieveAction(action string, actionPages map[string]string, configs Configs, flags *Flags, model string, natActionPrefix string, returnFact string) (string, error) {
	fact := ""

	// Get the specified page based on action
	page := returnActionPage(action, actionPages)

	// login is not required for most pages
	loginRequired := false

	// pages that require login
	loginPages := []string{"ipalloc", "nattable", "reset"}

	for _, loginPage := range loginPages {
		if page == loginPage {
			if configs.Password == "" {
				log.Fatal("Password is required")
			}
			loginRequired = true
		}
	}

	if loginRequired {
		debugLog(*flags.Debug, "LoginRequired true")
		performLogin(rc, configs)
	}

	fact, err := rc.getPage(action, configs, flags, model, natActionPrefix, page, returnFact)

	return fact, err
}