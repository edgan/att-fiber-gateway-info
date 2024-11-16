package main

import (
	"bufio"
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
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"golang.org/x/net/html"
)

type GatewayClient struct {
	client        *http.Client
	baseURL       string
	colorMode     bool
	cookieFile    string
	loadedCookies int
	loginPath     string
}

// debugLog prints debug messages when debug mode is enabled
func debugLog(enabled bool, message string) {
	if enabled {
		fmt.Printf("Debug: %s\n", message)
	}
}

func NewGatewayClient(baseURL string, colorMode bool, cookieFile string, debug bool, freshCookies bool, loginPath string) (*GatewayClient, error) {
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
		colorMode:     colorMode,
		cookieFile:    cookieFile,
		loadedCookies: loadedCookies,
		loginPath:     loginPath,
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
							cellText = strings.ReplaceAll(cellText, "<br/>", "\n  ")            // Replace <br /> with newline
							cellText = strings.ReplaceAll(cellText, "<br>", "\n")               // Handle <br> tag as well
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

func calculateHash(password, nonce string) string {
	// Replicate JavaScript hex_md5(password + nonce)
	hasher := md5.New()
	io.WriteString(hasher, password+nonce)
	return hex.EncodeToString(hasher.Sum(nil))
}

func (rc *GatewayClient) postForm(path string, formData url.Values) error {
	// Submit form
	resp, err := rc.client.PostForm(rc.baseURL+path, formData)

	if err != nil {
		return fmt.Errorf("failed to submit the form to %s: %v", path, err)
	}

	defer resp.Body.Close()

	// Check if submission was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("submission to %s failed with status: %d", path, resp.StatusCode)
	}

	if path == rc.loginPath {
		// Save session cookies
		if err := rc.saveSessionCookies(); err != nil {
			log.Printf("Failed to save cookies: %v", err)
		}
	}

	return nil
}

func (rc *GatewayClient) login(password string) error {
	page := "login"

	// Get nonce from page
	nonce, err := rc.getNonce(page)

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
		"Continue":     {"Continue"}, // submit button
	}

	if err := rc.postForm(rc.loginPath, formData); err != nil {
		log.Fatalf("Submission to %s failed: %v", rc.loginPath, err)
	}

	return nil
}

// getPath is a generic function that fetches and returns the response body as a string.
func (rc *GatewayClient) getPath(path string) (string, error) {
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
	if strings.Contains(bodyStr, rc.loginPath) {
		return "", fmt.Errorf("Login failed. Password likely wrong.")
	}

	return bodyStr, nil
}

func (rc *GatewayClient) submitForm(action string, answerYes bool, page string, path string, reset_action [4]string) error {
	buttonName := reset_action[0]
	buttonValue := reset_action[1]
	question := reset_action[2]
	warning := reset_action[3]

	if askYesNo(answerYes, rc.colorMode, question, warning) {
		// Get nonce from page
		nonce, err := rc.getNonce(page)
		if err != nil {
			return fmt.Errorf("failed to get nonce: %v", err)
		}

		// Prepare form data
		formData := url.Values{
			"nonce":    {nonce},
			buttonName: {buttonValue}, // Dynamically use the submit button
		}

		// Submit the form
		if err := rc.postForm(path, formData); err != nil {
			return fmt.Errorf("submission to %s in %s failed: %v", action, path, err)
		}
	}

	return nil
}

func formatQuestion(task string, resource string) string {
	return fmt.Sprintf("Do you want to %s the %s?", task, resource)
}

func (rc *GatewayClient) getPage(action string, answerYes bool, filter string, natActionPrefix string, page string, pretty bool) error {
	path := returnPath(page)

	if page == "reset" {
		may_warning := "This may take down your internet immediately."
		will_warning := "Note this will take down your internet immediately."

		parts := strings.Split(action, "-")
		task := parts[0]
		resource := parts[1]

		question := formatQuestion(task, resource)

		// The buttonNames and buttonValues have to be exact
		reset_actions := map[string][4]string{
			"reset-connection": {"ResetConn", "Reset Connection", question, may_warning},
			"reset-device":     {"Reset", "Reset Device...", question, may_warning},
			"reset-firewall":   {"FReset", "Reset Firewall Config", question, may_warning},
			"reset-ip":         {"ResetIP", "Reset IP", question, may_warning},
			"reset-wifi":       {"WReset", "Reset Wi-Fi Config", question, may_warning},
			"restart-gateway":  {"Restart", "Restart", question, will_warning},
		}

		err := rc.submitForm(action, answerYes, page, path, reset_actions[action])
		if err != nil {
			log.Fatalf("Submission failed: %v", err)
		}
	} else {

		// Get body using the new getPath function
		body, err := rc.getPath(path)

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

func returnPath(page string) string {
	path := "/cgi-bin/" + page + ".ha"

	return path
}

// Prompt the user for yes/no input
func askYesNo(answerYes bool, colorMode bool, question string, warning string) bool {
	if colorMode {
		red := color.New(color.FgRed)
		warning = red.Sprint(warning)
	}

	if !answerYes {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print(question + " " + warning + " (yes/no): ")
			input, err := reader.ReadString('\n')

			if err != nil {
				fmt.Println("Error reading input. Please try again.")
				continue
			}

			// Trim whitespace and convert to lowercase
			input = strings.TrimSpace(strings.ToLower(input))

			if input == "yes" || input == "y" {
				return true
			} else if input == "no" || input == "n" {
				return false
			} else {
				fmt.Println("Invalid input. Please type 'yes' or 'no'.")
			}
		}
	}

	return true
}

func ColoredUsage() {
	// Create color functions
	blue := color.New(color.FgBlue)
	boldGreen := color.New(color.FgGreen, color.Bold)
	cyan := color.New(color.FgCyan)
	green := color.New(color.FgGreen)

	// Print usage header
	fmt.Printf("Usage of %s:\n", green.Sprintf(os.Args[0]))

	flag.VisitAll(func(f *flag.Flag) {
		// Format flag name with color
		s := fmt.Sprintf("  ")
		s += boldGreen.Sprintf("-%s", f.Name)

		// Add default value if it exists and isn't empty
		if f.DefValue != "" {
			s += blue.Sprintf(" (default: %v)", f.DefValue)
		}

		// Add the usage description in blue
		if f.Usage != "" {
			s += "\n    \t" + cyan.Sprintf(f.Usage)
		}

		fmt.Println(s)
	})
}

func main() {
	colorTerminal := isColorTerminal()

	colorMode := true

	if !colorTerminal {
		colorMode = false
	}

	actions := []string{
		"broadband-status", "device-list", "fiber-status", "home-network-status",
		"ip-allocation", "nat-check", "nat-connections", "nat-destinations",
		"nat-sources", "nat-totals", "reset-connection", "reset-device",
		"reset-firewall", "reset-ip", "reset-wifi", "restart-gateway",
		"system-information",
	}

	actions_help := []string{}

	for _, action := range actions {
		actions_help = append(actions_help, action)
	}

	actionDescription := fmt.Sprintf("Action to perform (%s)", strings.Join(actions_help, ", "))

	filters := []string{"icmp", "ipv4", "ipv6", "tcp", "udp"}
	filters_help := []string{}

	for _, filter := range filters {
		filters_help = append(filters_help, filter)
	}

	filterDescription := fmt.Sprintf("Filter to perform (%s)", strings.Join(filters_help, ", "))

	// Define a map linking actions to their corresponding page names
	actionPages := map[string]string{
		"broadband-status":    "broadbandstatistics",
		"device-list":         "devices",
		"fiber-status":        "fiberstat",
		"home-network-status": "lanstatistics",
		"ip-allocation":       "ipalloc",
		"restart-gateway":     "reset",
		"system-information":  "sysinfo",
	}

	// All "nat-" prefixed actions use "nattable" page
	natActionPrefix := "nat-"

	// All "reset-" prefixed actions use "reset" page
	resetActionPrefix := "reset-"

	for _, action := range actions {
		if strings.HasPrefix(action, natActionPrefix) {
			actionPages[action] = "nattable"
		}
		if strings.HasPrefix(action, resetActionPrefix) {
			actionPages[action] = "reset"
		}
	}

	gatewayBaseURL := "https://192.168.1.254"

	var cookieFilename string

	// Set cookieFilename based on GOOS
	switch runtime.GOOS {
	case "windows":
		cookieFilename = "C:\\Windows\\Temp\\att-fiber-gateway-info_cookies.gob"
	default: // Linux and other Unix-like systems
		cookieFilename = "/var/tmp/.att-fiber-gateway-info_cookies.gob"
	}

	loginPath := returnPath("login")

	// Parse command line arguments
	baseURL := flag.String("url", gatewayBaseURL, "Gateway base URL")
	password := flag.String("password", "", "Gateway password")
	action := flag.String("action", "", actionDescription)
	filter := flag.String("filter", "", filterDescription)
	cookieFile := flag.String("cookiefile", cookieFilename, "File to save session cookies")
	debug := flag.Bool("debug", false, "Enable debug mode")
	answerYes := flag.Bool("yes", false, "Answer yes to any questions")
	pretty := flag.Bool("pretty", false, "Enable pretty mode for nat-connections")
	freshCookies := flag.Bool(
		"fresh", false,
		"Do not use existing cookies (Warning: If you use all the time you will run out of sessions. There is a max.)",
	)

	if colorMode {
		// Replace the default Usage with our colored version
		flag.Usage = ColoredUsage
	}

	flag.Parse()

	if *password == "" {
		log.Fatal("Password is required")
	}

	isValidAction := false

	for _, a := range actions {
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
		for _, f := range filters {
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
	client, err := NewGatewayClient(*baseURL, colorMode, *cookieFile, *debug, *freshCookies, loginPath)

	if err != nil {
		log.Fatalf("Failed to create router client: %v", err)
	}

	// Perform login
	if err := client.login(*password); err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	// Determine the page based on action
	page, exists := actionPages[*action]

	if !exists {
		log.Fatalf("Unknown action: %s", *action)
	}

	// Get the specified page
	if err := client.getPage(*action, *answerYes, *filter, natActionPrefix, page, *pretty); err != nil {
		log.Fatalf("Failed to get %s: %v", *action, err)
	}

	os.Exit(0)
}
