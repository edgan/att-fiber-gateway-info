package main

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/gob"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type GatewayClient struct {
	client        *http.Client
	baseURL       string
	cookieFile    string
	loadedCookies int
}

// debugLog prints debug messages when debug mode is enabled
func debugLog(enabled bool, message string) {
	if enabled {
		fmt.Printf("Debug: %s\n", message)
	}
}

func NewGatewayClient(baseURL string, cookieFile string, debug bool, freshCookies bool) (*GatewayClient, error) {
	loadedCookies := 0

	// Create cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %v", err)
	}

	if !freshCookies {
		// Load cookies from file if it exists
		if _, err := os.Stat(cookieFile); err == nil {
			if err := loadCookies(jar, baseURL, cookieFile); err != nil {
				log.Printf("Failed to load cookies: %v", err)
			} else {
				loadedCookies = 1
				debugLog(debug, "Stored cookies use")
			}
		}
	}

	// Configure transport to skip SSL verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// Create HTTP client
	client := &http.Client{
		Transport: tr,
		Jar:       jar,
	}

	return &GatewayClient{
		client:        client,
		baseURL:       baseURL,
		cookieFile:    cookieFile,
		loadedCookies: loadedCookies,
	}, nil
}

// Saves the current cookies to the file after login
func (rc *GatewayClient) saveSessionCookies() error {
	return saveCookies(rc.client.Jar, rc.baseURL, rc.cookieFile)
}

// Function to print a row with padding
func printRowWithPadding(row []string, columnWidths []int) {
	for i, cell := range row {
		fmt.Print(cell + strings.Repeat(" ", columnWidths[i]-len(cell)+2))
	}
	fmt.Println()
}

func extractTableData(action string, doc *goquery.Document, filter string, pretty bool, tableClass string) {
	doc.Find("table").Each(func(i int, table *goquery.Selection) {
		var tableData [][]string // A slice to hold all rows, with each row as a slice of cell values

		// Check if this table has the specified class
		if class, exists := table.Attr("class"); exists && class == tableClass {
			// If the table matches, process its rows
			table.Find("tr").Each(func(j int, row *goquery.Selection) {
				var rowData []string // A slice to hold all cells in this row

				// Extract each "th" or "td" element and add its text to rowData
				row.Find("th").Each(func(k int, cell *goquery.Selection) {
					rowData = append(rowData, strings.TrimSpace(cell.Text()))
				})
				row.Find("td").Each(func(k int, cell *goquery.Selection) {
					var cellText string

					// Check if there is a <pre> tag inside the <td>
					pre := cell.Find("pre")
					if pre.Length() > 0 {
						// Process <pre> content with <br> tags replaced by newlines
						htmlContent, err := pre.Html() // Get HTML inside <pre>, which returns (string, error)
						if err != nil {
							cellText = "" // or handle error appropriately
						} else {
							cellText = strings.ReplaceAll(htmlContent, "Wi-Fi<br/>", "Wi-Fi: ") // Special case for WiFi
							cellText = strings.ReplaceAll(cellText, "<br/>", "\n  ") // Replace <br /> with newline
							cellText = strings.ReplaceAll(cellText, "<br>", "\n")     // Handle <br> tag as well
						}
					} else {
						// No <pre> tag; get text normally
						cellText = strings.TrimSpace(cell.Text())
					}

					// Append processed cellText to rowData
					rowData = append(rowData, strings.TrimSpace(cellText))
				})

				// Append the row to the table data
				tableData = append(tableData, rowData)
			})

			if action == "nat-check" || action == "nat-totals" {
				if class == "table60" {
					for _, row := range tableData {
						if len(row) > 0 && row[0] == "Total sessions in use" {
							if action == "nat-check" {
								fmt.Printf("%s.0\n", row[1])

								maxConnections := 8192
								connections, err := strconv.Atoi(row[1])
								if err != nil {
									panic(err)
								}

								if connections >= maxConnections {
									fmt.Printf("\nError: Too many connections\n")
									os.Exit(1)
								}
							}

							if action == "nat-totals" {
								fmt.Printf("%s: %s\n", "Total number of connections", row[1])
							}
						}
					}
				}

				if action == "nat-totals" && class == "grid table100" {
					// Initialize counters
					var tcpCount, udpCount int

					for _, row := range tableData {
						if row[1] == "tcp" {
							tcpCount++
						}
						if row[1] == "udp" {
							udpCount++
						}
					}

					fmt.Printf("Total number of tcp connections: %d\n", tcpCount)
					fmt.Printf("Total number of udp connections: %d\n", udpCount)
				}
			} else if action == "nat-connections" && !pretty {
				for _, row := range tableData {
					line := strings.Join(row, ", ")
					fmt.Println(line)
				}
			} else if action == "nat-destinations" {
				sortedDestinationsIPs := CountIPsByColumn(tableData, 7)
				fmt.Println("Destinations IP addresses:")
				for _, row := range sortedDestinationsIPs {
					fmt.Printf("%d %s\n", row.Count, row.IP)
				}
			} else if action == "nat-sources" {
				sortedSourcesIPs := CountIPsByColumn(tableData, 5)
				fmt.Println("Source IP addresses:")
				for _, row := range sortedSourcesIPs {
					fmt.Printf("%d %s\n", row.Count, row.IP)
				}
			} else {
				if action == "home-network-status" {
					for _, row := range tableData {
						count := len(row)
						if row[0] != "" && count > 1 {
							row[0] = row[0] + ":"
						}
					}
				}

				// Print each row in a default format
				for _, row := range tableData {
					if action == "device-list" {
						count := len(row)

						// ipv4 address / name
						if row[0] == "IPv4 Address / Name" {
							row[0] = "IPv4 Address"
						}

						substring := " / "

						if count > 1 {
							if strings.Contains(row[1], substring) {
								row[1] = strings.Replace(row[1], substring, "Name: ", 1)
							}
						}
					}

					line := strings.Join(row, ": ")

					if action == "device-list" {
						// connection-type
						if row[0] == "Connection Type" {
							line = strings.Join(row, ": \n  ")
						}
					}

					if action == "home-network-status" || action == "ip-allocation" || (action == "nat-connections" && pretty) {
						// Find the maximum width of each column
						columnWidths := make([]int, len(tableData[0]))
						for _, row := range tableData {
							for i, cell := range row {
								if len(cell) > columnWidths[i] {
									columnWidths[i] = len(cell)
								}
							}
						}

						printRowWithPadding(row, columnWidths)
					} else if !strings.Contains(line, "Legal Disclaimer") {
						fmt.Println(line)
					}
				}
				fmt.Println()
			}
		}
	})
}

func extractData(action string, content string, filter string, natActionPrefix string, pretty bool) error {
	// Load the HTML content into goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return fmt.Errorf("failed to parse content: %v", err)
	}

	if action == "nat-check" {
		extractTableData(action, doc, filter, pretty, "table60")
	} else if action == "nat-totals" {
		extractTableData(action, doc, filter, pretty, "table60")
		extractTableData(action, doc, filter, pretty, "grid table100")
	} else if strings.HasPrefix(action, natActionPrefix) || action == "ip-allocation" {
		extractTableData(action, doc, filter, pretty, "grid table100")
	} else if action == "fiber-status" {
		extractTableData(action, doc, filter, pretty, "table75")

		// Extract <h1> elements with specific information about Temperature, Vcc, and Power outside the tables
		doc.Find("h1").Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())

			// Only process h1 text if it contains specific keywords
			if strings.Contains(text, "Currently") {
				text = strings.ReplaceAll(text, "\u00A0\u00A0Currently", ":")
				fmt.Println(text)
			}
		})
	} else if action == "home-network-status" || action == "device-list" {
		extractTableData(action, doc, filter, pretty, "table100")
	} else {
		// broadband-status, sys-info
		extractTableData(action, doc, filter, pretty, "table75")
	}

	return nil
}

func extractContentSub(htmlStr string) (string, error) {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return "", err
	}

	// Find the content-sub div
	contentSub := doc.Find("#content-sub").First()

	// Get the HTML of the content-sub div
	htmlContent, err := contentSub.Html()
	if err != nil {
		return "", err
	}

	return htmlContent, nil
}

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

func (rc *GatewayClient) getNonce(loginPath string) (string, error) {
	if rc.loadedCookies == 0 {
		// First request to load the login page and get cookies
		resp1, err1 := rc.client.Get(rc.baseURL + loginPath)
		if err1 != nil {
			return "", fmt.Errorf("failed to get login page: %v", err1)
		}
		defer resp1.Body.Close()
	}

	// Second request to get the nonce, using the cookies from the first request
	resp2, err2 := rc.client.Get(rc.baseURL + loginPath)
	if err2 != nil {
		return "", fmt.Errorf("failed to get nonce from login page: %v", err2)
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

func calculateHash(password, nonce string) string {
	// Replicate JavaScript hex_md5(password + nonce)
	hasher := md5.New()
	io.WriteString(hasher, password+nonce)
	return hex.EncodeToString(hasher.Sum(nil))
}

func (rc *GatewayClient) login(password string, loginPath string) error {
	// Get nonce from login page
	nonce, err := rc.getNonce(loginPath)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
	}

	// Calculate hash
	hash := calculateHash(password, nonce)

	// Prepare form data
	formData := url.Values{
		"nonce":        {nonce},
		"password":     {strings.Repeat("*", len(password))}, // Replicate JS behavior
		"hashpassword": {hash},
		"Continue":     {"Continue"},
	}

	// Submit login form
	resp, err := rc.client.PostForm(rc.baseURL+loginPath, formData)
	if err != nil {
		return fmt.Errorf("failed to submit login form: %v", err)
	}
	defer resp.Body.Close()

	// Check if login was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status: %d", resp.StatusCode)
	}

	// Save session cookies
	if err := rc.saveSessionCookies(); err != nil {
		log.Printf("Failed to save cookies: %v", err)
	}

	return nil
}

func (rc *GatewayClient) restart(resetPath string, loginPath string) error {

	// Get nonce from reset page
	nonce, err := rc.getNonce(resetPath)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
  }
  
	// Prepare form data
	formData := url.Values{
		"Restart":      {"Restart"},
		"nonce":        {nonce},
	}

	// Submit restart form
	resp, err := rc.client.PostForm(rc.baseURL+resetPath, formData)
	if err != nil {
		return fmt.Errorf("failed to submit restart form: %v", err)
	}
  
  defer resp.Body.Close()

	// Check if restart was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("restart failed with status: %d", resp.StatusCode)
	}

	return nil
}

// getPath is a generic function that fetches and returns the response body as a string.
func (rc *GatewayClient) getPath(path string, loginPath string) (string, error) {
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
	if strings.Contains(bodyStr, loginPath) {
		return "", fmt.Errorf("Login failed. Password likely wrong.")
	}

	return bodyStr, nil
}

func (rc *GatewayClient) getPage(page string, action string, filter string, loginPath string, natActionPrefix string, pretty bool) error {
	path := "/cgi-bin/" + page + ".ha"

	// Get body using the new getPath function
	body, err := rc.getPath(path, loginPath)
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

	return nil
}

// Define a function that takes a column index as a parameter and returns sorted IP counts
func CountIPsByColumn(tableData [][]string, column int) []struct {
	IP    string
	Count int
} {
	// Variable to count occurrences of each IP address in the specified column
	ipCount := make(map[string]int)

	for i, row := range tableData {
		if i != 0 {
			ipCount[row[column]]++
		}
	}

	// Convert the map to a slice of structs for sorting
	var sortedIPs []struct {
		IP    string
		Count int
	}
	for ip, count := range ipCount {
		sortedIPs = append(sortedIPs, struct {
			IP    string
			Count int
		}{IP: ip, Count: count})
	}

	// Sort the slice by count in descending order
	sort.Slice(sortedIPs, func(i, j int) bool {
		return sortedIPs[i].Count > sortedIPs[j].Count
	})

	return sortedIPs
}

// Function to load cookies from a file
func loadCookies(jar http.CookieJar, baseURL string, filePath string) error {
	file, err := os.Open(filePath)
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

	u, err := url.Parse(baseURL)
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

func main() {	
	actionSlice := []string{"broadband-status", "device-list", "fiber-status", "home-network-status", "ip-allocation", "nat-check", "nat-connections", "nat-destinations", "nat-sources", "nat-totals", "system-information", "restart"}
	actions_help := []string{}

	for _, action := range actionSlice {
		actions_help = append(actions_help, action)
	}

	actionDescription := fmt.Sprintf("Action to perform (%s)", strings.Join(actions_help, ", "))

	filterSlice := []string{"icmp", "ipv4", "ipv6", "tcp", "udp"}
	filters_help := []string{}

	for _, filter := range filterSlice {
		filters_help = append(filters_help, filter)
	}

	filterDescription := fmt.Sprintf("Filter to perform (%s)", strings.Join(filters_help, ", "))

	// Define a map linking actions to their corresponding page names
	actionPageMap := map[string]string{
		"broadband-status":    "broadbandstatistics",
		"device-list":         "devices",
		"fiber-status":        "fiberstat",
		"home-network-status": "lanstatistics",
		"ip-allocation":       "ipalloc",
		"system-information":  "sysinfo",
	}

	// All "nat-" prefixed actions use "nattable" page
	natActionPrefix := "nat-"

	for _, action := range actionSlice {
		if strings.HasPrefix(action, natActionPrefix) {
			actionPageMap[action] = "nattable"
		}
	}

	gatewayBaseURL := "https://192.168.1.254"
	loginPath := "/cgi-bin/login.ha"
	resetPath := "/cgi-bin/reset.ha"

	cookieFilename := "/var/tmp/.att-fiber-gateway-info_cookies.gob"

	// Parse command line arguments
	baseURL := flag.String("url", gatewayBaseURL, "Gateway base URL")
	password := flag.String("password", "", "Gateway password")
	action := flag.String("action", "", actionDescription)
	filter := flag.String("filter", "", filterDescription)
	cookieFile := flag.String("cookiefile", cookieFilename, "File to save session cookies")
	debug := flag.Bool("debug", false, "Enable debug mode")
	pretty := flag.Bool("pretty", false, "Enable pretty mode for nat-connections")
	freshCookies := flag.Bool("fresh", false, "Do not use existing cookies (Warning: If you use all the time you will run out of sessions. There is a max.)")
	flag.Parse()

	if *password == "" {
		log.Fatal("Password is required")
	}

	isValidAction := false

	for _, a := range actionSlice {
		if *action == a {
			isValidAction = true
			break
		}
	}

	if !isValidAction {
		actionError := fmt.Sprintf("Action must be one of these (%s)", strings.Join(actions_help, ", "))
		log.Fatal(actionError)
	}

	isValidFilter := false

	// Validate filter only if it's provided
	isValidFilter = *filter == "" // Default to valid if filter is empty (optional)
	if *filter != "" {
		for _, f := range filterSlice {
			if *filter == f {
				isValidFilter = true
				break
			}
		}
	}

	if !isValidFilter {
		filterError := fmt.Sprintf("Filter must be one of these (%s)", strings.Join(filters_help, ", "))
		log.Fatal(filterError)
	}

	// Create router client
	client, err := NewGatewayClient(*baseURL, *cookieFile, *debug, *freshCookies)
	if err != nil {
		log.Fatalf("Failed to create router client: %v", err)
	}

	// Perform login
	if err := client.login(*password, loginPath); err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	// Perform restart
	if *action == "restart" {
		if err := client.restart(resetPath, loginPath); err != nil {
			log.Fatalf("Restart failed: %v", err)
		}
		os.Exit(0)
  	}

	// Determine the page based on action
	page, exists := actionPageMap[*action]
	if !exists {
		log.Fatalf("Unknown action: %s", *action)
	}

	// Get the specified page
	if err := client.getPage(page, *action, *filter, loginPath, natActionPrefix, *pretty); err != nil {
		log.Fatalf("Failed to get %s: %v", *action, err)
	}

	os.Exit(0)
}
