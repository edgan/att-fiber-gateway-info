package main

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"runtime"

	"edgan/att-fiber-gateway-info/internal/logging"
)

// Saves the current cookies to the file after login
func (rc *gatewayClient) saveSessionCookies() error {
	return saveCookies(rc.client.Jar, rc.baseURL, rc.cookieFile)
}

func createAndLoadCookies(configs configs, flags *flags) (*cookiejar.Jar, bool, error) {
	// Create cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, false, fmt.Errorf("failed to create cookie jar: %v", err)
	}

	loadedCookies := false

	if !*flags.FreshCookies {
		// Load cookies from file if it exists
		if _, err := os.Stat(*flags.CookieFile); err == nil {
			if err := loadCookies(configs, flags, jar); err != nil {
				return nil, loadedCookies, fmt.Errorf("failed to load cookies: %v", err)
			}

			loadedCookies = true
			logging.DebugLog(*flags.Debug, "Stored cookies use") // Assuming debug is true for simplicity
		}
	}

	return jar, loadedCookies, nil
}

// Function to load cookies from a file
func loadCookies(configs configs, flags *flags, jar http.CookieJar) error {
	file, err := os.Open(*flags.CookieFile)

	if err != nil {
		return err
	}

	defer file.Close()

	var cookies []*http.Cookie
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&cookies)

	if err != nil {
		return err
	}

	u, err := url.Parse(configs.BaseURL)

	if err != nil {
		return err
	}
	jar.SetCookies(u, cookies)
	return nil
}

func saveCookies(jar http.CookieJar, baseURL string, filePath string) error {
	// Parse the URL
	u, err := url.Parse(baseURL)

	if err != nil {
		return fmt.Errorf("failed to parse base URL: %v", err)
	}

	// Retrieve cookies from the jar
	cookies := jar.Cookies(u)

	// Open the file to save cookies
	file, err := os.Create(filePath)

	if err != nil {
		return fmt.Errorf("failed to create cookie file: %v", err)
	}

	defer file.Close()

	// Encode and save cookies to the file
	encoder := gob.NewEncoder(file)

	if err := encoder.Encode(cookies); err != nil {
		return fmt.Errorf("failed to encode cookies: %v", err)
	}

	return nil
}

// determineCookiePath sets the cookie file path based on GOOS
func determineCookiePath() string {
	cookieFilename := "att-fiber-gateway-info_cookies.gob"
	var cookiePath string

	switch runtime.GOOS {
	case "windows":
		cookiePath = "C:\\Windows\\Temp\\" + cookieFilename
	default: // Linux and other Unix-like systems
		cookiePath = "/var/tmp/" + cookieFilename
	}

	return cookiePath
}
