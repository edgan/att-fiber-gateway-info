package main

import (
	"log"
)

func (rc *GatewayClient) retrieveAction(action string, actionPages map[string]string, answerNo bool, answerYes bool, datadog bool, filter string, metrics bool, model string, natActionPrefix string, password string, pretty bool, returnFact string, statsdIPPort string) (string, error) {
	fact := ""

	// Get the specified page based on action
	page := getActionPage(action, actionPages)

	// login is not required for most pages
	loginRequired := false

	// pages that require login
	loginPages := []string{"ipalloc", "nat-table", "reset"}

	for _, loginPage := range loginPages {
		if page == loginPage {
			if password == "" {
				log.Fatal("Password is required")
			}
			loginRequired = true
		}
	}

	if loginRequired {
		performLogin(rc, password)
	}

	fact, err := rc.getPage(action, answerNo, answerYes, datadog, filter, metrics, model, natActionPrefix, page, password, pretty, returnFact, statsdIPPort)

	return fact, err
}
