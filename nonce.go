package main

import (
	"fmt"

	"golang.org/x/net/html"
)

// Extracts the nonce value from the HTML document
func findNonce(n *html.Node) (string, error) {
	nonce := searchNonce(n)
	if nonce == "" {
		return "", fmt.Errorf("nonce not found in HTML")
	}
	return nonce, nil
}

// Recursively searches for the nonce value in the HTML node tree
func searchNonce(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "input" {
		if isNonceInput, nonceValue := checkNonceInput(n); isNonceInput && nonceValue != "" {
			return nonceValue
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := searchNonce(c); result != "" {
			return result
		}
	}
	return ""
}

// Checks if the input node is a nonce input and retrieves its value
func checkNonceInput(n *html.Node) (bool, string) {
	var isNonceInput bool
	var nonceValue string
	for _, attr := range n.Attr {
		if attr.Key == "name" && attr.Val == "nonce" {
			isNonceInput = true
		}
		if attr.Key == "value" {
			nonceValue = attr.Val
		}
	}
	return isNonceInput, nonceValue
}

func (rc *gatewayClient) getNonce(page string) (string, error) {
	path := returnPath(page)

	if page == "login" && rc.loadedCookies == 0 {
		// First request to load the login page and get cookies
		resp1, err1 := rc.client.Get(rc.baseURL + path)

		if err1 != nil {
			return "", fmt.Errorf("failed to get login page: %v", err1)
		}

		defer resp1.Body.Close()
	}

	// Second request to get the nonce, using the cookies from the first request
	resp2, err2 := rc.client.Get(rc.baseURL + path)

	if err2 != nil {
		return "", fmt.Errorf("failed to get nonce from page: %v", err2)
	}

	defer resp2.Body.Close()

	// Parse HTML to extract nonce
	doc, err := html.Parse(resp2.Body)

	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %v", err)
	}

	// Use the findNonce function to get the nonce
	nonce, err := findNonce(doc)

	if err != nil {
		return "", err
	}

	return nonce, nil
}
