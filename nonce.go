package main

import (
	"fmt"

	"golang.org/x/net/html"
)

// Extracts the nonce value from the HTML document
func findNonce(n *html.Node) (string, error) {
	var nonce string

	var searchNode func(*html.Node)
	searchNode = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "input" {
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
			if isNonceInput && nonceValue != "" {
				nonce = nonceValue
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			searchNode(c)
		}
	}

	searchNode(n)

	if nonce == "" {
		return "", fmt.Errorf("nonce not found in HTML")
	}
	return nonce, nil
}

func (rc *GatewayClient) getNonce(page string) (string, error) {
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
